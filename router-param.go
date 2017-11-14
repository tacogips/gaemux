package gaemux

import (
	"context"
	"errors"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

var NoRouterParamError = errors.New("no router param")

type RouteParams httprouter.Params

func (rp RouteParams) ByNameString(name string) string {
	return httprouter.Params(rp).ByName(name)
}

func WithRouterParam(c context.Context, params RouteParams) context.Context {
	return context.WithValue(c, routerParamKey, params)
}

func RouterParams(c context.Context) RouteParams {
	rp, _ := c.Value(routerParamKey).(RouteParams)
	return rp
}

func RouterParamString(c context.Context, name string) (string, error) {
	rp := RouterParams(c)
	s := rp.ByNameString(name)
	if len(s) == 0 {
		return "", NoRouterParamError
	}

	return s, nil
}

func RouterParamInt(c context.Context, name string) (int, error) {
	if v, err := RouterParamString(c, name); err != nil {
		return 0, err
	} else {
		if intValue, parseErr := strconv.Atoi(v); parseErr != nil {
			return 0, err
		} else {
			return intValue, nil
		}
	}
}

func RouterParamInt64(c context.Context, name string) (int64, error) {
	if v, err := RouterParamString(c, name); err != nil {
		return 0, err
	} else {
		if intValue, parseErr := strconv.ParseInt(v, 10, 64); parseErr != nil {
			return 0, err
		} else {
			return intValue, nil
		}
	}
}
