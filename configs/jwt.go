package configs

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Claims defines the JWT claims
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT for a given user ID
func GenerateToken(userID uint) (string, error) {
	jwtSecret := EnvConfig("JWT_SECRET")
	expiresInHours, _ := strconv.Atoi(EnvConfig("JWT_EXPIRES_IN_HOURS"))

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expiresInHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken parses and validates a JWT string
func ValidateToken(tokenString string) (*Claims, error) {
	jwtSecret := EnvConfig("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

// GetUserIDFromLocals retrieves the user ID set by the auth middleware
func GetUserIDFromLocals(c *fiber.Ctx) (uint, error) {
	id, ok := c.Locals("userID").(uint)
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Cannot parse userID from context")
	}
	return id, nil
}
