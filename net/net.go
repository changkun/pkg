// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package net

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// RequestParams ...
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

// HTTPRequest create a HTTP request
func HTTPRequest(url, method string, data []byte, params *RequestParams, response interface{}) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("utils: HTTPRequest error: %v", err)
		}
	}()

	client := http.Client{}
	if params.Timeout < 0 {
		params.Timeout = 0
	}
	client.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, time.Second*time.Duration(params.Timeout))
			if err != nil {
				return nil, err
			}
			c.SetDeadline(time.Now().Add(time.Second * time.Duration(params.Timeout)))
			return c, nil
		},
		DisableKeepAlives: true,
	}
	request, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewReader(data))
	if err != nil {
		return
	}
	request.SetBasicAuth(params.AuthUser, params.AuthPass)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer func() {
		err = resp.Body.Close()
	}()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(respBytes, response)
	return
}
