package schemas

import "time"

type PublicUserSchema struct {
	ID    string `json:"id"`
	Email string `json:"email,omitempty"`
	Name  string `json:"name"`

	CreatedAt time.Time `json:"created_at"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

type UserCreateRequest struct {
	Name            string `json:"name" validate:"required,min=1,max=45"`
	Email           string `json:"email" validate:"required,min=1,email"`
	Password        string `json:"password" validate:"required,min=14"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=14"`
}

type UserFetchRequest struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty" validate:"email"`
}

type UserAuthenticateRequest struct {
	Email    string `json:"email" validate:"required,min=1,email"`
	Password string `json:"password" validate:"required,min=14"`
}

type UserAuthorizationRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type UserAuthenticateResponse struct {
	User        PublicUserSchema `json:"user"`
	AccessToken string           `json:"access_token"`
}