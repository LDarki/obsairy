package auth

import (
	"obsairy/internal/middleware"
	"time"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app *fiber.App, handler *Handler, repo Repository) {
	authGroup := app.Group("/auth")
	registerPublicRoutes(authGroup, handler)
	registerPrivateRoutes(authGroup, handler, repo)
}

func registerPublicRoutes(group fiber.Router, handler *Handler) {
	group.Post("/register", func(c fiber.Ctx) error {
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

		// Auto login after registration
		userAgent := c.Get("User-Agent")
		ipAddress := c.IP()
		session, _, err := handler.Login(c.Context(), body.Email, body.Password, userAgent, ipAddress)
		if err != nil {
			return c.JSON(user)
		}

		setAuthCookie(c, session.Token, session.ExpiresAt)
		return c.JSON(user)
	})

	group.Post("/login", func(c fiber.Ctx) error {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.Bind().Body(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		userAgent := c.Get("User-Agent")
		ipAddress := c.IP()

		session, user, err := handler.Login(c.Context(), body.Email, body.Password, userAgent, ipAddress)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
		}

		setAuthCookie(c, session.Token, session.ExpiresAt)
		return c.JSON(user)
	})
}

func registerPrivateRoutes(group fiber.Router, handler *Handler, repo Repository) {
	private := group.Group("/sessions", middleware.JWTMiddleware(repo))

	private.Get("/me", func(c fiber.Ctx) error {
		userID := middleware.GetUserID(c)
		user, err := handler.GetUser(c.Context(), userID)
		if err != nil || user == nil {
			return c.Status(404).JSON(fiber.Map{"error": "user not found"})
		}
		return c.JSON(user)
	})

	private.Get("/", func(c fiber.Ctx) error {
		userID := middleware.GetUserID(c)
		sessions, err := handler.GetSessions(c.Context(), userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(sessions)
	})

	registerDeletionRoutes(private, handler)
}

func registerDeletionRoutes(router fiber.Router, handler *Handler) {
	router.Delete("/me", func(c fiber.Ctx) error {
		userID := middleware.GetUserID(c)
		sessAny := middleware.GetSession(c)

		if sess, ok := sessAny.(*Session); ok {
			_ = handler.RevokeSession(c.Context(), userID, sess.ID)
		}

		c.ClearCookie("auth_token")
		return c.SendStatus(204)
	})

	router.Delete("/:sessionID", func(c fiber.Ctx) error {
		userID := middleware.GetUserID(c)
		sessionID := c.Params("sessionID")

		if err := handler.RevokeSession(c.Context(), userID, sessionID); err != nil {
			status := 500
			if err == ErrUnauthorized {
				status = 403
			}
			return c.Status(status).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(204)
	})

	router.Delete("/all/me", func(c fiber.Ctx) error {
		userID := middleware.GetUserID(c)
		_ = handler.RevokeAllSessions(c.Context(), userID)
		c.ClearCookie("auth_token")
		return c.SendStatus(204)
	})
}

func setAuthCookie(c fiber.Ctx, token string, expiresAt time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  expiresAt,
		HTTPOnly: true,
		SameSite: "Lax",
	})
}
