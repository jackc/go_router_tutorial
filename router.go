package router

import (
	"io"
	"net/http"
)

type Router struct {
}

// ServeHTTP makes Router implement standard http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
	io.WriteString(w, "404 Not Found")
}

func (r *Router) AddRoute(method string, path string, handler http.Handler) {
}

func NewRouter() (r *Router) {
	return new(Router)
}
