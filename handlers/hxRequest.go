package handlers

import "net/http"

type HxRequest struct {
	ByHtmx  bool
	Boosted bool
}

func parseHxRequest(r *http.Request) *HxRequest {
	return &HxRequest{
		ByHtmx:  r.Header.Get("HX-Request") == "true",
		Boosted: r.Header.Get("HX-Boosted") == "true",
	}
}
