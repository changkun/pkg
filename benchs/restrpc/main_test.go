package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/changkun/gobase/benchs/restrpc/rpcs"
	"github.com/changkun/gobase/benchs/restrpc/ser"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func bootServer(success chan bool) {
	go func() {
		defer func() {
			if m := recover(); m != nil {
				success <- false
			}
		}()
		go ser.RunRPC()
		ser.RunHTTP()
	}()

	timeout := time.Second * 10
	start := time.Now()

	for {
		if time.Now().Sub(start) > timeout {
			logrus.Println("timeout")
			success <- false
			break
		}

		res, err := http.Get("http://localhost:12345/api/v1/ping")
		if err != nil {
			logrus.Printf("1 err: %v", err)
			success <- false
			break
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logrus.Printf("2 err: %v", err)
			success <- false
			break
		}

		var m map[string]string
		if err := json.Unmarshal(body, &m); err != nil {
			logrus.Printf("3 err: %v", err)
			success <- false
			break
		}

		msg, ok := m["msg"]
		if msg != "pong" || !ok {
			logrus.Println("msg is not ok")
			success <- false
			break
		}

		success <- true
		break
	}
}

func init() {
	bootcheck := make(chan bool, 1)
	go bootServer(bootcheck)
	if success := <-bootcheck; success == false {
		logrus.Fatal("fail to boot the service")
	}
}

func BenchmarkAPIRestful(b *testing.B) {
	for i := 0; i < b.N; i++ {
		requestBody, err := json.Marshal(map[string]float64{
			"a": 42.0,
			"b": 99.9,
		})
		if err != nil {
			b.Fatal("prepare body fail")
		}
		body := bytes.NewBuffer(requestBody)
		req, err := http.NewRequest("POST", "http://0.0.0.0:12345/api/v1/add", body)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}
}

func BenchmarkAPIgRPC(b *testing.B) {
	conn, err := grpc.Dial("0.0.0.0:12346", grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("did not connect: \n\t%v", err)
	}
	defer conn.Close()
	client := rpcs.NewArithmeticClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Add(ctx, &rpcs.AddInput{
			A: 42.0,
			B: 99.9,
		})
	}
}
