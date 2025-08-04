package server

import (
	"context"
	"log/slog"

	"github.com/ashtonx86/mocker/internal/auth"
	v1 "github.com/ashtonx86/mocker/internal/handlers/v1"
	"github.com/ashtonx86/mocker/internal/supervisor"
	"github.com/gofiber/fiber/v2"
)

var ProtectedRoutes = []string{
	"/api/v1/auth/protected",
}

type WebServer struct {
	App        *fiber.App
	Supervisor *supervisor.Supervisor
}

func NewWebServer() *WebServer {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := fiber.New()

	su, err := supervisor.New(ctx)
	if err != nil {
		slog.Error("[func NewWebServer] init failed ", "error", err)
		return nil
	}

	su.Init()

	authMiddleware := auth.New(auth.Config{
		DB:      su.SQLite.DB,
		Include: ProtectedRoutes,
	})
	app.Use(authMiddleware)

	v1.MapV1Handlers(app, su)

	return &WebServer{
		App:        app,
		Supervisor: su,
	}
}
