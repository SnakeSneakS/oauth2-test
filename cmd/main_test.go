package cmd_test

import (
	"log"
	"net/http"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/snakesneaks/oauth2-test/core"
	"github.com/snakesneaks/oauth2-test/service/handler"
)

func TestRun(t *testing.T) {
	config := core.GetConfig()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	router := gin.Default()
	router.Use(sessions.Sessions("auth-session", core.Store))

	handler := handler.NewMockHandler(mockCtrl)

	/*
		router.Static("/public", "web/static")
		router.LoadHTMLGlob("web/template/*")

		router.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "home.html", nil)
		})
	*/
	router.GET("/login", handler.Login)
	router.GET("/callback", handler.Callback)
	router.GET("/user", handler.User)
	router.GET("/logout", handler.Logout)

	log.Printf("server listening on http://localhost:%s\n", config.Host.PORT)
	if err := http.ListenAndServe("0.0.0.0:3000", router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
