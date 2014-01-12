package router

import (
	"io"
	"net/http"
	"strings"
)

type endpoint struct {
	handler    http.Handler
	parameters []string
}

type Router struct {
	endpoint        *endpoint
	staticBranches  map[string]*Router
	parameterBranch *Router
}

// ServeHTTP makes Router implement standard http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	segments := segmentizePath(req.URL.Path)
	if endpoint, arguments, ok := r.findEndpoint(segments, []string{}); ok {
		addRouteArgumentsToRequest(endpoint.parameters, arguments, req)
		endpoint.handler.ServeHTTP(w, req)
	} else {
		w.WriteHeader(404)
		io.WriteString(w, "404 Not Found")
	}
}

func (r *Router) AddRoute(method string, path string, handler http.Handler) {
	segments := segmentizePath(path)
	parameters := extractParameterNames(segments)
	endpoint := &endpoint{handler: handler, parameters: parameters}
	r.addRouteFromSegments(method, segments, endpoint)
}

func NewRouter() (r *Router) {
	r = new(Router)
	r.staticBranches = make(map[string]*Router)
	return r
}

func (r *Router) addRouteFromSegments(method string, segments []string, endpoint *endpoint) {
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

func (r *Router) findEndpoint(segments []string, pathArguments []string) (*endpoint, []string, bool) {
	if len(segments) > 0 {
		head, tail := segments[0], segments[1:]
		if subrouter, present := r.staticBranches[head]; present {
			return subrouter.findEndpoint(tail, pathArguments)
		} else if r.parameterBranch != nil {
			pathArguments = append(pathArguments, head)
			return r.parameterBranch.findEndpoint(tail, pathArguments)
		} else {
			return nil, nil, false
		}
	}
	return r.endpoint, pathArguments, true
}

func addRouteArgumentsToRequest(names, values []string, req *http.Request) {
	query := req.URL.Query()
	for i := 0; i < len(names); i++ {
		query.Set(names[i], values[i])
	}
	req.URL.RawQuery = query.Encode()
}

func extractParameterNames(segments []string) (parameters []string) {
	for _, s := range segments {
		if strings.HasPrefix(s, ":") {
			parameters = append(parameters, s[1:])
		}
	}
	return
}
