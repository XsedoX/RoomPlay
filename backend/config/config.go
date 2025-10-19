package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
)

const (
	envPrefix      = ""
	envDevelopment = "development"
	envProduction  = "production"
)

type Configuration struct {
	Server struct {
		Port string `json:"port" envconfig:"SERVER_PORT"`
		Host string `json:"host" envconfig:"SERVER_HOST"`
	} `json:"server"`
	Database struct {
		ConnectionString string `json:"connectionString" envconfig:"DATABASE_CONNECTION_STRING"`
	} `json:"database"`
	Environment    string `json:"environment" envconfig:"ENVIRONMENT"`
	Authentication struct {
		JwtSecret string `json:"jwtSecret" envconfig:"JWT_SECRET"`
	}
}

func (conf *Configuration) IsDevelopment() bool {
	return conf.Environment == envDevelopment
}
func (conf *Configuration) IsProduction() bool {
	return conf.Environment == envProduction
}

func Load() *Configuration {
	var config Configuration
	readFile(&config)
	loadEnv(&config)
	return &config
}
func loadEnv(cfg *Configuration) {
	err := envconfig.Process(envPrefix, cfg)
	if err != nil {
		processError(err)
	}
}
func readFile(cfg *Configuration) {
	f, err := os.Open("./config/config.json")
	if err != nil {
		processError(err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}
func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
