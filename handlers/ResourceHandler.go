package handlers

import (
	"context"
	"linn221/shop/myctx"
	"linn221/shop/views"
	"net/http"
	"strconv"
)

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
