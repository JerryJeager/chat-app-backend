package main

import (
	"log"
	"os"

	"github.com/JerryJeager/chat-app-backend/config"
	"github.com/JerryJeager/chat-app-backend/db"
	"github.com/JerryJeager/chat-app-backend/internal/user"
	"github.com/JerryJeager/chat-app-backend/internal/ws"
	"github.com/JerryJeager/chat-app-backend/router"
)

func main() {
	config.LoadEnv()
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userController := user.NewController(userSvc)

	hub := ws.NewHub()
	wsController := ws.NewController(hub)

	go hub.Run()

	router.InitRouter(userController, wsController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Start(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
