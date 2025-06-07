package handlers

import (
	"context"
	"linn221/shop/myctx"
	"linn221/shop/views"
	"net/http"
)

type DefaultSession struct {
	UserId int
}

type DefaultHandlerFunc func(ctx context.Context, r *http.Request, session *DefaultSession, vr *views.Renderer) error

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
