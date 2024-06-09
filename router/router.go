package router

import (
	"github.com/JerryJeager/chat-app-backend/internal/user"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userController *user.Controller){
	r = gin.Default()
	r.POST("/signup", userController.CreateUser)
	r.POST("/login", userController.LoginUser)
	r.GET("/logout", userController.Logout)
}

func Start(addr string) error{
	return r.Run(addr)
}