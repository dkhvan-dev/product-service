package config

import "fmt"

type AppConfig struct {
	ServerPort  string `env:"SERVER_PORT" envDefault:"8000"`
	ReleaseMode bool   `env:"RELEASE_MODE" envDefault:"false"`

	DSName string `env:"DS_NAME" envDefault:"postgres"`
	DBPort string `env:"DB_PORT" envDefault:"5432"`
	DBHost string `env:"DB_HOST" envDefault:"localhost"`
	DBName string `env:"DB_NAME" envDefault:"products_db"`
	DBUser string `env:"DB_USERNAME" envDefault:"postgres"`
	DBPass string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBURL  string `env:"DB_URL"`
}

func (a *AppConfig) ComputeDBUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", a.DBUser, a.DBPass, a.DBHost, a.DBPort, a.DBName)
}
