package config

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

// env-default:"production" -> should be used in production
type Config struct { // struct tags
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	Storagepath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

// should not return error, must directly end the program if error occurs
func MustLoad() *Config {
	slog.Info("Loading config started")
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	// Takes config from terminal
	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			slog.Error("loading config path")
			log.Fatal("Config path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not esist: %s", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		slog.Error("reading the config file")
		log.Fatalf("can not read config file: %s", err.Error())
	}
	slog.Info("Successfully loaded the config")
	return &cfg
}
