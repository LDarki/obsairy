package auth

import (
	"obsairy/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app *fiber.App, handler *Handler, repo Repository) {
	auth := app.Group("/auth")

	auth.Post("/register", func(c fiber.Ctx) error {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.Bind().Body(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		user, err := handler.Register(c.Context(), body.Email, body.Password)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(user)
	})

	auth.Post("/login", func(c fiber.Ctx) error {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.Bind().Body(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		userAgent := c.Get("User-Agent")
		ipAddress := c.IP()

		session, err := handler.Login(c.Context(), body.Email, body.Password, userAgent, ipAddress)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
		}

		return c.JSON(session)
	})

	// Private routes (protected by auth middleware)
	private := auth.Group("/sessions", middleware.JWTMiddleware(repo))

	private.Get("/me", func(c fiber.Ctx) error {
		userID := middleware.GetUserID(c)
		sessions, err := handler.GetSessions(c.Context(), userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(sessions)
	})

	private.Delete("/:sessionID", func(c fiber.Ctx) error {
		userID := middleware.GetUserID(c)
		sessionID := c.Params("sessionID")

		if err := handler.RevokeSession(c.Context(), userID, sessionID); err != nil {
			if err == ErrUnauthorized {
				return c.Status(403).JSON(fiber.Map{"error": "you can only revoke your own sessions"})
			}
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(204)
	})

	private.Delete("/all/me", func(c fiber.Ctx) error {
		userID := middleware.GetUserID(c)
		if err := handler.RevokeAllSessions(c.Context(), userID); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(204)
	})
}
