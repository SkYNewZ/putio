package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/SkYNewZ/putio/handlers"
	"github.com/SkYNewZ/putio/services"
	"github.com/gobuffalo/packr/v2"
	muxHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	staticDir     string = "/assets/"
	staticDirPath string = "./assets"
)

// Define our struct
type authenticationMiddleware struct {
	users map[string]string
}

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

// Initialize it somewhere
func (amw *authenticationMiddleware) Populate() error {
	client, err := services.NewFirestoreClient()
	if err != nil {
		return err
	}

	users, err := client.GetUsers()
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return errors.New("No user found in Firestore")
	}

	log.Infof("Populate users database with %d users", len(users))
	for _, u := range users {
		amw.users[u.Username] = u.Password
	}
	return nil
}

// HTTP Basic auth middleware
func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		// If auth not provided
		if !ok || !checkPasswordHash(pass, amw.users[user]) {
			w.Header().Add("WWW-Authenticate", "Basic realm=Protected area")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Ensure Content-Type is HTML
func contentTypeMiddleware(next http.Handler) http.Handler {
	return muxHandlers.ContentTypeHandler(next, "text/html")
}

// Log requests in Apache Common format
func loggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := statusWriter{ResponseWriter: w}
		next.ServeHTTP(&sw, r)
		duration := time.Since(start)

		log.Printf("%s %s %d %d - %d ms", r.Method, r.RequestURI, sw.status, sw.length, duration.Milliseconds())
	})
}

// Disable directory listing when using http.FileServer
func neuterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func newRouter(ticker *time.Ticker) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	// Embed assets in binary
	box := packr.New(staticDir, staticDirPath)
	r.PathPrefix(staticDir).Methods(http.MethodGet).Handler(http.StripPrefix(staticDir, http.FileServer(box)))

	ofuscationToken := os.Getenv("OFUSCATION_TOKEN")
	if ofuscationToken == "" {
		log.Fatalln("You must set OFUSCATION_TOKEN")
	}

	// Routes /XXXXXXXXXXXXX
	log.Printf("Ofuscation token: %s", ofuscationToken)
	routes := r.PathPrefix(fmt.Sprintf("/%s", ofuscationToken)).Subrouter()
	routes.HandleFunc("/{folder:[0-9]+}", handlers.FilesHandler).Name("folder")

	// Middleware
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(muxHandlers.CompressHandler)
	r.Use(contentTypeMiddleware)
	r.Use(loggingMiddleWare)
	r.Use(neuterMiddleware)

	if os.Getenv("NO_AUTH") != "1" {
		// Ensure Firestore access is specified
		if os.Getenv("GOOGLE_CLOUD_PROJECT") == "" {
			log.Fatalln("You must specify set GOOGLE_CLOUD_PROJECT for firestore works")
		}

		// Only protect main routes, not /icons routes
		amw := authenticationMiddleware{make(map[string]string)}

		// Scheduled refreshing users
		go func() {
			log.Debugln("Start routine for refresh in-memory users")
			for ; true; <-ticker.C {
				if err := amw.Populate(); err != nil {
					log.Warnln(err)
				}
			}
		}()
		routes.Use(amw.Middleware)
	}

	return r
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
