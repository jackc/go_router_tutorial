package router

import (
	"io"
	"net/http"
	"strings"
)

type Router struct {
	endpoint        http.Handler
	staticBranches  map[string]*Router
	parameterBranch *Router
}

// ServeHTTP makes Router implement standard http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	segments := segmentizePath(req.URL.Path)
	if endpoint, ok := r.findEndpoint(segments); ok {
		endpoint.ServeHTTP(w, req)
	} else {
		w.WriteHeader(404)
		io.WriteString(w, "404 Not Found")
	}
}

func (r *Router) AddRoute(method string, path string, handler http.Handler) {
	segments := segmentizePath(path)
	endpoint := handler
	r.addRouteFromSegments(method, segments, endpoint)
}

func NewRouter() (r *Router) {
	r = new(Router)
	r.staticBranches = make(map[string]*Router)
	return r
}

func (r *Router) addRouteFromSegments(method string, segments []string, endpoint http.Handler) {
	if len(segments) > 0 {
		head, tail := segments[0], segments[1:]

		var subrouter *Router
		if strings.HasPrefix(head, ":") {
			if r.parameterBranch == nil {
				r.parameterBranch = NewRouter()
			}
			subrouter = r.parameterBranch
		} else {
			if _, present := r.staticBranches[head]; !present {
				r.staticBranches[head] = NewRouter()
			}
			subrouter = r.staticBranches[head]
		}

		subrouter.addRouteFromSegments(method, tail, endpoint)
	} else {
		r.endpoint = endpoint
	}
}

func segmentizePath(path string) (segments []string) {
	for _, s := range strings.Split(path, "/") {
		if len(s) != 0 {
			segments = append(segments, s)
		}
	}
	return
}

func (r *Router) findEndpoint(segments []string) (http.Handler, bool) {
	if len(segments) > 0 {
		head, tail := segments[0], segments[1:]
		if subrouter, present := r.staticBranches[head]; present {
			return subrouter.findEndpoint(tail)
		} else if r.parameterBranch != nil {
			return r.parameterBranch.findEndpoint(tail)
		} else {
			return nil, false
		}
	}
	return r.endpoint, true
}
