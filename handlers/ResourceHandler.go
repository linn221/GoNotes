package handlers

import (
	"linn221/shop/myctx"
	"net/http"
	"strconv"
)

type Session struct {
	UserId int
	ResId  int
}

func ResourceHandler(handle func(http.ResponseWriter, *http.Request, *Session) error) http.HandlerFunc {
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

		err = handle(w, r, &session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
