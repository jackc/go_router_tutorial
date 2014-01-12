package router

import (
	"io"
	"net/http"
)

type Router struct {
	endpoints map[string]http.Handler
}

// ServeHTTP makes Router implement standard http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, ok := r.endpoints[req.URL.Path]; ok {
		handler.ServeHTTP(w, req)
	} else {
		w.WriteHeader(404)
		io.WriteString(w, "404 Not Found")
	}
}

func (r *Router) AddRoute(method string, path string, handler http.Handler) {
	r.endpoints[path] = handler
}

func NewRouter() (r *Router) {
	r = new(Router)
	r.endpoints = make(map[string]http.Handler)
	return r
}
