package middleware

import (
	"context"
	"net/http"

	"github.com/tacogips/gaemux"
)

func FromHttpHandlerFunc(hdrFunc http.HandlerFunc) gaemux.Middleware {
	return func(c context.Context, hdr gaemux.ChainableHandlers) error {
		w, r := gaemux.ResponseWriter(c), gaemux.Request(c)

		hdrFunc(w, r)

		next, rest := hdr[0], hdr[1:]
		return next(c, rest)
	}
}
