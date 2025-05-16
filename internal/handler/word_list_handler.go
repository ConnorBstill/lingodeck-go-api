package handler

import (
	"encoding/json"
	"net/http"

	"lingodeck-go-api/internal/service"
)

type WordListData struct {
	List        []service.Translation `json:"list"`
	AudioSample []byte                `json:"audioSample"`
}

func GetWordListDataHandler(writer http.ResponseWriter, request *http.Request) {
	genaiChan := make(chan []service.Translation)
	ttsChan := make(chan []byte)

	go service.GetRelatedWordList("desert", "fr-CA", 20, genaiChan)
	go service.GetTextToSpeech("desert", ttsChan)

	list := <-genaiChan
	audioSample := <-ttsChan

	wordListData := WordListData{List: list, AudioSample: audioSample}

	writer.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(writer).Encode(wordListData); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
