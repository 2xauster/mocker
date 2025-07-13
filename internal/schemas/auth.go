package schemas

type PublicUserSchema struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

type UserCreateRequest struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UserAuthenticateRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserAuthorizationRequest struct {
	AccessToken string `json:"access_token"`
}

type UserAuthenticateResponse struct {
	User PublicUserSchema `json:"user"`
	AccessToken string `json:"access_token"`
}