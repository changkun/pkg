// Copyright 2020 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"changkun.de/x/pkg/benchs/restrpc/rpcs"
	"changkun.de/x/pkg/benchs/restrpc/ser"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpAddr    = "http://0.0.0.0:12345"
	grpcAddr    = "0.0.0.0:12346"
	bootTimeout = 10 * time.Second
	bootRetry   = 50 * time.Millisecond
)

// bootServer starts both servers and polls the REST health endpoint until it
// answers. The endpoint is not up the instant Listen returns, so a single
// probe races the server and the poll has to retry rather than give up on the
// first refused connection.
func bootServer(ctx context.Context) error {
	go func() {
		if err := ser.RunRPC(); err != nil {
			log.Printf("grpc server: %v", err)
		}
	}()
	go func() {
		if err := ser.RunHTTP(); err != nil {
			log.Printf("http server: %v", err)
		}
	}()

	deadline := time.Now().Add(bootTimeout)
	var last error
	for time.Now().Before(deadline) {
		if last = ping(ctx); last == nil {
			return nil
		}
		time.Sleep(bootRetry)
	}
	return fmt.Errorf("server not ready after %v: %w", bootTimeout, last)
}

// ping reports whether the REST health endpoint answers correctly.
func ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, httpAddr+"/api/v1/ping", nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var m map[string]string
	if err := json.Unmarshal(body, &m); err != nil {
		return fmt.Errorf("decode %q: %w", body, err)
	}
	if m["msg"] != "pong" {
		return fmt.Errorf("ping answered %q, want pong", m["msg"])
	}
	return nil
}

func TestMain(m *testing.M) {
	if err := bootServer(context.Background()); err != nil {
		log.Fatalf("fail to boot the service: %v", err)
	}
	os.Exit(m.Run())
}

// addREST posts an addition and returns the sum the server reports.
func addREST(ctx context.Context, a, b float64) (float64, error) {
	requestBody, err := json.Marshal(map[string]float64{"a": a, "b": b})
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, httpAddr+"/api/v1/add", bytes.NewReader(requestBody))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var out map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return 0, err
	}
	return out["sum"], nil
}

// dialGRPC opens a client connection to the gRPC endpoint.
func dialGRPC(tb testing.TB) rpcs.ArithmeticClient {
	tb.Helper()
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		tb.Fatalf("did not connect: %v", err)
	}
	tb.Cleanup(func() { conn.Close() })
	return rpcs.NewArithmeticClient(conn)
}

// TestBothTransportsAgree guards the benchmarks: a transport that answers
// wrongly, or not at all, would otherwise still produce a respectable number.
func TestBothTransportsAgree(t *testing.T) {
	const a, b = 42.0, 99.5
	const want = a + b

	got, err := addREST(t.Context(), a, b)
	if err != nil {
		t.Fatalf("REST add: %v", err)
	}
	if got != want {
		t.Errorf("REST sum = %v, want %v", got, want)
	}

	out, err := dialGRPC(t).Add(t.Context(), &rpcs.AddInput{A: a, B: b})
	if err != nil {
		t.Fatalf("gRPC add: %v", err)
	}
	if float64(out.GetSum()) != want {
		t.Errorf("gRPC sum = %v, want %v", out.GetSum(), want)
	}
}

func BenchmarkAPIRestful(b *testing.B) {
	ctx := b.Context()
	for b.Loop() {
		if _, err := addREST(ctx, 42.0, 99.9); err != nil {
			b.Fatalf("REST add: %v", err)
		}
	}
}

func BenchmarkAPIgRPC(b *testing.B) {
	client := dialGRPC(b)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	b.ResetTimer()
	for b.Loop() {
		if _, err := client.Add(ctx, &rpcs.AddInput{A: 42.0, B: 99.9}); err != nil {
			b.Fatalf("gRPC add: %v", err)
		}
	}
}
