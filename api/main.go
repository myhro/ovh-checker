package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/auth"
	"github.com/myhro/ovh-checker/api/hardware"
)

func main() {
	r := gin.Default()
	port := ":8080"

	authHandler, err := auth.NewHandler()
	if err != nil {
		log.Fatal("auth: ", err)
	}
	r.POST("/auth/signup", authHandler.Signup)

	hardwareHandler, err := hardware.NewHandler()
	if err != nil {
		log.Fatal("hardware: ", err)
	}
	r.GET("/hardware/offers", hardwareHandler.Offers)

	if gin.Mode() == gin.ReleaseMode {
		log.Print("Starting server on port ", port)
	}
	r.Run(port)
}
