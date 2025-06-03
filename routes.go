package main

import (
	"linn221/shop/handlers"
	"linn221/shop/middlewares"
	"linn221/shop/models"
	"log"
	"net/http"
	"time"
)

func (app *App) Serve() {

	authMux := http.NewServeMux()

	authMux.Handle("GET /me", handlers.Me(app.DB))
	authMux.HandleFunc("POST /change-password", handlers.ChangePassword(app.DB))
	authMux.HandleFunc("POST /logout", handlers.Logout(app.DB, app.Cache))
	authMux.HandleFunc("POST /update-profile", handlers.UpdateUserInfo(app.DB))

	//categories
	authMux.HandleFunc("POST /categories", handlers.HandleCategoryCreate(app.DB,
		app.Readers.CategoryListService.CleanCache,
	))
	authMux.HandleFunc("PUT /categories/{id}", handlers.HandleCategoryUpdate(app.DB,
		app.Readers.AfterCategoryUpdate,
	))
	authMux.HandleFunc("DELETE /categories/{id}", handlers.HandleCategoryDelete(app.DB,
		app.Readers.AfterCategoryUpdate,
	))
	authMux.HandleFunc("PATCH /categories/{id}/toggle",
		handlers.HandleToggleActive[models.Category](app.DB,
			app.Readers.AfterCategoryUpdate,
		))
	authMux.HandleFunc("GET /categories/{id}", handlers.DefaultGetHandler(app.Readers.CategoryGetService))
	authMux.HandleFunc("GET /categories", handlers.DefaultListHandler(app.Readers.CategoryListService))
	authMux.HandleFunc("GET /categories/inactive",
		handlers.ListInactiveHandler[models.Category, models.CategoryResource](app.DB),
	)

	//units
	authMux.HandleFunc("POST /units", handlers.HandleUnitCreate(app.DB,
		app.Readers.UnitListService.CleanCache,
	))
	authMux.HandleFunc("PUT /units/{id}", handlers.HandleUnitUpdate(app.DB,
		app.Readers.AfterUnitUpdate,
	))
	authMux.HandleFunc("DELETE /units/{id}", handlers.HandleUnitDelete(app.DB,
		app.Readers.AfterUnitUpdate,
	))
	authMux.HandleFunc("PATCH /units/{id}/toggle", handlers.HandleToggleActive[models.Unit](app.DB,
		app.Readers.AfterUnitUpdate,
	))
	authMux.HandleFunc("GET /units/{id}", handlers.DefaultGetHandler(app.Readers.UnitGetService))
	authMux.HandleFunc("GET /units", handlers.DefaultListHandler(app.Readers.UnitListService))
	authMux.HandleFunc("GET /units/inactive", handlers.ListInactiveHandler[models.Unit, models.UnitResource](app.DB))

	// items
	authMux.HandleFunc("POST /items", handlers.HandleItemCreate(app.DB,
		app.Readers.ItemListService.CleanCache,
	))
	authMux.HandleFunc("PUT /items/{id}", handlers.HandleItemUpdate(app.DB,
		app.Readers.ItemGetService.CleanCache,
		app.Readers.ItemListService.CleanCache,
	))
	authMux.HandleFunc("DELETE /items/{id}", handlers.HandleItemDelete(app.DB,
		app.Readers.ItemGetService.CleanCache,
		app.Readers.ItemListService.CleanCache,
	))
	authMux.HandleFunc("PATCH /items/{id}/toggle",
		handlers.HandleToggleActive[models.Item](app.DB, app.Readers.AfterUpdateItem),
	)
	authMux.HandleFunc("GET /items/{id}", handlers.DefaultGetHandler(app.Readers.ItemGetService))
	authMux.HandleFunc("GET /items", handlers.HandleItemIndex(app.DB))
	authMux.HandleFunc("GET /items/all", handlers.DefaultListHandler(app.Readers.ItemListService))
	authMux.HandleFunc("GET /items/inactive",
		handlers.ListCustomInactiveHandler(app.DB, models.FetchInactiveItemResources),
	)

	mainMux := http.NewServeMux()
	// public routes
	mainMux.HandleFunc("POST /upload-single", handlers.HandleImageUploadSingle(app.DB, app.ImageDirectory))
	mainMux.HandleFunc("POST /register", handlers.Register(app.DB))
	mainMux.HandleFunc("POST /login", handlers.Login(app.DB, app.Cache))

	mainMux.Handle("/api/", http.StripPrefix("/api", middlewares.Auth(app.ResourceRateLimit(authMux))))

	srv := http.Server{
		Addr:         ":" + app.Port,
		Handler:      middlewares.Default(app.GeneralRateLimit(mainMux), app.Cache),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Fatal(srv.ListenAndServe())
}
