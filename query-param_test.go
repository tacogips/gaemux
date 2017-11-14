package gaemux

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetQueryParam(t *testing.T) {

	{
		c := WithQueryParam(context.Background(), url.Values{"foo": []string{"bar"}})
		v, err := QueryParamString(c, "foo")
		assert.Nil(t, err)
		assert.Equal(t, "bar", v)
	}

	{
		c := WithQueryParam(context.Background(), url.Values{"foo": []string{"10"}})
		v, err := QueryParamInt(c, "foo")
		assert.Nil(t, err)
		assert.Equal(t, 10, v)
	}

	{
		c := WithQueryParam(context.Background(), url.Values{"foo": []string{"10"}})
		v, err := QueryParamInt64(c, "foo")
		assert.Nil(t, err)
		assert.Equal(t, int64(10), v)
	}

}

func TestUnmarshallQueryParam(t *testing.T) {

	type Dst struct {
		Foo       string `query:"foo"`
		FooInt    int    `query:"foo_int"`
		FooInt64  int64  `query:"foo_int64"`
		FooInt64P *int64 `query:"foo_int64p"`
	}
	var dst Dst

	c := WithQueryParam(context.Background(),
		url.Values{
			"foo":        []string{"bar"},
			"foo_int":    []string{"10"},
			"foo_int64":  []string{"100"},
			"foo_int64p": []string{"1000"},
		})
	err := UnmarshallQueryParam(c, &dst)
	assert.Nil(t, err)

	var i int64 = 1000
	assert.Equal(t,
		Dst{
			Foo:       "bar",
			FooInt:    10,
			FooInt64:  100,
			FooInt64P: &i,
		}, dst)
}

func TestUnmarshallQueryParamWithDefault(t *testing.T) {

	type Dst struct {
		Foo       string `query:"foo" default:"baa"`
		FooInt    int    `query:"foo_int" default:"20"`
		FooInt64  int64  `query:"foo_int64" default:"200"`
		FooInt64P *int64 `query:"foo_int64p" default:"2000"`
	}
	var dst Dst

	{
		c := WithQueryParam(context.Background(),
			url.Values{})
		err := UnmarshallQueryParam(c, &dst)
		assert.Nil(t, err)

		var i int64 = 2000
		assert.Equal(t,
			Dst{
				Foo:       "baa",
				FooInt:    20,
				FooInt64:  200,
				FooInt64P: &i,
			}, dst)

	}

	{
		c := WithQueryParam(context.Background(),
			url.Values{
				"foo":        []string{"bar"},
				"foo_int":    []string{"10"},
				"foo_int64":  []string{"100"},
				"foo_int64p": []string{"1000"},
			})
		err := UnmarshallQueryParam(c, &dst)
		assert.Nil(t, err)

		var i int64 = 1000
		assert.Equal(t,
			Dst{
				Foo:       "bar",
				FooInt:    10,
				FooInt64:  100,
				FooInt64P: &i,
			}, dst)
	}
}
