package main

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/auth"
	"github.com/myhro/ovh-checker/api/hardware"
	"github.com/myhro/ovh-checker/storage"
)

// API is the main server structure
type API struct {
	router *gin.Engine
	port   string

	cache storage.Cache
	db    storage.DB
	store storage.CookieStore
}

func main() {
	cache, err := storage.NewCache()
	if err != nil {
		log.Fatal("cache: ", err)
	}

	db, err := storage.NewDB()
	if err != nil {
		log.Fatal("database: ", err)
	}

	store, err := storage.NewCookieStore()
	if err != nil {
		log.Print("store: ", err)
	}

	api := API{
		router: gin.Default(),
		port:   ":8080",
		cache:  cache,
		db:     db,
		store:  store,
	}
	api.router.Use(sessions.Sessions("session", api.store))
	api.loadRoutes()

	if gin.Mode() == gin.ReleaseMode {
		log.Print("Starting server on port ", api.port)
	}
	api.router.Run(api.port)
}

func (a *API) loadRoutes() {
	authHandler, err := auth.NewHandler(a.cache, a.db)
	if err != nil {
		log.Fatal("auth: ", err)
	}
	a.router.GET("/auth/tokens", authHandler.AuthRequired, authHandler.Tokens)
	a.router.GET("/auth/user", authHandler.AuthRequired, authHandler.User)
	a.router.POST("/auth/login", authHandler.Login)
	a.router.POST("/auth/logout", authHandler.AuthRequired, authHandler.Logout)
	a.router.POST("/auth/signup", authHandler.Signup)

	hardwareHandler, err := hardware.NewHandler(a.db)
	if err != nil {
		log.Fatal("hardware: ", err)
	}
	a.router.GET("/hardware/offers", hardwareHandler.Offers)
}
