package app

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-redis/redis"
	"github.com/goohopeteam/auth-service/internal/handler"
	"github.com/goohopeteam/auth-service/internal/repository"
	"github.com/goohopeteam/auth-service/internal/service/auth"
	"github.com/goohopeteam/auth-service/internal/service/mailer"
	"github.com/goohopeteam/auth-service/internal/service/storage"
	"github.com/goohopeteam/auth-service/internal/service/verifier"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	DB struct {
		Name     string `env:"DB_NAME,required"`
		Host     string `env:"DB_HOST" envDefault:"localhost"`
		Port     string `env:"DB_PORT" envDefault:"5432"`
		User     string `env:"DB_USER,required"`
		Password string `env:"DB_PASSWORD,required"`
	}
	Redis struct {
		Host     string `env:"REDIS_HOST" envDefault:"localhost"`
		Port     string `env:"REDIS_PORT" envDefault:"6379"`
		Password string `env:"REDIS_PASSWORD,required"`
		DB       int    `env:"REDIS_DB" envDefault:"0"`
	}
	Mailer struct {
		Sender   string `env:"MAILER_SENDER,required"`
		Password string `env:"MAILER_PASSWORD,required"`
		SMTP     struct {
			Host string `env:"SMTP_HOST" envDefault:"smtp.gmail.com"`
			Port string `env:"SMTP_PORT" envDefault:"587"`
		}
	}
	HTTPHost   string `env:"HTTP_HOST" envDefault:"localhost:8080"`
	GlobalSalt string `env:"GLOBAL_SALT,required"`
}

type App struct {
	config *Config
	router *gin.Engine
	db     *sqlx.DB
	redis  *redis.Client
}

func loadConfig() *Config {
	var config Config
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}
	return &config
}

func loadDB(config *Config) *sqlx.DB {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable sslmode=disable",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.Name)
	db, err := sqlx.Connect("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func loadRedis(config *Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	err := client.Ping().Err()
	if err != nil {
		panic(err)
	}
	return client
}

func loadMailer(config *Config) mailer.Mailer {
	return mailer.New(config.Mailer.SMTP.Host, config.Mailer.SMTP.Port, config.Mailer.Sender, config.Mailer.Password)
}

func loadRepositories(db *sqlx.DB) (repository.UserRepository, repository.TokenRepository) {
	userRep := repository.NewUserRepository(db)
	tokenRep := repository.NewTokenRepository(db)
	return userRep, tokenRep
}

func loadStorage(client *redis.Client) storage.Storage {
	return storage.NewRedisStorage(client)
}

func loadVerifier(mailer mailer.Mailer, storage storage.Storage) verifier.Verifier {
	return verifier.NewEmailVerifier(mailer, storage)
}

func loadAuthService(userRep repository.UserRepository, tokenRep repository.TokenRepository, emailVerifier verifier.Verifier, globalSalt string) auth.AuthService {
	return auth.New(userRep, tokenRep, emailVerifier, globalSalt)
}

func loadRouter(authService auth.AuthService) *gin.Engine {
	handler := handler.New(authService)
	router := gin.Default()
	router.POST("/register", handler.HandleRegistration)
	router.POST("/login", handler.HandleLogin)
	router.POST("/verify_email", handler.HandleEmailVerification)
	router.POST("/verify_token", handler.HandleTokenVerification)
	return router
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
	router := loadRouter(authService)
	return &App{config: config, router: router, db: db}
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
