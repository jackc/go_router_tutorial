package go_router_tutorial

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func stubHandler(responseBody string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, responseBody)
		for key, values := range r.URL.Query() {
			fmt.Fprintf(w, " %s: %s", key, values[0])
		}
	})
}

func testRequest(t *testing.T, router *Router, method string, path string, expectedCode int, expectedBody string) {
	response := httptest.NewRecorder()
	request, err := http.NewRequest(method, "http://example.com"+path, nil)
	if err != nil {
		t.Errorf("Unable to create test %s request for %s", method, path)
	}

	router.ServeHTTP(response, request)
	if response.Code != expectedCode {
		t.Errorf("%s %s: expected HTTP code %d, received %d", method, path, expectedCode, response.Code)
	}
	if response.Body.String() != expectedBody {
		t.Errorf("%s %s: expected HTTP response body \"%s\", received \"%s\"", method, path, expectedBody, response.Body.String())
	}
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	r.Get("/widgets", stubHandler("widgetIndex"))
	r.Post("/widgets", stubHandler("widgetCreate"))
	r.Get("/widgets/:id", stubHandler("widgetShow"))
	r.Get("/widgets/:id/edit", stubHandler("widgetEdit"))
	r.Put("/widgets/:id", stubHandler("widgetUpdate"))
	r.Delete("/widgets/:id", stubHandler("widgetDelete"))

	testRequest(t, r, "GET", "/widgets", 200, "widgetIndex")
	testRequest(t, r, "POST", "/widgets", 200, "widgetCreate")
	testRequest(t, r, "GET", "/widgets/1", 200, "widgetShow id: 1")
	testRequest(t, r, "GET", "/widgets/2", 200, "widgetShow id: 2")
	testRequest(t, r, "GET", "/widgets/1/edit", 200, "widgetEdit id: 1")
	testRequest(t, r, "PUT", "/widgets/1", 200, "widgetUpdate id: 1")
	testRequest(t, r, "DELETE", "/widgets/1", 200, "widgetDelete id: 1")
	testRequest(t, r, "GET", "/missing", 404, "404 Not Found")
}

var benchmarkRouter *Router

func getBenchmarkRouter() *Router {
	if benchmarkRouter != nil {
		return benchmarkRouter
	}

	benchmarkRouter := NewRouter()
	handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	benchmarkRouter.AddRoute("GET", "/", handler)
	benchmarkRouter.AddRoute("GET", "/foo", handler)
	benchmarkRouter.AddRoute("GET", "/foo/bar", handler)
	benchmarkRouter.AddRoute("GET", "/foo/baz", handler)
	benchmarkRouter.AddRoute("GET", "/foo/bar/baz/quz", handler)
	benchmarkRouter.AddRoute("GET", "/people", handler)
	benchmarkRouter.AddRoute("GET", "/people/search", handler)
	benchmarkRouter.AddRoute("GET", "/people/:id", handler)
	benchmarkRouter.AddRoute("GET", "/users", handler)
	benchmarkRouter.AddRoute("GET", "/users/:id", handler)
	benchmarkRouter.AddRoute("GET", "/widgets", handler)
	benchmarkRouter.AddRoute("GET", "/widgets/important", handler)

	return benchmarkRouter
}

func getBench(b *testing.B, handler http.Handler, path string, expectedCode int) {
	response := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "http://example.com"+path, nil)
	if err != nil {
		b.Fatalf("Unable to create test GET request for %v", path)
	}

	handler.ServeHTTP(response, request)
	if response.Code != expectedCode {
		b.Fatalf("GET %v: expected HTTP code %v, received %v", path, expectedCode, response.Code)
	}
}

func BenchmarkRoutedRequest(b *testing.B) {
	router := getBenchmarkRouter()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		getBench(b, router, "/widgets/important", 200)
	}
}

func BenchmarkFindEndpointRoot(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.findEndpoint("GET", segmentizePath("/"), []string{})
	}
}

func BenchmarkFindEndpointSegment1(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.findEndpoint("GET", segmentizePath("/foo"), []string{})
	}
}

func BenchmarkFindEndpointSegment2(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.findEndpoint("GET", segmentizePath("/people/search"), []string{})
	}
}

func BenchmarkFindEndpointSegment2Placeholder(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.findEndpoint("GET", segmentizePath("/people/1"), []string{})
	}
}

func BenchmarkFindEndpointSegment4(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.findEndpoint("GET", segmentizePath("/foo/bar/baz/quz"), []string{})
	}
}
