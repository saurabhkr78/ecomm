package api

import (
	"ecomm/configs"
	"ecomm/internal/api/rest"
	"ecomm/internal/api/rest/handlers"
	"ecomm/internal/domain"
	"ecomm/internal/helper"
	"log"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(config configs.AppConfig) {
	//create a new fiber app
	app := fiber.New()

	//connect go orm here
	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v\n", err)
	}

	log.Println("Database connected successfully")

	//if database connection successful then runthe migration(here auto migration automatically detect the changes in user.go domain file and create table accordingly)
	db.AutoMigrate(&domain.User{})

	//befor resthandler we gonna create a auth instance
	//so that we can use this auth instance in user service
	//auth instance will have the secret key which we will use to generate token and hash password
	//we can create a new field in rest handler struct to hold this auth instance but instead we can directly pass this auth instance to user service while creating its instance in user handler
	//because only user service need this auth instance

	auth := helper.SetUpAuth(config.AppSecret)

	//intantiate rest handler
	rh := &rest.RestHandler{
		App:  app,
		DB:   db,
		Auth: auth,
	}
	//to create table we need migration also using gorm

	//connect the databse for ORM after start of the server

	SetupRoutes(rh)
	//start the server
	app.Listen(config.ServerPort)
}

// function to setup all the routes
func SetupRoutes(rh *rest.RestHandler) {
	//setup user routes
	handlers.SetupUserRoutes(rh)
	//transaction routes
	//catalog routes
}
