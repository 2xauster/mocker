package v1

import (
	"github.com/ashtonx86/mocker/internal/supervisor"
	"github.com/gofiber/fiber/v2"
)

func MapV1Handlers(app *fiber.App, su *supervisor.Supervisor) {
	router := app.Group("/api/v1")

	authHandler := NewAuthHandler(su)
	authHandler.MapRoutes(router.Group("/auth").(*fiber.Group))
}