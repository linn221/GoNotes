package handlers

import (
	"linn221/shop/myctx"
	"net/http"
)

type CreateResource[T any] struct {
	parseInput       func(*http.Request) (*T, formErrors)
	handleParseError func(http.ResponseWriter, *http.Request, *DefaultSession, *T, formErrors) error
	handle           func(http.ResponseWriter, *http.Request, *DefaultSession, *T) error
}

func CreateHandler[T any](res CreateResource[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, err := myctx.GetUserId(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// resIdStr := r.PathValue("id")
		// if resIdStr == "" {
		// 	http.Error(w, "resource id is required", http.StatusBadRequest)
		// 	return
		// }
		// resId, err := strconv.Atoi(resIdStr)
		// if err != nil {
		// 	http.Error(w, "resource id is required", http.StatusBadRequest)
		// 	return

		// }
		session := DefaultSession{
			UserId: userId,
		}

		input, ferrs := res.parseInput(r)
		if len(ferrs) > 0 {
			finalErrHandle(w,
				res.handleParseError(w, r, &session, input, ferrs),
			)
			return
		}

		err = res.handle(w, r, &session, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// type CreateFunc[T any] func(w http.ResponseWriter, r *http.Request, session Session, input *T) error

// func CreateHandler[T any](handle CreateFunc[T]) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		ctx := r.Context()
// 		userId, shopId, err := myctx.GetIdsFromContext(ctx)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		input, ok, err := parseJson[T](w, r)
// 		if !ok {
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 			}
// 			return
// 		}

// 		CreateSession := Session{
// 			UserId: userId,
// 		}
// 		err = handle(w, r, CreateSession, input)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }
