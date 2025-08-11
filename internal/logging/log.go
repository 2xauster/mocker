package logging

import (
	"context"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func Log(level slog.Level, c *fiber.Ctx, msg string, args ...any) {
	var apiPath string
	var reqMethod string

	if c != nil {
		apiPath = c.Path()
		reqMethod = c.Method()
	}

	attrs := []any{
		slog.String("apiPath", apiPath),
		slog.String("reqMethod", reqMethod),
	}

	attrs = append(attrs, args...)

	logger.Log(context.Background(), level, msg, attrs...)
}
