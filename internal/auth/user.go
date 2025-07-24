package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ashtonx86/mocker/internal/data"
	"github.com/ashtonx86/mocker/internal/entities"
	"github.com/ashtonx86/mocker/internal/errs"
	"github.com/ashtonx86/mocker/internal/schemas"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(ctx context.Context, db *sql.DB, userData schemas.UserCreateRequest) (*entities.User, error) {
	if userData.Password != userData.ConfirmPassword {
		return nil, errs.NewError(errors.New("data mismatch: Password != ConfirmPassword"), errs.DataErrorType, errs.ErrDataMismatch)
	}	

	// bcrypt.DefaultCost = 10
	bytes, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.NewError(fmt.Errorf("failed to hash password :: %w", err), errs.DataErrorType, errs.ErrInternalFailure)
	}

	user := entities.User{
		ID:           uuid.NewString(),
		Name:         userData.Name,
		Email:        userData.Email,
		PasswordHash: string(bytes),

		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
	}
	_, err = data.Insert(ctx, db, user)

	if err != nil && errors.Is(err, errs.Error{Code: errs.ErrAlreadyExists, Type: errs.SQLErrorType.String()}) {

		return nil, errs.NewError(fmt.Errorf("data already exists :: %w", err), errs.DataErrorType, errs.ErrAlreadyExists)
	} else if err != nil {

		return nil, errs.NewError(fmt.Errorf("[pkg auth func CreateUser] failed to insert data :: %w", err), errs.DataErrorType, errs.ErrInternalFailure)
	}

	return &user, nil
}
