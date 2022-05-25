package app

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/goohopeteam/auth-service/internal/handler"
	"github.com/goohopeteam/auth-service/internal/repository"
	"github.com/goohopeteam/auth-service/internal/service/auth"
	"github.com/goohopeteam/auth-service/internal/service/mailer"
	"github.com/goohopeteam/auth-service/internal/service/storage"
	"github.com/goohopeteam/auth-service/internal/service/verifier"
	"github.com/jmoiron/sqlx"
)

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

func loadRouter(config *Config, authService auth.AuthService) *gin.Engine {
	handler := handler.New(authService)
	router := gin.Default()
	if config.EnvType == "DEV" {
		router.Use(CORS())
	}
	router.POST("/register", handler.HandleRegistration)
	router.POST("/login", handler.HandleLogin)
	router.POST("/verify_email", handler.HandleEmailVerification)
	router.POST("/verify_token", handler.HandleTokenVerification)
	router.POST("/change_password", handler.HandlePasswordChange)
	return router
}
