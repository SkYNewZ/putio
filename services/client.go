package services

import (
	"context"
	"os"

	"github.com/putdotio/go-putio"
	"golang.org/x/oauth2"
)

var (
	client *putio.Client
)

// NewPutioClient return a Put.io client
func NewPutioClient() *putio.Client {
	if client != nil {
		return client
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("PUT_IO_TOKEN")})
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client = putio.NewClient(oauthClient)
	return client
}
