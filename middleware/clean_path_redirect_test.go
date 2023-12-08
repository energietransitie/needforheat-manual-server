package middleware

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

func TestCleanPathRedirect(t *testing.T) {
	r := chi.NewRouter()

	r.Use(CleanPathRedirect)

	returnRequestPath := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.Path))
	}

	r.Get("/test", returnRequestPath)
	r.Get("/test/", returnRequestPath)
	r.Get("/test/multiple/", returnRequestPath)

	ts := httptest.NewServer(r)
	defer ts.Close()

	testRequest(t, ts, "GET", "/test", http.StatusOK, "/test/")
	testRequest(t, ts, "GET", "/test/", http.StatusOK, "/test/")
	testRequest(t, ts, "GET", "/test/multiple", http.StatusOK, "/test/multiple/")
	testRequest(t, ts, "GET", "/test/multiple/", http.StatusOK, "/test/multiple/")
	testRequest(t, ts, "GET", "/test//multiple/", http.StatusOK, "/test/multiple/")
	testRequest(t, ts, "GET", "/test//multiple//", http.StatusOK, "/test/multiple/")
}

func testRequest(t *testing.T, ts *httptest.Server, method string, path string, expectedStatusCode int, expectedBody string) {
	t.Run(method+" "+path, func(t *testing.T) {
		req, err := http.NewRequest(method, ts.URL+path, nil)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return errors.New("redirect not allowed")
			},
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != expectedStatusCode {
			t.Fatalf("expected status code %d, got %d", expectedStatusCode, resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()

		if string(body) != expectedBody {
			t.Fatalf("expected body %s, got %s", expectedBody, string(body))
		}
	})
}
