package main

import (
	awshelpers "cloud-posture-scanner/aws-helpers"
	_ "cloud-posture-scanner/docs"
	"cloud-posture-scanner/routers"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to environment variables")
	}

	router := routers.InitRouters()
	awshelpers.InitS3Config()

	log.Println("Cloud Posture Scanner is running")
	router.Run(":8080")
}
