// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package net

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// RequestParams carries the per-request settings HTTPRequest cannot infer
// from the URL: the timeout in seconds and optional basic auth credentials.
type RequestParams struct {
	Timeout  int
	AuthUser string
	AuthPass string
}

// QueryEncoder encodes a key value map to URL query string
func QueryEncoder(m map[string]string) string {
	query := url.Values{}
	for k, v := range m {
		query.Set(k, v)
	}
	return query.Encode()
}

// HTTPRequest performs an HTTP request and decodes the JSON body into
// response. Cancelling ctx aborts the request; params.Timeout additionally
// bounds the dial and the connection as a whole.
func HTTPRequest(ctx context.Context, url, method string, data []byte, params *RequestParams, response any) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("utils: HTTPRequest error: %w", err)
		}
	}()

	timeout := time.Duration(max(params.Timeout, 0)) * time.Second
	dialer := &net.Dialer{Timeout: timeout}
	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				c, derr := dialer.DialContext(ctx, network, addr)
				if derr != nil {
					return nil, derr
				}
				if timeout > 0 {
					if derr := c.SetDeadline(time.Now().Add(timeout)); derr != nil {
						c.Close()
						return nil, derr
					}
				}
				return c, nil
			},
			DisableKeepAlives: true,
		},
	}

	request, err := http.NewRequestWithContext(ctx, strings.ToUpper(method), url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	request.SetBasicAuth(params.AuthUser, params.AuthPass)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	// Report a close failure only when the request itself succeeded, so a
	// real error is not replaced by a secondary one.
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(respBytes, response)
}
