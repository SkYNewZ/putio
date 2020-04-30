package handlers

import (
	"context"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/SkYNewZ/putio/services"
	"github.com/gorilla/mux"
	"github.com/putdotio/go-putio"
)

type file struct {
	putio.File
	URL string
}

// FilesHandler handle list files
func FilesHandler(w http.ResponseWriter, r *http.Request) {
	// Get directory from request URL
	directory, err := parseDirectoryID(r)
	if err != nil {
		Fatal(w, err)
		return
	}

	// Read files in cache
	var filesTpl filesTemplate
	var cache services.CacheInterface = services.NewCache(os.Getenv("NO_CACHE") != "1")
	data, found := cache.Get(strconv.Itoa(int(directory)))
	if found {
		log.Debugln("Use cached files")
		filesTpl = *data.(*filesTemplate)
		if err := RenderTemplate(w, &filesTpl); err != nil {
			Fatal(w, err)
		}
		return
	}

	// Not found in cache, get real files on server
	client := services.NewPutioClient()
	log.Debugf("List Put.io files in directory %d", directory)
	children, parent, err := client.Files.List(context.Background(), directory)

	// Handle Put.io 404
	if e, ok := err.(*putio.ErrorResponse); ok {
		log.Debugf("Put.io throw 404: %s", e.Message)
		http.Error(w, e.Response.Status, e.Response.StatusCode)
		return
	}

	if err != nil {
		Fatal(w, err)
		return
	}

	log.Debugf("Found %d items in '%s'", len(children), parent.Name)

	// Get current route for redirection
	route := mux.CurrentRoute(r)
	u, err := route.URL("folder", strconv.Itoa(int(parent.ParentID)))
	if err != nil {
		Fatal(w, err)
		return
	}

	// Create HTML template data
	filesTpl = filesTemplate{Path: parent.Name, Previous: u.String()}

	// Get stream URL for each files found
	log.Debugln("Get download/steam URLs")
	var files []file = make([]file, 0)
	var url string = u.String()
	for _, f := range children {
		if !f.IsDir() {
			url, err = client.Files.URL(context.Background(), f.ID, false)
			if err != nil {
				Fatal(w, err)
				return
			}
		}
		files = append(files, file{f, url})
	}

	// Render template
	filesTpl.Items = &files

	// Store is cache
	cache.Set(strconv.Itoa(int(directory)), &filesTpl)

	if err := RenderTemplate(w, &filesTpl); err != nil {
		Fatal(w, err)
		return
	}
}
