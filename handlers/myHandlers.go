package handlers

import (
	"context"
	"linn221/shop/myctx"
	"linn221/shop/services"
	"linn221/shop/views"
	"net/http"
	"strconv"
)

type DefaultSession struct {
	UserId int
}

type DefaultHandlerFunc func(ctx context.Context, r *http.Request, session *DefaultSession, vr *views.Renderer) error

func handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

type MinHandlerFunc func(w http.ResponseWriter, r *http.Request, userId int) error

func MinHandler(handle MinHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, err := myctx.GetUserId(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = handle(w, r, userId)
		if err != nil {
			handleError(w, err)
		}
	}
}

func DefaultHandler(t *views.Templates, handle DefaultHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, err := myctx.GetUserId(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session := DefaultSession{
			UserId: userId,
		}

		renderer := t.NewRenderer(w, userId)

		err = handle(ctx, r, &session, renderer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type Session struct {
	UserId int
	ResId  int
}

type ResourceHandlerFunc func(ctx context.Context, r *http.Request, session *Session, vr *views.Renderer) error

func ResourceHandler(t *views.Templates, handle ResourceHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, err := myctx.GetUserId(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resIdStr := r.PathValue("id")
		if resIdStr == "" {
			http.Error(w, "resource id is required", http.StatusBadRequest)
			return
		}
		resId, err := strconv.Atoi(resIdStr)
		if err != nil {
			http.Error(w, "resource id is required", http.StatusBadRequest)
			return

		}
		session := Session{
			UserId: userId,
			ResId:  resId,
		}
		renderer := t.NewRenderer(w, userId)

		err = handle(ctx, r, &session, renderer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type DeleteHandlerFunc func(ctx context.Context, r *http.Request, userId int, resId int) error

func DeleteHandler(handle DeleteHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, err := myctx.GetUserId(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resIdStr := r.PathValue("id")
		if resIdStr == "" {
			http.Error(w, "resource id is required", http.StatusBadRequest)
			return
		}
		resId, err := strconv.Atoi(resIdStr)
		if err != nil {
			http.Error(w, "resource id is required", http.StatusBadRequest)
			return

		}

		err = handle(ctx, r, userId, resId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func CreateHandler[T any](t *views.Templates,
	parseInput func(*http.Request) (*T, services.FormErrors),
	handle func(w http.ResponseWriter, r *http.Request, session *DefaultSession, input *T, fe services.FormErrors, vr *views.Renderer) error,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, err := myctx.GetUserId(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session := DefaultSession{
			UserId: userId,
		}

		renderer := t.NewRenderer(w, userId)
		input, ferrs := parseInput(r)

		err = handle(w, r, &session, input, ferrs, renderer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func UpdateHandler[T any](t *views.Templates,
	parseInput func(r *http.Request) (*T, services.FormErrors),
	handle func(w http.ResponseWriter, r *http.Request, session *Session, input *T, fe services.FormErrors, renderer *views.Renderer) error,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, err := myctx.GetUserId(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resIdStr := r.PathValue("id")
		if resIdStr == "" {
			http.Error(w, "resource id is required", http.StatusBadRequest)
			return
		}
		resId, err := strconv.Atoi(resIdStr)
		if err != nil {
			http.Error(w, "resource id is required", http.StatusBadRequest)
			return

		}
		session := Session{
			UserId: userId,
			ResId:  resId,
		}

		renderer := t.NewRenderer(w, userId)
		input, ferrs := parseInput(r)

		err = handle(w, r, &session, input, ferrs, renderer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
