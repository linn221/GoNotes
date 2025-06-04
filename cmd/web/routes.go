package main

import (
	"linn221/shop/handlers"
	"linn221/shop/middlewares"
	"log"
	"net/http"
	"time"
)

func (app *App) Serve() {

	authMux := http.NewServeMux()
	authMux.HandleFunc("GET /me", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world, you are authed"))
	})

	mainMux := http.NewServeMux()
	// public routes
	mainMux.HandleFunc("/login", handlers.HandleLogin(app.Renderer, app.DB, app.Cache))

	// mainMux.Handle("/api/", http.StripPrefix("/api", middlewares.Auth(authMux)))
	// file upload
	fileHandler := http.StripPrefix("/static", http.FileServer(http.Dir("files"))) // trailing slash or not to slash
	mainMux.Handle("/static/", fileHandler)

	srv := http.Server{
		Addr:         ":" + app.Port,
		Handler:      middlewares.Default(mainMux, app.Cache),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Fatal(srv.ListenAndServe())
}
