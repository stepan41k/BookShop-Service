package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`
	HttpServer `yaml:"http_server"`
	DataBase `yaml:"db"`
}

type HttpServer struct {
	Adress string	`yaml:"adress" env-default:"localhost:8080"`
    Timeout time.Duration `yaml:"timeout" env-default:"4s"`
    Idle_timeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DataBase struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Username string `yaml:"username"`
	Password string
    DBName string `yaml:"dbname"`
    SSLMode string `yaml:"sslmode"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())	
	}

	os.Setenv("CONFIG_PATH", "./config/local.yaml")
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("config path is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}