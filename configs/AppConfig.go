package configs

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type AppConfig struct {
	ServerPort string
}

// here we will be reading env file and setting up the APPConfig struct
func SetUpEnv() (cfg AppConfig, err error) {

	//some condition to check env file

	if os.Getenv("ENV_FILE") == "dev" {
		// Load development-specific configuration
		godotenv.Load()
	}
	//read the env variable
	httpPort := os.Getenv("HTTP_PORT")

	//check if httpPort is empty
	if len(httpPort) < 1 {
		return AppConfig{}, errors.New("env_variable HTTP_PORT not set")
	}
	/*
	   AppConfig{} creates an empty struct of type AppConfig with all fields set to their zero values.
	   errors.New(...) creates a new error object with the given message.
	*/
	return AppConfig{ServerPort: httpPort}, nil
}
