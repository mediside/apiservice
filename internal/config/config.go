package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPathEnv = "CONFIG_PATH"

type Config struct {
	Env    string `yaml:"env" env:"ENV"`
	Domain string `yaml:"domain" env:"DOMAIN"`
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
