package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Listen     ListenConfig
	Postgresql PostgresConfig
}

type ListenConfig struct {
	BindIP string `env:"BIND_IP" env-required:"true"`
	Port   string `env:"LISTEN_PORT" env-required:"true"`
}

func (l ListenConfig) Addr() string {
	return fmt.Sprintf("%s:%s", l.BindIP, l.Port)
}

type PostgresConfig struct {
	Host     string `env:"PSQL_HOST" env-required:"true"`
	Port     string `env:"PSQL_PORT" env-required:"true"`
	Database string `env:"PSQL_NAME" env-required:"true"`
	Username string `env:"PSQL_USER" env-required:"true"`
	Password string `env:"PSQL_PASSWORD" env-required:"true"`
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		p.Username, p.Password, p.Host, p.Port, p.Database)
}

func MustLoad(path string) *Config {
	var cfg Config

	log.Printf("Reading config from %s...", path)

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	log.Println("Config loaded successfully")
	return &cfg
}
