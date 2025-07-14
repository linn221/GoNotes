package app

import (
	"context"
	"fmt"
	"linn221/shop/handlers"
	"linn221/shop/middlewares"
	"linn221/shop/myctx"
	"log"
	"net/http"
	"time"
)

func (app *App) Serve() {

	authMux := http.NewServeMux()
	authMux.HandleFunc("GET /me", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world, you are authed"))
	})
	authMux.HandleFunc("POST /logout", handlers.HandleLogout(app.Cache))

	t := app.Templates
	myServices := app.Services
	authMux.HandleFunc("/change-password", handlers.HandleChangePassword(t, myServices.UserService, app.Cache))

	//labels
	authMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/labels", http.StatusPermanentRedirect)
	})
	authMux.HandleFunc("GET /labels", handlers.ShowLabelIndex(t, myServices.LabelService))
	authMux.HandleFunc("GET /labels/new", handlers.ShowLabelCreate(t))
	authMux.HandleFunc("GET /labels/{id}/edit", handlers.ShowLabelEdit(t, myServices.LabelService))
	authMux.HandleFunc("POST /labels/{id}/toggle", handlers.HandleLabelToggleActive(t, myServices.LabelService))
	authMux.HandleFunc("POST /labels", handlers.HandleLabelCreate(t, myServices.LabelService))
	authMux.HandleFunc("PUT /labels/{id}", handlers.HandleLabelUpdate(t, myServices.LabelService))
	authMux.HandleFunc("DELETE /labels/{id}", handlers.HandleLabelDelete(myServices.LabelService))

	//notes
	getTimezone := func(ctx context.Context) string {
		token, err := myctx.GetToken(ctx)
		if err != nil {
			return "UTC"
		}
		s, err := app.Cache.GetH(fmt.Sprintf("Token:%s", token), "timezone")
		if err != nil {
			return "UTC"
		}
		return s
	}

	authMux.HandleFunc("GET /notes/new", handlers.ShowNoteCreate(t, myServices.LabelService))
	authMux.HandleFunc("GET /notes/{id}/edit", handlers.ShowNoteEdit(t, myServices.NoteService, myServices.LabelService))
	authMux.HandleFunc("GET /notes", handlers.ShowNoteIndex(t, myServices.NoteService, myServices.LabelService, getTimezone))
	authMux.HandleFunc("POST /notes", handlers.HandleNoteCreate(t, myServices.NoteService, myServices.LabelService))
	authMux.HandleFunc("POST /notes/export", handlers.HandleNoteExport(myServices.NoteService))
	authMux.HandleFunc("GET /notes/import", handlers.ShowNoteImport(t))
	authMux.HandleFunc("POST /notes/import", handlers.HandleNoteImport(myServices.NoteService))
	authMux.HandleFunc("PATCH /notes/{id}", handlers.HandleNotePartialUpdate(t, myServices.NoteService, myServices.LabelService, getTimezone))
	authMux.HandleFunc("PUT /notes/{id}", handlers.HandleNoteUpdate(t, myServices.NoteService, myServices.LabelService))
	authMux.HandleFunc("DELETE /notes/{id}", handlers.HandleNoteDelete(myServices.NoteService))

	mainMux := http.NewServeMux()
	// public routes

	mainMux.Handle("/login", middlewares.RedirectIfAuthenticated(handlers.HandleLogin(t, myServices.UserService, app.Cache)))
	mainMux.Handle("/register", middlewares.RedirectIfAuthenticated(handlers.HandleRegister(t, myServices.UserService)))

	// mainMux.Handle("/api/", http.StripPrefix("/api", middlewares.Auth(authMux)))
	// file upload
	fileHandler := http.StripPrefix("/static", http.FileServer(http.Dir(app.AssetDirectory))) // trailing slash or not to slash
	mainMux.Handle("/static/", fileHandler)
	mainMux.Handle("/", middlewares.Auth(authMux))

	srv := http.Server{
		Addr:         ":" + app.Port,
		Handler:      middlewares.Default(mainMux, app.Cache),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	fmt.Println("server started. visit http://localhost:" + app.Port)
	log.Fatal(srv.ListenAndServe())
}
