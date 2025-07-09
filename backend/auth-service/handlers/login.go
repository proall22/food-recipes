package handlers

import (
	"database/sql"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"

	"food-recipes/backend/auth-service/utils"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	input := new(LoginInput)
	if err := c.BodyParser(input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "DB error")
	}
	defer db.Close()

	var id string
	var hashed string
	err = db.QueryRow(`SELECT id, password FROM users WHERE email = $1`, input.Email).Scan(&id, &hashed)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "User not found")
	}

	if !utils.CheckPasswordHash(input.Password, hashed) {
		return fiber.NewError(fiber.StatusUnauthorized, "Wrong password")
	}

	token, err := utils.GenerateJWT(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Token error")
	}

	return c.JSON(fiber.Map{
		"user_id": id,
		"token":   token,
	})
}
