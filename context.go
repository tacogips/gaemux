package gaemux

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
)

type (
	routerParamKeyType    int
	queryParamKeyType     int
	requestKeyType        int
	resopnseWriterKeyType int
)

var (
	routerParamKey    routerParamKeyType    = 0
	queryParamKey     queryParamKeyType     = 1
	requestKey        requestKeyType        = 2
	responseWriterKey resopnseWriterKeyType = 3
)

func ctx(w http.ResponseWriter, r *http.Request, params httprouter.Params) context.Context {
	c := context.Background()
	c = appengine.WithContext(c, r)

	c = WithRouterParam(c, RouteParams(params))
	c = WithQueryParam(c, r.URL.Query())

	c = WithRequest(c, r)
	c = WithResponseWriter(c, w)

	return c
}

func Request(c context.Context) *http.Request {
	r, _ := c.Value(requestKey).(*http.Request)
	return r
}

func WithRequest(c context.Context, r *http.Request) context.Context {
	return context.WithValue(c, requestKey, r)
}

func ResponseWriter(c context.Context) http.ResponseWriter {
	w, _ := c.Value(responseWriterKey).(http.ResponseWriter)
	return w
}

func WithResponseWriter(c context.Context, w http.ResponseWriter) context.Context {
	return context.WithValue(c, responseWriterKey, w)
}
