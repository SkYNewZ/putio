// +build embedded

package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var putIOToken string

// Insert Put.io token when building for friends
func init() {
	if putIOToken == "" {
		log.Fatalln(`Please specify PUT_IO_TOKEN during build. -ldflags="-X 'main.putIOToken=XXX'"`)
	}

	t := "files"
	vars := map[string]string{
		"PUT_IO_TOKEN":     putIOToken,
		"OFUSCATION_TOKEN": t,
		"NO_AUTH":          "1",
	}
	for k, v := range vars {
		os.Setenv(k, v)
	}

	log.Printf("Go to http://localhost:8000/%s/0", t)
}
