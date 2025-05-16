package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Translation struct {
	ID          int    `json:"id"`
	Translation string `json:"translation"`
	Word        string `json:"word"`
}

func GetRelatedWordList(term string, languageCode string, numOfResults int16, c chan []Translation) {
	ctx := context.Background()

	genaiApiKey := os.Getenv("GEMINI_API_KEY")
	geminiClientOption := option.WithAPIKey(genaiApiKey)

	client, err := genai.NewClient(ctx, geminiClientOption)
	if err != nil {
		fmt.Println("Error getting AI client", err)
		defer client.Close()

		return
	}

	model := client.GenerativeModel("gemini-1.5-flash")

	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"id": {
					Type: genai.TypeInteger,
				},
				"word": {
					Type: genai.TypeString,
				},
				"translation": {
					Type: genai.TypeString,
				},
			},
			Required: []string{"id", "word", "translation"},
		},
	}

	prompt := fmt.Sprintf(
		`Create a list of %d words related to "%s" in all lowercase, with no proper nouns, along with their translations in %s, formatted as an array of JSON objects with the properties 'word' for the word and 'translation' for the word's translation, as well as a property called 'id' that is equal to the object's index + 1. The output should conform to the given JSON schema.`,
		numOfResults,
		term,
		languageCode,
	)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		fmt.Println("Error getting related words", err)
		defer client.Close()

		return
	}

	var respJsonString string
	if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		respJsonString = string(textPart)
	} else {
		fmt.Println("Error: Expected a Text part")
		return
	}

	var translationsList []Translation

	unmarshalErr := json.Unmarshal([]byte(respJsonString), &translationsList)
	if unmarshalErr != nil {
		fmt.Println("Error parsing json", err)

		return
	}

	c <- translationsList
}
