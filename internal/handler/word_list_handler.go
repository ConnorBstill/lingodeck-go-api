package handler

import (
	"encoding/json"
	"net/http"

	"lingodeck-go-api/internal/service"
)

// HelloHandler responds with a simple greeting.
func GetWordListDataHandler(writer http.ResponseWriter, request *http.Request) {
	list := service.GenerateRelatedWordList("desert", "fr-CA", 10)

	writer.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(writer).Encode(list.Parts); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
