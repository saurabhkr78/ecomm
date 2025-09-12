package main

import (
	"ecomm/configs"
	"ecomm/internal/api"
	"log"
)

func main() {
	//check and load env file
	cfg, err := configs.SetUpEnv()
	if err != nil {
		log.Fatalf("Error setting up environment: %v\n", err) //exist the program  if error occurs,it's os level printf
	}
	//if everything is fine start the server

	api.StartServer(cfg)
}
