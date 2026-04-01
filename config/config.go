package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	JWT      JWTConfig      `yaml:"jwt"`
	Logger   LoggerConfig   `yaml:"logger"`
	Database DatabaseConfig `yaml:"database"`
	Slack    SlackConfig    `yaml:"slack"`
}

type environment string
type ServerConfig struct {
	Port        int         `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
	Environment environment `yaml:"environment" env:"ENVIRONMENT" env-default:"development"`
}

type JWTConfig struct {
	Secret string        `env:"JWT_SECRET" env-required:"true"`
	Expiry time.Duration `yaml:"expiry" env:"JWT_EXPIRY" env-default:"24h"`
}

type LoggerConfig struct {
	Level     string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	AddSource bool   `yaml:"add_source" env:"LOG_ADD_SOURCE" env-default:"false"`
}

type SlackConfig struct {
	BotToken string `env:"SLACK_BOT_TOKEN"`
}

type DatabaseConfig struct {
	URL      string `env:"DATABASE_URL"` // Production only
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     int    `yaml:"port" env:"DB_PORT"`
	Name     string `yaml:"name" env:"DB_NAME"`
	User     string `yaml:"user" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
}

func (db *DatabaseConfig) ConnectionString() string {
	if db.URL != "" {
		return db.URL
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		db.Host, db.User, db.Password, db.Name, db.Port)
}

func Load() *Config {
	// load `.env`
	_ = godotenv.Load()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/environments/dev.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config

	// Read YAML values then override conflict values
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot load config: %v", err)
	}

	return &cfg
}
