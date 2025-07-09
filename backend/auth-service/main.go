package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"os"

	"github.com/joho/godotenv"
	"food-recipes/backend/auth-service/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	app := fiber.New()

	app.Post("/signup", handlers.SignUp)
	app.Post("/login", handlers.Login)

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}
