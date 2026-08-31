// Copyright 2021 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

// pollInterval is how long to wait between progress checks while the service
// transcribes the upload.
const pollInterval = 30 * time.Second

// statusDone is the progress code the service reports when a task is finished.
const statusDone = 9

func main() {
	os.Exit(run())
}

// run holds the body so that deferred cleanup happens; os.Exit skips it.
func run() int {
	var (
		appid     = os.Getenv("XFYUN_APPID")
		secretKey = os.Getenv("XFYUN_SECRETKEY")
	)
	if appid == "" || secretKey == "" {
		log.Println("set XFYUN_APPID and XFYUN_SECRETKEY")
		return 1
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	client := New(appid, secretKey)
	taskid, uerr := client.UploadAudio(ctx, "./testdata/test.mp3", "cn")
	if uerr != nil {
		log.Printf("failed to upload the audio: %v", uerr)
		return 1
	}
	log.Printf("audio has uploaded, task id: %v", taskid)

	// The first progress check is part of the same poll as the rest; it was a
	// separate call whose result was logged and then thrown away.
	for {
		status, err := client.GetProgress(ctx, taskid)
		if err != nil {
			log.Printf("failed to get progress: %v", err)
			return 1
		}
		if status == statusDone {
			break
		}
		log.Printf("status: %v, go sleep...", status)

		select {
		case <-ctx.Done():
			log.Printf("interrupted: %v", ctx.Err())
			return 1
		case <-time.After(pollInterval):
		}
	}

	content, err := client.GetResult(ctx, taskid)
	if err != nil {
		log.Printf("failed to get result: %v", err)
		return 1
	}
	log.Println("result:", content)
	return 0
}
