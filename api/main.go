package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/auth"
	"github.com/myhro/ovh-checker/api/hardware"
	"github.com/myhro/ovh-checker/storage"
)

func main() {
	r := gin.Default()
	port := ":8080"

	cache, err := storage.NewCache()
	if err != nil {
		log.Fatal("cache: ", err)
	}

	db, err := storage.NewDB()
	if err != nil {
		log.Fatal("database: ", err)
	}

	authHandler, err := auth.NewHandler(cache, db)
	if err != nil {
		log.Fatal("auth: ", err)
	}
	r.GET("/auth/tokens", authHandler.AuthRequired, authHandler.Tokens)
	r.GET("/auth/user", authHandler.AuthRequired, authHandler.User)
	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/logout", authHandler.AuthRequired, authHandler.Logout)
	r.POST("/auth/signup", authHandler.Signup)

	hardwareHandler, err := hardware.NewHandler(db)
	if err != nil {
		log.Fatal("hardware: ", err)
	}
	r.GET("/hardware/offers", hardwareHandler.Offers)

	if gin.Mode() == gin.ReleaseMode {
		log.Print("Starting server on port ", port)
	}
	r.Run(port)
}
