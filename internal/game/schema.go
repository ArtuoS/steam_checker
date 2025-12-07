package game

import (
	"github.com/go-playground/validator/v10"
)

type CreateInput struct {
	AppID int `json:"app_id" validate:"required"`
}

func (c *CreateInput) Validate() error {
	return validator.New().Struct(c)
}
