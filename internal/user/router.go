package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	Create(ctx context.Context, input *CreateInput) error
	Authenticate(ctx context.Context, input *AuthenticateInput) (AuthenticateOutput, error)
}

type Router struct {
	service UserService
}

func NewRouter(service UserService) *Router {
	return &Router{
		service: service,
	}
}

func (r *Router) Create(c *gin.Context) {
	var input CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := r.service.Create(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.Status(http.StatusCreated)
}

func (r *Router) Authenticate(c *gin.Context) {
	var input AuthenticateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	output, err := r.service.Authenticate(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, output)
}
