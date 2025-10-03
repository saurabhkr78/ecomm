package configs

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

// DSN-Data Source name
type AppConfig struct {
	ServerPort string
	Dsn        string
	AppSecret  string
}

// here we will be reading env file and setting up the APPConfig struct
func SetUpEnv() (cfg AppConfig, err error) {

	//some condition to check env file

	if os.Getenv("APP_ENV") == "dev" {
		// Load development-specific configuration
		godotenv.Load()
	}
	//read the env variable
	httpPort := os.Getenv("HTTP_PORT")

	//check if httpPort is empty
	if len(httpPort) < 1 {
		return AppConfig{}, errors.New("env variable HTTP_PORT not found")
	}
	//dsn for database connection
	Dsn := os.Getenv("DSN")
	if len(Dsn) < 1 {
		return AppConfig{}, errors.New("env variable DSN not found")
	}
	/*
	   AppConfig{} creates an empty struct of type AppConfig with all fields set to their zero values.
	   errors.New(...) creates a new error object with the given message.
	*/

	AppSecret := os.Getenv("APP_SECRET")

	if len(AppSecret) < 1 {
		return AppConfig{}, errors.New("env variable APP_SECRET not found")
	}

	return AppConfig{ServerPort: httpPort, Dsn: Dsn, AppSecret: AppSecret}, nil
}
