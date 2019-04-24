package base

import (
	"fmt"
	"net/http"
	"net/url"
)

type requestOption interface {
	updateRequest(*http.Request) error
}

type requestOptionFunc func(*http.Request) error

func (f requestOptionFunc) updateRequest(r *http.Request) error {
	return f(r)
}

func PathParams(parameters ...string) requestOptionFunc {
	return func(req *http.Request) error {
		var encoded []interface{}
		for _, param := range parameters {
			encoded = append(encoded, url.QueryEscape(param))
		}

		path := fmt.Sprintf(req.URL.Path, encoded...)

		parsed, err := url.Parse(path)
		if err != nil {
			return err
		}

		req.URL = req.URL.ResolveReference(parsed)
		return nil
	}
}

type QueryValues []struct{ Key, Value string }

func (values QueryValues) updateRequest(req *http.Request) error {
	qs := req.URL.Query()
	for _, item := range values {
		qs.Add(item.Key, item.Value)
	}

	req.URL.RawQuery = qs.Encode()
	return nil
}
