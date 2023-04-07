package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"snippets/data"
	"snippets/data/wikipedia"
	"snippets/handlers"
	"syscall"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
)

// bindAddress is the address the server will listen on. It is set using the BIND_ADDRESS environment variable. By default, it will listen on port 9095.
var bindAddress = env.String("BIND_ADDRESS", false, ":9095", "Bind address for the server")

// backend is the backend to use. It is set using the BACKEND environment variable. By default, it will use the wikipedia API backend, which is the only backend currently implemented.
var backend = env.String("BACKEND", false, "wikipedia", "Backend to use. Default value is \"wikipedia\"")

// backend_api_url is the API url of the backend. It is set using the BACKEND_API_URL environment variable. By default, it will use the english wikipedia API url.
var backend_api_url = env.String("BACKEND_API_URL", false, "https://en.wikipedia.org/w/api.php", "API url of the backend. Default value is \"https://en.wikipedia.org/w/api.php\"")

func main() {

	// create a new logger
	l := hclog.Default()

	// parse environment variables
	err := env.Parse()
	if err != nil {
		l.Error("Unable to parse environment variables", "error", err)
		os.Exit(1)
	}

	var spBack data.SnippetsBackend

	// create the backend
	switch *backend {
	case "wikipedia":
		spBack, err = wikipedia.NewBackend(l, *backend_api_url)

	// TODO: Add more backends here. Maybe support other wikipedia languages, or other APIs like wikidata, or even a custom backend that uses a database.

	default:
		l.Error("Backend not supported", "backend", *backend)
		os.Exit(1)
	}

	if err != nil {
		l.Error("Unable to create backend", "error", err)
		os.Exit(1)
	}

	ss := handlers.NewSnippetsService(l, spBack)

	router := mux.NewRouter()

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/openapi.yaml"}
	sh := middleware.Redoc(opts, nil)

	router.Handle("/docs", sh)
	router.Handle("/openapi.yaml", http.FileServer(http.Dir("./")))
	router.HandleFunc("/snippets/{name}", ss.GetSnippetHandler).Methods("GET")

	s := &http.Server{
		Addr:         *bindAddress,
		Handler:      router,
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start the server in a goroutine so that it doesn't block
	go func() {
		l.Info("Starting server on", "address", *bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm to gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	// block until a signal is received
	sig := <-c
	l.Info("Got signal, gracefully shutting down", "signal", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	// we ignore the error returned from Shutdown because there is nothing useful we can do with it
	_ = s.Shutdown(tc)

}
