package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "localhost:8000", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	rec := httptest.NewRecorder()

	index(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}
}

func TestFollowHandler(t *testing.T) {
	tt := []struct {
		name   string
		urlID  string
		status int
	}{
		{name: "unknown url id", urlID: "0", status: http.StatusNotFound},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "localhost:8000/f/"+tc.urlID, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			req = mux.SetURLVars(req, map[string]string{"urlID": tc.urlID})
			rec := httptest.NewRecorder()

			follow(rec, req)

			res := rec.Result()
			if res.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, res.StatusCode)
			}
		})
	}
}
