package gaemux

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

var NoQueryParamError = errors.New("no query param")

func WithQueryParam(c context.Context, queryParams url.Values) context.Context {
	return context.WithValue(c, queryParamKey, queryParams)
}

func QueryParams(c context.Context) url.Values {
	if value, ok := c.Value(queryParamKey).(url.Values); !ok {
		return url.Values{}
	} else {
		return value
	}
}

func QueryParamString(c context.Context, name string) (string, error) {
	vals := QueryParams(c)[name]
	if len(vals) == 0 {
		return "", NoQueryParamError
	}
	return vals[0], nil
}

func QueryParamInt(c context.Context, name string) (int, error) {
	if v, err := QueryParamString(c, name); err != nil {
		return 0, err
	} else {
		if intValue, parseErr := strconv.Atoi(v); parseErr != nil {
			return 0, err
		} else {
			return intValue, nil
		}
	}
}

func QueryParamInt64(c context.Context, name string) (int64, error) {
	if v, err := QueryParamString(c, name); err != nil {
		return 0, err
	} else {
		if intValue, parseErr := strconv.ParseInt(v, 10, 64); parseErr != nil {
			return 0, err
		} else {
			return intValue, nil
		}
	}
}

func UnmarshallQueryParam(c context.Context, dst interface{}) error {
	queryParmaVals := QueryParams(c)

	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
		return fmt.Errorf("dst must be not nil pointer to struct")
	}

	indirectDstVal := reflect.Indirect(dstVal)
	tp := indirectDstVal.Type()
	k := tp.Kind()
	if k != reflect.Struct {
		return fmt.Errorf("dst must be pointer to struct")
	}
	queryTag := "query"
	defaultTag := "default"

	for i := 0; i < tp.NumField(); i++ {
		field := indirectDstVal.Field(i)
		if !field.CanSet() {
			continue
		}

		typeOfField := tp.Field(i)

		queryParamName := typeOfField.Tag.Get(queryTag)
		defaultVal := typeOfField.Tag.Get(defaultTag)

		if len(queryParamName) == 0 || queryParamName == "-" {
			continue
		}
		if queryParamVals, ok := queryParmaVals[queryParamName]; ok {
			unmarshalField(field, typeOfField.Type.Kind(), queryParamVals[0])
		} else if len(defaultVal) != 0 {
			unmarshalField(field, typeOfField.Type.Kind(), defaultVal)
		}
	}
	return nil
}

func unmarshalField(field reflect.Value, fieldKind reflect.Kind, val string) error {
	switch fieldKind {

	case reflect.Ptr:
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		unmarshalField(field.Elem(), field.Elem().Kind(), val)
	case reflect.String:
		field.SetString(val)
	case reflect.Int:
		v, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return err
		}
		field.SetInt(v)
	case reflect.Int64:
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(v)
	case reflect.Bool:
		field.SetBool(val == "true")

	case reflect.Float64:
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		field.SetFloat(v)
	default:
		return fmt.Errorf("not supported query param type [%s]", field.Type().Name())
	}
	return nil
}
