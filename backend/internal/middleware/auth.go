package middleware

import (
	"context" // Added context import
	"errors"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// AuthRepo interface to interact with session storage
type AuthRepo interface {
	FindSessionByToken(ctx context.Context, token string) (any, error)
}

const UserIDKey = "userID"
const SessionKey = "session"

func JWTMiddleware(repo AuthRepo) fiber.Handler {
	return func(c fiber.Ctx) error {
		tokenString, err := extractToken(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		// Verify token exists in DB
		session, err := repo.FindSessionByToken(c.Context(), tokenString)
		if err != nil || session == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "session not found or revoked"})
		}

		userID, err := validateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}

		c.Locals(UserIDKey, userID)
		c.Locals(SessionKey, session)
		return c.Next()
	}
}

// GetUserID retrieves the user ID from the fiber context.
func GetUserID(c fiber.Ctx) string {
	val, ok := c.Locals(UserIDKey).(string)
	if !ok {
		return ""
	}
	return val
}

// GetSession retrieves the session from the fiber context.
func GetSession(c fiber.Ctx) any {
	return c.Locals(SessionKey)
}

func extractToken(c fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}
	}

	// Try cookie
	token := c.Cookies("auth_token")
	if token != "" {
		return token, nil
	}

	return "", errors.New("missing authorization header or session cookie")
}

func validateToken(tokenString string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "very-secret-key"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("user id not found in token")
	}

	return userID, nil
}
