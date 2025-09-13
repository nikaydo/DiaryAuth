package config

import (
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Env struct {
	Port       string `env:"PORT"`
	Host       string `env:"HOST"`
	Postgresql string `env:"POSTGRESQL"`

	JWTSecret     string `env:"SECRET_JWT"`
	JWTTTL        string `env:"JWT_TTL"`
	RefreshSecret string `env:"SECRET_REFRESH"`
	RefreshTTL    string `env:"REFRESH_TTL"`
}

func ReadEnv() (Env, error) {
	err := godotenv.Load()
	cfg := Env{
		Port:       os.Getenv("PORT"),
		Host:       os.Getenv("HOST"),
		Postgresql: os.Getenv("POSTGRESQL"),
	}
	if err != nil {
		log.Println("Error read env:", err, ". Use default values")
		return cfg, nil
	}
	if err := env.Parse(&cfg); err != nil {
		return cfg, fmt.Errorf("error parse env: %w", err)
	}
	return cfg, nil
}
