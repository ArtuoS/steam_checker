package user

import "github.com/go-playground/validator/v10"

type CreateInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateOutput struct {
	Token string `json:"token"`
}

type TrackInput struct {
	AppID int `json:"app_id" uri:"app_id"`
}

func (c *CreateInput) Validate() error {
	return validator.New().Struct(c)
}

func (c *AuthenticateInput) Validate() error {
	return validator.New().Struct(c)
}

func (c *TrackInput) Validate() error {
	return validator.New().Struct(c)
}
