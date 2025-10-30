package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// DSN-Data Source name
type AppConfig struct {
	ServerPort        string
	Dsn               string
	AppSecret         string
	TwilioSid         string
	TwilioAuthToken   string
	TwilioPhoneNumber string
	StripeSecretKey   string
}

// here we will be reading env file and setting up the APPConfig struct
func SetUpEnv() (cfg AppConfig, err error) {

	//some condition to check env file

	if os.Getenv("APP_ENV") == "dev" {
		// Load development-specific configuration
		godotenv.Load()
	}
	//read the env variable
	httpPort := os.Getenv("SERVER_PORT")
	//dsn for database connection
	Dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v TimeZone=Asia/Kolkata",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	/*
	   AppConfig{} creates an empty struct of type AppConfig with all fields set to their zero values.
	   errors.New(...) creates a new error object with the given message.
	*/
	AppSecret := os.Getenv("APP_SECRET")

	return AppConfig{
		ServerPort:        httpPort,
		Dsn:               Dsn,
		AppSecret:         AppSecret,
		TwilioSid:         os.Getenv("TWILIO_SID"),
		TwilioAuthToken:   os.Getenv("TWILIO_AUTH_TOKEN"),
		TwilioPhoneNumber: os.Getenv("TWILIO_PHONE_NUMBER"),
		StripeSecretKey:   os.Getenv("STRIPE_SECRET_KEY"),
	}, nil
}
