// Copyright 2021 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Speech to text
// https://cloud.google.com/speech-to-text/docs/async-recognize
// https://cloud.google.com/speech-to-text/docs/apis
// https://pkg.go.dev/cloud.google.com/go/speech/apiv1p1beta1/speechpb
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	speech "cloud.google.com/go/speech/apiv1p1beta1"
	speechpb "cloud.google.com/go/speech/apiv1p1beta1/speechpb"
)

func main() {
	os.Exit(run())
}

// run holds the body so that the deferred close happens; os.Exit, which
// log.Fatalf calls, skips deferred functions.
func run() int {
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Printf("failed to create client: %v", err)
		return 1
	}
	defer client.Close()

	gsf := "gs://changkun.de/test2.mp3"

	f, err := os.Create("out.txt")
	if err != nil {
		log.Printf("failed to create output file: %v", err)
		return 1
	}
	defer f.Close()

	if err := send(ctx, f, client, gsf); err != nil {
		log.Printf("failed to send the file to google: %v", err)
		return 1
	}
	return 0
}

func send(ctx context.Context, w io.Writer, client *speech.Client, gsf string) error {
	req := &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_MP3,
			SampleRateHertz: 44100,
			LanguageCode:    "zh",
			DiarizationConfig: &speechpb.SpeakerDiarizationConfig{
				EnableSpeakerDiarization: true,
				MinSpeakerCount:          2,
				MaxSpeakerCount:          4,
			},
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: gsf},
		},
	}
	log.Printf("long running recognize...")
	op, err := client.LongRunningRecognize(ctx, req)
	if err != nil {
		return err
	}

	log.Printf("long running recognize... op.wait")
	resp, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	log.Printf("long running recognize...")
	b, err := json.Marshal(resp.Results)
	if err != nil {
		return err
	}
	err = os.WriteFile("all.txt", b, 0o600)
	if err != nil {
		return err
	}

	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Fprintf(w, "\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}
	return nil
}
