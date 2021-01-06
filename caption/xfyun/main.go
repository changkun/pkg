// Copyright 2021 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"time"
)

func main() {
	var (
		appid     = os.Getenv("XFYUN_APPID")
		secretKey = os.Getenv("XFYUN_SECRETKEY")
	)
	if len(appid) == 0 || len(secretKey) == 0 {
		panic("cannot read configuration")
	}

	client := New(appid, secretKey)
	taskid, err := client.UploadAudio("./testdata/test.mp3", "cn")
	if err != nil {
		log.Fatalf("failed to upload the audio: %v", err)
	}
	log.Printf("audio has uploaded, task id: %v", taskid)

	status, err := client.GetProgress(taskid)
	if err != nil {
		log.Fatalf("failed to get progress: %v", err)
	}
	log.Println("status:", status)

	for {
		status, err := client.GetProgress(taskid)
		if err != nil {
			log.Fatalf("failed to get progress: %v", err)
		}
		if status == 9 {
			break
		}
		log.Printf("status: %v, go sleep...", status)
		time.Sleep(30 * time.Second)
	}
	content, err := client.GetResult(taskid)
	if err != nil {
		log.Fatalf("failed to get result: %v", err)
	}
	log.Println("result:", content)
}
