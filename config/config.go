package config

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App struct {
		LogLevel string `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
	} `yaml:"app"`
}

const (
	EnvConfigPathName  = "CONFIG_PATH"
	FlagConfigPathName = "config"
)

var configPath string
var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		// Определяем флаг, но парсинг будет в main.go
		flag.StringVar(&configPath, FlagConfigPathName, "configs/config.local.yaml", "path to config file")

		log.Print("config init")

		if path := os.Getenv(EnvConfigPathName); path != "" {
			configPath = path
		}

		instance = &Config{}

		if err := cleanenv.ReadConfig(configPath, instance); err != nil {
			helpText := "Consultant Parser"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
