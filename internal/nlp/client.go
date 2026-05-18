package nlp

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type AnalyzeRequest struct {
	Text string `json:"text"`
}

func ProcessText(text string) Document {

	requestBody := AnalyzeRequest{
		Text: text,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return Document{}
	}

	response, err := http.Post(
		"http://localhost:8000/analyze",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return Document{}
	}

	defer response.Body.Close()

	var document Document

	err = json.NewDecoder(response.Body).Decode(&document)
	if err != nil {
		return Document{}
	}

	return document
}
