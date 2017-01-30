package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebHandler(t *testing.T) {
	cases := []struct{ in, out string }{
		{"asdf@golang.org", "works!"},
		{"foobar", "works!"},
	}

	for _, c := range cases {
		req, err := http.NewRequest(http.MethodGet,
			"http://localhost:8080/"+c.in,
			nil,
		)

		if err != nil {
			t.Fatalf("Unable to create http.NewRequest: %v", err)
		}

		rec := httptest.NewRecorder()

		MainHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("got status: %d", rec.Code)
		}

		if strings.Contains(rec.Body.String(), c.out) == false {
			t.Fatalf("String not found")
		}
	}
}
