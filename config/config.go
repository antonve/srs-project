package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"

	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/labstack/echo.v1/middleware"
)

// Environment in which the application runs
type Environment string

// Environments contains all the possible environments
var Environments = map[string]Environment{
	"dev":  "dev.yml",
	"prod": "prod.yml",
	"test": "test.yml",
}

// Config contains the necessary application configuration
type Config struct {
	ConnectionString string `yaml:"connection_string"`
	Database         string `yaml:"database"`
	Debug            bool   `yaml:"debug"`
	Environment      Environment
	JWTKey           string `yaml:"JWT_key"`
	MigrationsPath   string `yaml:"migrations_path"`
}

// GetCompleteConnectionString returns the connection string based on the current config
func (config Config) GetCompleteConnectionString() string {
	return fmt.Sprintf("%s%s", config.ConnectionString, config.Database)
}

// Local static variables
var config = Config{}
var environment = Environments["dev"]

// SetEnviroment sets the application wide environment
func SetEnviroment(env Environment) {
	environment = env
}

// GetConfig returns the application wide configuration
func GetConfig() Config {
	// Return config if we have already loaded it
	if (Config{}) != config {
		return config
	}

	// Load config file
	configData, err := ioutil.ReadFile(fmt.Sprintf("%s/src/srs-project/app/config/%s", os.Getenv("GOPATH"), environment))
	if err != nil {
		log.Fatalf("Could not load config for environment `%s`", environment)
	}

	// Parse config file
	config = Config{
		Environment: environment,
	}
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf("Could not parse config for environment `%s`", environment)
	}

	return config
}

// GetJWTConfig returns the JWT config
func GetJWTConfig(claims jwt.Claims) middleware.JWTConfig {
	return middleware.JWTConfig{
		Claims:     claims,
		SigningKey: []byte(GetConfig().JWTKey),
		ContextKey: "user",
	}
}
