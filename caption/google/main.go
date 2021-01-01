// Speech to text
// https://cloud.google.com/speech-to-text/docs/async-recognize
// https://cloud.google.com/speech-to-text/docs/apis
// https://pkg.go.dev/google.golang.org/genproto@v0.0.0-20201214200347-8c77b98c765d/googleapis/cloud/speech/v1
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	speech "cloud.google.com/go/speech/apiv1p1beta1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1p1beta1"
)

func main() {
	ctx := context.Background()

	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	gsf := "gs://changkun.de/test2.mp3"

	f, err := os.Create("out.txt")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer f.Close()

	err = send(f, client, gsf)
	if err != nil {
		log.Fatalf("failed to send file to googld: %v", err)
	}
}

func send(w io.Writer, client *speech.Client, gsf string) error {
	ctx := context.Background()
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
	err = os.WriteFile("all.txt", b, os.ModePerm)
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
