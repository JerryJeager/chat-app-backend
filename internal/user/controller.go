package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service
}

func NewController(s Service) *Controller {
	return &Controller{
		Service: s,
	}
}

func (c *Controller) CreateUser(ctx *gin.Context) {
	var u CreateUserReq
	if err := ctx.ShouldBindJSON(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.Service.CreateUser(ctx.Request.Context(), &u)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, *user)
}
