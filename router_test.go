package router

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	router := NewRouter()

	widgetIndexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "widgetIndex")
	})
	router.AddRoute("GET", "/widgets", widgetIndexHandler)

	testRequest := func(method string, path string, expectedCode int, expectedBody string) {
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

	testRequest("GET", "/widgets", 200, "widgetIndex")
	testRequest("GET", "/missing", 404, "404 Not Found")
}
