package gaemux

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine/log"
)

type Handler func(context.Context) error

type ChainableHandler func(context.Context, ChainableHandlers) error
type ChainableHandlers []ChainableHandler
type Middleware ChainableHandler
type PanicHandler func(context.Context, interface{})

type GaeMux struct {
	r            *httprouter.Router
	middlewares  []ChainableHandler //TODO Group
	PanicHandler PanicHandler
}

var defaultPanicHandler PanicHandler = func(c context.Context, panicbody interface{}) {
	log.Errorf(c, "panic occure [%#v]", panicbody)
	InternalError(c)
}

func New() *GaeMux {
	return &GaeMux{
		r: httprouter.New(),
	}
}

func (gm *GaeMux) Handler() http.Handler {
	return gm.r
}

func (gm *GaeMux) Use(m Middleware) {
	gm.middlewares = append(gm.middlewares, ChainableHandler(m))
}

func (gm *GaeMux) Get(path string, handler Handler) {
	gm.r.GET(path, gm.httpRouterize(handler))
}

func (gm *GaeMux) Post(path string, handler Handler) {
	gm.r.POST(path, gm.httpRouterize(handler))
}

func (gm *GaeMux) Put(path string, handler Handler) {
	gm.r.PUT(path, gm.httpRouterize(handler))
}

func (gm *GaeMux) Delete(path string, handler Handler) {
	gm.r.DELETE(path, gm.httpRouterize(handler))
}

func (gm *GaeMux) httpRouterize(gaeHandler Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		c := ctx(w, r, params)

		err := gm.handle(c, gaeHandler)
		if err != nil {
			//TODO custom error handler
			http.Error(w, "error", http.StatusBadRequest)
		} else {
		}
	}
}

func (gm *GaeMux) handle(c context.Context, hdr Handler) error {
	hdrs := append(gm.middlewares, func(c context.Context, _ ChainableHandlers) error {
		return hdr(c)
	})

	defer func() {
		if r := recover(); r != nil {
			gm.PanicHandler(c, r)
		}
	}()

	next, rest := hdrs[0], hdrs[1:]
	return next(c, rest)
}
