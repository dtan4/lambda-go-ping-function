package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRealHandler_success(t *testing.T) {
	Stdout = new(bytes.Buffer)
	HTTPClient = http.DefaultClient

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "192.168.0.1")
	}))
	defer ts.Close()

	PingURL = ts.URL

	ctx := context.Background()

	if err := RealHandler(ctx); err != nil {
		t.Fatalf("want: no error, got: %q", err)
	}
}

func TestRealHandler_error(t *testing.T) {
	Stdout = new(bytes.Buffer)
	HTTPClient = http.DefaultClient

	testcases := []struct {
		subject    string
		statusCode int
		body       string
		error      string
	}{
		{
			subject:    "response_400",
			statusCode: http.StatusBadRequest,
			body:       "bad request",
			error:      `unexpected response status: 400, body: "bad request\n"`,
		},
		{
			subject:    "response_502",
			statusCode: http.StatusBadGateway,
			body:       "bad gateway",
			error:      `unexpected response status: 502, body: "bad gateway\n"`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.subject, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				fmt.Fprintln(w, tc.body)
			}))
			defer ts.Close()

			PingURL = ts.URL

			ctx := context.Background()

			err := RealHandler(ctx)

			if err == nil {
				t.Fatal("want: error, got: nil")
			}

			if err.Error() != tc.error {
				t.Fatalf("want: %q, got: %q", tc.error, err.Error())
			}
		})
	}
}
