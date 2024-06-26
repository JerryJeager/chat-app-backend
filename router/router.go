package router

import (
	"net/http"

	"github.com/JerryJeager/chat-app-backend/internal/user"
	"github.com/JerryJeager/chat-app-backend/internal/ws"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userController *user.Controller, wsController *ws.Controller) {
	r = gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "chat app v1 server"})
	})
	r.POST("/signup", userController.CreateUser)
	r.POST("/login", userController.LoginUser)
	r.GET("/logout", userController.Logout)



	r.POST("/ws/createRoom", wsController.CreateRoom)
	r.GET("/ws/joinRoom/:roomId", wsController.JoinRoom)
}

func Start(addr string) error {
	return r.Run(addr)
}
