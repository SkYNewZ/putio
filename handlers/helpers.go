package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
)

type filesTemplate struct {
	Path     string
	Items    *[]file
	Previous string
}

// RenderTemplate return template and send as response
func RenderTemplate(w http.ResponseWriter, t *filesTemplate) error {
	// Embed templates files in binary
	box := packr.New("templates", "./templates")

	s, err := box.FindString("index.go.html")
	if err != nil {
		return err
	}

	tpl, err := template.New("index").Parse(s)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tpl.Execute(w, t)
}

// Fatal throw internal server error
func Fatal(w http.ResponseWriter, err error) {
	log.Errorf("Unexpected error: %s", err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func parseDirectoryID(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	result, err := strconv.Atoi(vars["folder"])
	if err != nil {
		return 0, err
	}
	return int64(result), nil
}
