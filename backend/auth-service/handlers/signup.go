package handlers

import (
	"database/sql"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"food-recipes/backend/auth-service/utils"
)

type SignupInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUp(c *fiber.Ctx) error {
	input := new(SignupInput)
	if err := c.BodyParser(input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error hashing password")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "DB error")
	}
	defer db.Close()

	id := uuid.New()
	_, err = db.Exec(`INSERT INTO users(id, name, email, password) VALUES($1, $2, $3, $4)`,
		id, input.Name, input.Email, hashed)
	if err != nil {
		return fiber.NewError(fiber.StatusConflict, "Email already exists")
	}

	token, err := utils.GenerateJWT(id.String())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Token generation error")
	}

	return c.JSON(fiber.Map{
		"user_id": id,
		"token":   token,
	})
}
