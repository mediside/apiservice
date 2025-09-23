package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPathEnv = "CONFIG_PATH"

type Config struct {
	Env              string         `yaml:"env" env:"ENV"`
	PathologyLevel   float32        `yaml:"pathology_level"`
	ResearchSavePath string         `yaml:"research_save_path"`
	Http             httpConfig     `yaml:"http"`
	Postgres         postgresConfig `yaml:"postgres"`
	Redis            redisConfig    `yaml:"redis"`
}

type httpConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type postgresConfig struct {
	Host     string `yaml:"host" env:"DATABASE_HOST"`
	Port     string `yaml:"port" env:"DATABASE_PORT"`
	Name     string `yaml:"name" env:"DATABASE_NAME"`
	User     string `yaml:"user" env:"DATABASE_USER"`
	Password string `yaml:"password" env:"DATABASE_PASSWORD"`
}

type redisConfig struct {
	Host     string `yaml:"host" env:"REDIS_HOST"`
	Port     string `yaml:"port" env:"REDIS_PORT"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
}

func MustLoad() *Config {
	configPath := getConfigPath()
	if configPath == "" {
		log.Fatal("config path does not exist")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("config file does not exist in %s", configPath)
	}

	return &cfg
}

func getConfigPath() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv(configPathEnv)
	}

	return path
}
