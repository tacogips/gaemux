package gaemux

import (
	"context"
	"encoding/json"
	"net/http"
)

const (
	contentType     = "Content-Type"
	applicationJson = "application/json; charset=UTF-8"
)

func OK(c context.Context) error {
	w := ResponseWriter(c)

	w.WriteHeader(http.StatusOK)
	return nil
}

func BadRequest(c context.Context) error {
	w := ResponseWriter(c)
	w.WriteHeader(http.StatusBadRequest)
	return nil
}

func InternalError(c context.Context) error {
	w := ResponseWriter(c)
	w.WriteHeader(http.StatusInternalServerError)
	return nil
}

func OKJson(c context.Context, body interface{}) error {
	w := ResponseWriter(c)

	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(http.StatusOK)

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func BadRequestJson(c context.Context, body interface{}) error {
	w := ResponseWriter(c)

	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(http.StatusBadRequest)

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}
