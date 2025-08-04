package v1

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/ashtonx86/mocker/internal/auth"
	"github.com/ashtonx86/mocker/internal/data"
	"github.com/ashtonx86/mocker/internal/errs"
	"github.com/ashtonx86/mocker/internal/schemas"
	"github.com/ashtonx86/mocker/internal/supervisor"
	"github.com/gofiber/fiber/v2"
)

const TIMEOUT = 40 * time.Second

// assert: AuthHandler implements Handler interface.
var _ Handler = (*AuthHandler)(nil)

type AuthHandler struct {
	Supervisor *supervisor.Supervisor

	SQLite *data.SQLite
}

func NewAuthHandler(su *supervisor.Supervisor) *AuthHandler {
	return &AuthHandler{
		Supervisor: su,

		SQLite: su.SQLite,
	}
}

func (h *AuthHandler) MapRoutes(router *fiber.Group) {
	router.Get("/", h.handleGET)
	router.Post("/", h.handlePOST)

	router.Post("/access-token", h.handleGETAccessToken)
}

func (h *AuthHandler) handleGET(c *fiber.Ctx) error {
	req := new(schemas.UserFetchRequest)
	c.BodyParser(&req)

	err := errs.Validate(req)

	if err != nil && errors.Is(err, errs.Error{Code: errs.ErrDataIllegal, Type: errs.DataErrorType.String()}) {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.NewErrorAPIResponse(err, "Bad request"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	user, err := auth.GetUser(ctx, h.SQLite.DB, *req)
	if err != nil {
		var e errs.Error
		if errors.As(err, &e) {
			slog.Error("[pkg handlers : auth.go : func handleGET] User creation failed ", "error", err)

			switch e.Code {
			case errs.ErrAlreadyExists:
				return c.Status(fiber.StatusForbidden).JSON(schemas.NewErrorAPIResponse(err, "Already exists"))
			case errs.ErrDataMismatch:
				return c.Status(fiber.StatusBadRequest).JSON(schemas.NewErrorAPIResponse(err, "Data mismatch"))
			case errs.ErrInternalFailure:
				return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewErrorAPIResponse(err, "Internal failure"))
			case errs.ErrNotFound:
				return c.Status(fiber.StatusNotFound).JSON(schemas.NewErrorAPIResponse(err, "Not found"))
			}
		}
	}
	return c.JSON(schemas.NewAPIResponse(true, schemas.PublicUserSchema{
		Name: user.Name,
		ID: user.ID,
	}, ""))
}

func (h *AuthHandler) handlePOST(c *fiber.Ctx) error {
	req := new(schemas.UserCreateRequest)
	c.BodyParser(&req)

	err := errs.Validate(req)

	if err != nil && errors.Is(err, errs.Error{Code: errs.ErrDataIllegal, Type: errs.DataErrorType.String()}) {
		return c.Status(403).JSON(schemas.NewErrorAPIResponse(err, "Bad request"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	user, err := auth.CreateUser(ctx, h.SQLite.DB, *req)
	if err != nil {
		var e errs.Error
		if errors.As(err, &e) {
			slog.Error("[pkg handlers : auth.go : func handlePOST] User creation failed ", "error", err)

			switch e.Code {
			case errs.ErrAlreadyExists:
				return c.Status(fiber.StatusForbidden).JSON(schemas.NewErrorAPIResponse(err, "Already exists"))
			case errs.ErrDataMismatch:
				return c.Status(fiber.StatusBadRequest).JSON(schemas.NewErrorAPIResponse(err, "Data mismatch"))
			case errs.ErrInternalFailure:
				return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewErrorAPIResponse(err, "Internal failure"))
			}
		}
	}
	return c.JSON(schemas.NewAPIResponse(true, schemas.PublicUserSchema{
		ID: user.ID,
		Name: user.Name,
	}, ""))
}

func (h *AuthHandler) handleGETAccessToken(c *fiber.Ctx) error {
	req := new(schemas.UserAuthenticateRequest)
	c.BodyParser(&req)

	err := errs.Validate(req)

	if err != nil && errors.Is(err, errs.Error{Code: errs.ErrDataIllegal, Type: errs.DataErrorType.String()}) {
		return c.Status(403).JSON(schemas.NewErrorAPIResponse(err, "Bad request"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	user, err := auth.GetUser(ctx, h.SQLite.DB, schemas.UserFetchRequest{
		Email: req.Email,
	})
	if err != nil {
		var e errs.Error
		if errors.As(err, &e) {
			slog.Error("[pkg handlers : auth.go : func handlePOST] User creation failed ", "error", err)

			switch e.Code {
			case errs.ErrAlreadyExists:
				return c.Status(fiber.StatusForbidden).JSON(schemas.NewErrorAPIResponse(err, "Already exists"))
			case errs.ErrDataMismatch:
				return c.Status(fiber.StatusBadRequest).JSON(schemas.NewErrorAPIResponse(err, "Data mismatch"))
			case errs.ErrInternalFailure:
				return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewErrorAPIResponse(err, "Internal failure"))
			}
		}
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(schemas.NewErrorAPIResponse(err, "Internal server error"))
	}

	res := schemas.UserAuthenticateResponse{
		User: schemas.PublicUserSchema{
			ID: user.ID,
			Email: user.Email,

			CreatedAt: user.CreatedAt,
			LastUpdatedAt: user.LastUpdatedAt,
		},
		AccessToken: token,
	}
	return c.JSON(schemas.NewAPIResponse(true, res, ""))
}