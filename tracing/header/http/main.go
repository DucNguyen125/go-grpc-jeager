package http

import "net/http"

type Header http.Header

func (h Header) Get(key string) string {
	return http.Header(h).Get(key)
}

func (h Header) Set(key string, value string) {
	http.Header(h).Set(key, value)
}

func (h Header) Keys() []string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	return keys
}
