package main

import (
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"ecom-local/internal/json"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	slogchi "github.com/samber/slog-chi"
)

type config struct {
	addr string
}

type application struct {
	config config
	data   StaticData
}

type StaticData struct {
	BirthDate time.Time `json:"birth_date"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(slogchi.New(slog.Default()))
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"message": "Welcome to Eloi's API!",
			"version": "1.0.0",
			"endpoints": []map[string]string{
				{
					"path":        "/",
					"method":      "GET",
					"description": "API entrypoint, lists available endpoints",
				},
				{
					"path":        "/health",
					"method":      "GET",
					"description": "Health check endpoint",
				},
				{
					"path":        "/infos",
					"method":      "GET",
					"description": "Information about the client connection and API static data",
				},
				{
					"path":        "/myip",
					"method":      "GET",
					"description": "Returns your IP address",
				},
			},
		}
		json.Write(w, http.StatusOK, response)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("All good"))
	})

	r.Get("/myip", func(w http.ResponseWriter, r *http.Request) {

		ip := r.RemoteAddr

		if strings.Contains(ip, "::1") {
			ip = "127.0.0.1 (Localhost)"
		}

		response := map[string]string{
			"ip": ip,
		}

		json.Write(w, http.StatusOK, response)

	})

	r.Get("/infos", func(w http.ResponseWriter, r *http.Request) {

		ip := r.RemoteAddr
		userAgent := r.UserAgent()

		if strings.Contains(ip, "::1") {
			ip = "127.0.0.1 (Localhost)"
		}

		response := map[string]string{
			"ip":         ip,
			"user_agent": userAgent,
			"message":    "Welcome to Eloi's API !",
		}

		json.Write(w, http.StatusOK, response)

	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}
