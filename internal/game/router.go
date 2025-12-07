package game

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GameService interface {
	Create(ctx context.Context, input *CreateInput) (Game, error)
	GetAll(ctx context.Context) ([]Game, error)
}

type Router struct {
	service GameService
}

func NewRouter(service GameService) *Router {
	return &Router{
		service: service,
	}
}

func (r *Router) Track(c *gin.Context) {
	var input CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if _, err := r.service.Create(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.Status(http.StatusCreated)
}
