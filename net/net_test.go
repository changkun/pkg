// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package net_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"changkun.de/x/pkg/net"
)

func TestQueryEncoder(t *testing.T) {
	want := "Test=http%3A%2F%2Fchangkun.test"
	got := net.QueryEncoder(map[string]string{
		"Test": "http://changkun.test",
	})
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %v, got: %v", want, got)
	}
}

// newServer returns a server that echoes the method and the basic auth
// credentials it was called with, so a test can check what reached it.
func newServer(t *testing.T) *httptest.Server {
	t.Helper()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		fmt.Fprintf(w, `{"method":%q,"user":%q,"pass":%q}`, r.Method, user, pass)
	}))
	t.Cleanup(s.Close)
	return s
}

type echo struct {
	Method string `json:"method"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
}

func TestHTTPRequest(t *testing.T) {
	s := newServer(t)

	var got echo
	err := net.HTTPRequest(t.Context(), s.URL, "get", nil,
		&net.RequestParams{Timeout: 10, AuthUser: "user", AuthPass: "pass"}, &got)
	if err != nil {
		t.Fatalf("HTTPRequest: %v", err)
	}

	want := echo{Method: http.MethodGet, User: "user", Pass: "pass"}
	if got != want {
		t.Errorf("server saw %+v, want %+v", got, want)
	}
}

// TestHTTPRequestNegativeTimeout guards the clamp: a negative timeout used to
// be passed straight to the dialer as a negative duration.
func TestHTTPRequestNegativeTimeout(t *testing.T) {
	s := newServer(t)

	var got echo
	err := net.HTTPRequest(t.Context(), s.URL, http.MethodGet, nil,
		&net.RequestParams{Timeout: -100}, &got)
	if err != nil {
		t.Fatalf("HTTPRequest with a negative timeout: %v", err)
	}
	if got.Method != http.MethodGet {
		t.Errorf("server saw method %q, want %q", got.Method, http.MethodGet)
	}
}

func TestHTTPRequestNonJSONBody(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "<html>not json</html>")
	}))
	defer s.Close()

	if err := net.HTTPRequest(t.Context(), s.URL, http.MethodGet, nil,
		&net.RequestParams{Timeout: 10}, &struct{}{}); err == nil {
		t.Fatal("decoding a non-JSON body returned nil error")
	}
}

// TestHTTPRequestCancelled is the reason the context parameter exists: a
// caller must be able to abandon a request that the server is sitting on.
func TestHTTPRequestCancelled(t *testing.T) {
	release := make(chan struct{})
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-release
	}))
	defer s.Close()
	defer close(release)

	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()

	err := net.HTTPRequest(ctx, s.URL, http.MethodGet, nil,
		&net.RequestParams{Timeout: 10}, &struct{}{})
	if err == nil {
		t.Fatal("a cancelled request returned nil error")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("error = %v, want a context deadline", err)
	}
}
