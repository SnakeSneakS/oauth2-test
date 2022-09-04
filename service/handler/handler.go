//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/snakesneaks/oauth2-test/core"
	authenticator "github.com/snakesneaks/oauth2-test/service/authenticater"
)

type Handler interface {
	Home(*gin.Context)
	Login(*gin.Context)
	Callback(*gin.Context)
	User(*gin.Context)
	Logout(*gin.Context)
}

type handler struct {
	config        core.Config
	authenticator authenticator.Authenticator
}

func NewHandler(config core.Config) Handler {
	return handler{
		config:        config,
		authenticator: authenticator.NewAuthenticator(config),
	}
}

// Home redirect to /login
func (h handler) Home(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}

// Login redirect to login page
func (h handler) Login(ctx *gin.Context) {
	state, err := core.GenerateRandomState()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Save the state inside the session.
	session := sessions.Default(ctx)
	session.Set("state", state)
	if err := session.Save(); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, h.authenticator.AuthCodeURL(state))
}

// Callback after login
func (h handler) Callback(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if ctx.Query("state") != session.Get("state") {
		ctx.String(http.StatusBadRequest, "invalid state parameter.")
		return
	}

	//exchange an authorization code for a token
	token, err := h.authenticator.Exchange(ctx.Request.Context(), ctx.Query("code"))
	if err != nil {
		ctx.String(http.StatusUnauthorized, fmt.Sprintf("Failed to exchange an authorization code for a token. %v", err))
		return
	}

	idToken, err := h.authenticator.VerifyIDToken(ctx.Request.Context(), token)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("failed to verify id token. %v", err))
		return
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	session.Set("access_token", token.AccessToken)
	session.Set("profile", profile)
	if err := session.Save(); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("User login callback receive.\nprofile: %+#v", profile)

	ctx.Redirect(http.StatusTemporaryRedirect, "/user")
}

func (h handler) User(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	ctx.HTML(http.StatusOK, "user.html", profile)
}
func (h handler) Logout(ctx *gin.Context) {
	logoutUrl, err := url.Parse(fmt.Sprintf("https://%s/v2/logout", h.config.Auth0.AUTH0_DOMAIN))
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(fmt.Sprintf("%s://%s", scheme, ctx.Request.Host))
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", h.config.Auth0.AUTH0_CLIENT_ID)
	logoutUrl.RawQuery = parameters.Encode()

	session := sessions.Default(ctx)
	session.Clear()
	session.Save()

	ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}
