package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type filesTemplate struct {
	Path     string
	Items    *[]file
	Previous string
}

// RenderTemplate return template and send as response
func RenderTemplate(w http.ResponseWriter, t *filesTemplate) error {
	tpl := template.Must(template.ParseFiles("templates/index.go.html"))
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
