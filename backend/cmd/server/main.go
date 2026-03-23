package main

import (
	"log"
	"os"

	"obsairy/internal/auth"
	"obsairy/internal/config"
	"obsairy/internal/db"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	// Load env variables
	if err := config.LoadDotEnv(".env"); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	// DB
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	dbConn, err := db.InitPostgres(dsn)
	if err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}

	if err := db.RunMigrations(dbConn); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Auth module
	authRepo := auth.NewRepo()
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
	app.Use(logger.New())

	auth.RegisterRoutes(app, authHandler, authRepo)

	log.Fatal(app.Listen(":8080"))
}
