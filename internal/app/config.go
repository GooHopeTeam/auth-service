package app

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
