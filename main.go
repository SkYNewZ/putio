package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Define different logger for Google Cloud Run
	if v := os.Getenv("K_SERVICE"); v != "" {
		log.SetFormatter(joonix.NewFormatter())
	}

	log.SetLevel(logrus.DebugLevel)
}

func main() {
	if os.Getenv("PUT_IO_TOKEN") == "" {
		log.Fatalln("You must specify a PUT_IO_TOKEN. Please refer to https://app.put.io/settings/account/oauth/authorizations")
	}

	var port string = "8000"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}

	// Ticker for refresh in-memory users
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      newRouter(ticker),
	}

	go func() {
		log.Printf("Listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_ = srv.Shutdown(ctx)
	log.Println("Shutting down server")
	os.Exit(0)
}
