package auth

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/ashtonx86/mocker/internal/errs"
	"github.com/ashtonx86/mocker/internal/logging"
	"github.com/ashtonx86/mocker/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Include         []string
	Unauthorized    fiber.Handler
	DB              *sql.DB
	CompiledInclude []*regexp.Regexp
}

var ConfigDefault = Config{
	Include: []string{},
	Unauthorized: func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusUnauthorized).JSON(schemas.NewErrorAPIResponse(
			errors.New("unauthorized access"),
			"Unauthorized",
		))
	},
}

func pathToRegex(path string) (*regexp.Regexp, error) {
	escaped := regexp.QuoteMeta(path)
	pattern := "^" + strings.ReplaceAll(escaped, "\\*", ".*") + "$"
	return regexp.Compile(pattern)
}

func compileIncludes(patterns []string) []*regexp.Regexp {
	var regexes []*regexp.Regexp
	for _, p := range patterns {
		regex, err := pathToRegex(p)
		if err == nil {
			regexes = append(regexes, regex)
		}
	}
	return regexes
}

func configDefault(config ...Config) Config {
	if len(config) == 0 {
		return ConfigDefault
	}
	cfg := config[0]
	if cfg.Include == nil {
		cfg.Include = ConfigDefault.Include
	}
	if cfg.Unauthorized == nil {
		cfg.Unauthorized = ConfigDefault.Unauthorized
	}
	return cfg
}

func isPathIncluded(path string, includes []*regexp.Regexp) bool {
	for _, regex := range includes {
		if regex.MatchString(path) {
			return true
		}
	}
	return false
}

func New(config Config) fiber.Handler {
	cfg := configDefault(config)
	cfg.CompiledInclude = compileIncludes(cfg.Include)

	return func(c *fiber.Ctx) error {
		if len(cfg.Include) > 0 && !isPathIncluded(c.Path(), cfg.CompiledInclude) {
			return c.Next()
		}
		logging.Log(slog.LevelInfo, c, "Requiring authorization")

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return cfg.Unauthorized(c)
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return cfg.Unauthorized(c)
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		t, claims, err := VerifyJWT(tokenStr)
		var jwtErr errs.Error
		if err != nil && errors.As(err, &jwtErr) {
			return c.Status(fiber.StatusUnauthorized).JSON(schemas.NewErrorAPIResponse(err, "Unauthorized"))
		}
		if t == nil {
			return cfg.Unauthorized(c)
		}	
		if !t.Valid {
			return cfg.Unauthorized(c)
		}

		id, ok := claims["sub"].(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewErrorAPIResponse(
				errors.New("missing 'id' in token claims"),
				"Internal Error",
			))
		}

		ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
		defer cancel()

		user, err := GetUser(ctx, cfg.DB, schemas.UserFetchRequest{ID: id})
		if err != nil {
			var e errs.Error
			if errors.As(err, &e) && e.Code == errs.ErrNotFound {
				return c.Status(fiber.StatusNotFound).JSON(schemas.NewErrorAPIResponse(err, "User not found"))
			}
			return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewErrorAPIResponse(err, "Failed to fetch user"))
		}

		c.Locals("user", &user)
		return c.Next()
	}
}