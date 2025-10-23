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

type Server struct {
	Port string `json:"port" envconfig:"SERVER_PORT"`
	Host string `json:"host" envconfig:"SERVER_HOST"`
}
type Database struct {
	ConnectionString string `json:"connectionString" envconfig:"DATABASE_CONNECTION_STRING"`
}
type Authentication struct {
	ClientSecret  string `json:"clientSecret" envconfig:"CLIENT_SECRET"`
	ClientId      string `json:"clientId" envconfig:"CLIENT_ID"`
	ClientOrigin  string `json:"clientOrigin" envconfig:"CLIENT_ORIGIN"`
	EncryptionKey string `json:"encryptionKey" envconfig:"ENCRYPTION_KEY"`
	JwtKey        string `json:"jwtKey" envconfig:"JWT_KEY"`
}
type Configuration struct {
	Server         Server         `json:"server" envconfig:"SERVER"`
	Database       Database       `json:"database" envconfig:"DATABASE"`
	Environment    string         `json:"environment" envconfig:"ENVIRONMENT"`
	Authentication Authentication `json:"authentication" envconfig:"AUTHENTICATION"`
	Scopes         string         `json:"scopes" envconfig:"SCOPES"`
}

var Config *Configuration

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
