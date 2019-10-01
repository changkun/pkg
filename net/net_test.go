// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package net_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/changkun/gobase/net"
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

func TestHTTPRequest(t *testing.T) {
	got := net.HTTPRequest("https://google.de/", "GET", []byte{}, &net.RequestParams{
		Timeout:  -100,
		AuthUser: "test",
		AuthPass: "test",
	}, &struct{}{})
	if errors.Is(nil, got) {
		t.Fatalf("want error, got: %v", got)
	}
}
