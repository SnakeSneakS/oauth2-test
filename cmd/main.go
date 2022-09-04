package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/snakesneaks/oauth2-test/core"
	"github.com/snakesneaks/oauth2-test/service/handler"
)

func Init() {
	core.InitSessionStore()
}

func Run() {
	config := core.GetConfig()
	router := gin.Default()
	router.Use(sessions.Sessions("auth-session", core.Store))

	handler := handler.NewHandler(*config)

	router.Static("/public", "web/static")
	router.LoadHTMLGlob("web/html/template/*")

	router.GET("/", handler.Home)
	router.GET("/login", handler.Login)
	router.GET("/callback", handler.Callback)
	router.GET("/user", handler.User)
	router.GET("/logout", handler.Logout)

	log.Printf("server listening on http://localhost:%s\n", config.Host.PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.Host.PORT), router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
