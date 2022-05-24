package app

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type App struct {
	config *Config
	router *gin.Engine
	db     *sqlx.DB
	redis  *redis.Client
}

func Init() *App {
	config := loadConfig()
	db := loadDB(config)
	redis := loadRedis(config)
	mailer := loadMailer(config)
	storage := loadStorage(redis)
	verifier := loadVerifier(mailer, storage)
	userRep, tokenRep := loadRepositories(db)
	authService := loadAuthService(userRep, tokenRep, verifier, config.GlobalSalt)
	router := loadRouter(config, authService)
	return &App{config: config, router: router, db: db, redis: redis}
}

func (app *App) RunMigrations() {
	driver, err := postgres.WithInstance(app.db.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file:///"+app.config.MigrationsDir,
		"postgres", driver)
	if err != nil {
		panic(err)
	}

	log.Println("Running migrations...")
	err = migration.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Noting to migrate")
		} else {
			panic(err)
		}
		return
	}
	log.Println("Migration completed")
}

func (app *App) Close() {
	err := app.db.Close()
	if err != nil {
		panic(err)
	}
	log.Println("DB closed")

	err = app.redis.Close()
	if err != nil {
		panic(err)
	}
	log.Println("Redis closed")
}

func (app *App) Run() {
	runServer(&http.Server{
		Addr:    app.config.HTTPHost,
		Handler: app.router,
	}, func() {
		app.Close()
		log.Println("App closed")
	})
}

func runServer(server *http.Server, disposer func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()
	log.Print("Server started")

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %+v", err)
	}

	disposer()
	cancel()
	log.Print("Server stopped")
}
