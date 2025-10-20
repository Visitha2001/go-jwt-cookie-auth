package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/visitha2001/go-jwt-auth/configs"
	"github.com/visitha2001/go-jwt-auth/database"
	"github.com/visitha2001/go-jwt-auth/models"
	"gorm.io/gorm"
)

// --- DTOs (Data Transfer Objects) for Validation ---

type SignUpInput struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignInInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// --- Handlers ---

func SignUp(c *fiber.Ctx) error {
	// 1. Parse and Validate Input
	input := new(SignUpInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Error parsing JSON",
		})
	}
	// TODO: Add validation logic here (e.g., using 'go-playground/validator')

	// 2. Create User struct and Hash Password
	user := &models.User{
		Name:     input.Name,
		Username: input.Username,
		Email:    input.Email,
	}
	if err := user.HashPassword(input.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Could not hash password",
		})
	}

	// 3. Save to Database
	result := database.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status": "error", "message": "Email or Username already exists",
		})
	}

	// 4. Return Response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success", "message": "User created successfully",
	})
}

func SignIn(c *fiber.Ctx) error {
	// 1. Parse and Validate Input
	input := new(SignInInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Error parsing JSON",
		})
	}

	// 2. Find User by Email
	var user models.User
	result := database.DB.Where("email = ?", input.Email).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error", "message": "Invalid email or password",
		})
	}

	// 3. Check Password
	if err := user.CheckPassword(input.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error", "message": "Invalid email or password",
		})
	}

	// 4. Generate JWT
	token, err := configs.GenerateToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Could not generate token",
		})
	}

	// 5. Set HTTP-Only Cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 72), // Match JWT expiration
		HTTPOnly: true,                           // <-- Key for security
		Secure:   true,                           // <-- Set to true in production (HTTPS)
		SameSite: "Lax",
	}
	c.Cookie(&cookie)

	// 6. Return Response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Signed in successfully",
	})
}

func SignOut(c *fiber.Ctx) error {
	// Expire the cookie by setting its expiration to the past
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}
	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success", "message": "Signed out successfully",
	})
}

// GetProfile is a protected handler to get the current user's details
func GetProfile(c *fiber.Ctx) error {
	// 1. Get user ID from middleware (set in c.Locals)
	userID, err := configs.GetUserIDFromLocals(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error", "message": "Unauthorized",
		})
	}

	// 2. Find User by ID
	var user models.User
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error", "message": "User not found",
		})
	}

	// 3. Return user data (password is automatically omitted due to `json:"-"`)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   user,
	})
}
