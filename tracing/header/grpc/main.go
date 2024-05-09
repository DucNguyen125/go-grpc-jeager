package grpc

import "google.golang.org/grpc/metadata"

type Header metadata.MD

func (h Header) Get(key string) string {
	vals := metadata.MD(h).Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (h Header) Set(key, value string) {
	metadata.MD(h).Set(key, value)
}

func (h Header) Keys() []string {
	keys := make([]string, 0, len(h))
	for k := range metadata.MD(h) {
		keys = append(keys, k)
	}
	return keys
}
