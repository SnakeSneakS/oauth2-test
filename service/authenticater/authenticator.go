package authenticator

import (
	"context"
	"errors"
	"log"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/snakesneaks/oauth2-test/core"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	*oidc.Provider
	oauth2.Config
}

func NewAuthenticator(config core.Config) Authenticator {
	provider, err := oidc.NewProvider(
		context.Background(),
		"https://"+config.Auth0.AUTH0_DOMAIN+"/",
	)
	if err != nil {
		log.Fatal(err)
	}

	return Authenticator{
		Provider: provider,
		Config: oauth2.Config{
			ClientID:     config.Auth0.AUTH0_CLIENT_ID,
			ClientSecret: config.Auth0.AUTH0_CLIENT_SECRET,
			RedirectURL:  config.Auth0.AUTH0_CALLBACK_URL,
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		},
	}
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}
