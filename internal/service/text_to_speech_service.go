package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type ElevenLabsPayload struct {
	Text    string `json:"text"`
	ModelId string `json:"model_id"`
}

func GetTextToSpeech(text string, ttsChan chan []byte) {
	url := "https://api.elevenlabs.io/v1/text-to-speech/a5n9pJUnAhX4fn7lx3uo?output_format=mp3_44100_128"

	genaiApiKey := os.Getenv("GEMINI_API_KEY")

	payload := &ElevenLabsPayload{Text: "The first move is what sets everything in motion.", ModelId: "eleven_flash_v2_5"}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	payloadReader := strings.NewReader(string(payloadJson))

	req, _ := http.NewRequest("POST", url, payloadReader)

	req.Header.Add("xi-api-key", genaiApiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)

	ttsChan <- resBody
}
