package api

import (
	"ecomm/configs"
	"ecomm/internal/api/rest"
	"ecomm/internal/api/rest/handlers"
	"ecomm/internal/domain"
	"ecomm/internal/helper"
	"ecomm/pkg/payment"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
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
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Address{},
		&domain.BankAccount{},
		&domain.Category{},
		&domain.Product{},
		&domain.Cart{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Payment{},
	)
	if err != nil {
		log.Fatalf("Error on running migration %v", err.Error())
	}

	log.Println("Migration completed successfully")

	//cors configuration shoould be added here
	// c := cors.New(cors.Config{
	// 	AllowOrigins:     "http://localhost:3030,http://localhost:3000",
	// 	AllowHeaders:     "Content-Type,Accept,Authorization",
	// 	AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	// 	AllowCredentials: true,
	// })
	// app.Use(c)
	// Define allowed origins
	allowedOrigins := map[string]bool{
		"http://localhost:3030": true,
		"http://localhost:3000": true,
	}

	// Configure CORS middleware
	c := cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return allowedOrigins[origin]
		},
		AllowHeaders: []string{
			"Content-Type",
			"Accept",
			"Authorization",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowCredentials: true, //if using cookies or auth tokens
	})

	// Use CORS middleware
	app.Use(c)

	//befor resthandler we gonna create a auth instance
	//so that we can use this auth instance in user service
	//auth instance will have the secret key which we will use to generate token and hash password
	//we can create a new field in rest handler struct to hold this auth instance but instead we can directly pass this auth instance to user service while creating its instance in user handler
	//because only user service need this auth instance

	auth := helper.SetUpAuth(config.AppSecret)

	log.Printf("Stripe Config: Key=%s, SuccessURL=%s, CancelURL=%s",
		config.StripeSecretKey,
		config.SuccessUrl,
		config.CancelUrl,
	)

	paymentClient := payment.NewPaymentClient(config.StripeSecretKey, config.SuccessUrl, config.CancelUrl)
	//intantiate rest handler
	rh := &rest.RestHandler{
		App:    app,
		DB:     db,
		Auth:   auth,
		Config: config,
		Pc:     paymentClient,
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
	handlers.SetupTransactionRoutes(rh)
	//catalog routes
	handlers.SetupCatalogRoutes(rh)
}
