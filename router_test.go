package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	router := NewRouter()

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

	testRequest("GET", "/missing", 404, "404 Not Found")
}
