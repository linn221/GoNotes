package handlers

import (
	"context"
	"linn221/shop/myctx"
	"net/http"
	"strconv"
)

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
