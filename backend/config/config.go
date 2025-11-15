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
	Port     string `json:"port" envconfig:"PORT"`
	Host     string `json:"host" envconfig:"SERVER_HOST"`
	BasePath string `json:"BasePath" envconfig:"BASE_PATH"`
}
type Database struct {
	ConnectionString string `json:"connectionString" envconfig:"DATABASE_CONNECTION_STRING"`
}
type Authentication struct {
	ClientSecret      string `json:"clientSecret" envconfig:"CLIENT_SECRET"`
	ClientId          string `json:"clientId" envconfig:"CLIENT_ID"`
	ClientOrigin      string `json:"clientOrigin" envconfig:"CLIENT_ORIGIN"`
	ClientRedirectUri string `json:"clientRedirectUri" envconfig:"CLIENT_REDIRECT_URI"`
	EncryptionKey     string `json:"encryptionKey" envconfig:"ENCRYPTION_KEY"`
	JwtKey            string `json:"jwtKey" envconfig:"JWT_KEY"`
	ScopesField       string `json:"scopes" envconfig:"CLIENT_SCOPES"`
	AudienceField     string `json:"audience" envconfig:"CLIENT_AUDIENCE"`
	Issuer            string `json:"issuer" envconfig:"CLIENT_ISSUER"`
}
type Configuration struct {
	ServerField   Server         `json:"server" envconfig:"SERVER"`
	DatabaseField Database       `json:"database" envconfig:"DATABASE"`
	Environment   string         `json:"environment" envconfig:"ENVIRONMENT"`
	AuthField     Authentication `json:"authentication" envconfig:"AUTHENTICATION"`
}

func (conf *Configuration) Server() Server {
	return conf.ServerField
}
func (conf *Configuration) Database() Database {
	return conf.DatabaseField
}
func (conf *Configuration) Authentication() Authentication {
	return conf.AuthField
}
func (conf *Configuration) Scopes() string {
	return conf.AuthField.ScopesField
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
