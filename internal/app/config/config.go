package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"
	"os"
)

type AppConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080/"`
	FileStorage   string `env:"FILE_STORAGE_PATH" envDefault:"./data/storage.dat"`
	Database_dsn  string `env:"DATABASE_DSN" envDefault:"127.0.0.1:5432"`
}

func (config *AppConfig) Init() error {
	fmt.Println(os.Args)
	if err := env.Parse(config); err != nil {
		fmt.Println("can't load service config", err)
		return err
	}

	pflag.StringVarP(&config.ServerAddress, "a", "a", config.ServerAddress, "Http-server address")
	pflag.StringVarP(&config.BaseURL, "b", "b", config.BaseURL, "Base URL")
	pflag.StringVarP(&config.FileStorage, "f", "f", config.FileStorage, "File storage path")
	pflag.StringVarP(&config.FileStorage, "d", "d", config.Database_dsn, "Database connection string")
	pflag.Parse()

	if config.BaseURL == "" || config.FileStorage == "" || config.ServerAddress == "" || config.Database_dsn == "" {
		if err := env.Parse(&config); err != nil {
			fmt.Println("can't load service config", err)
			return err
		}
	}
	if config.BaseURL[len(config.BaseURL)-1:] != "/" {
		config.BaseURL += "/"
	}
	return nil
}

func NewConfig() *AppConfig {
	return &AppConfig{}
}
