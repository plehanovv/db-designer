package nlp

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type AnalyzeRequest struct {
	Text string `json:"text"`
}

type Token struct {
	Text       string `json:"text"`
	Lemma      string `json:"lemma"`
	Pos        string `json:"pos"`
	Dependency string `json:"dependency"`
	Head       string `json:"head"`
}

type AnalyzeResponse struct {
	Tokens []Token `json:"tokens"`
}

func AnalyzeText(text string) (*AnalyzeResponse, error) {

	requestBody := AnalyzeRequest{
		Text: text,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(
		"http://localhost:8000/analyze",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var result AnalyzeResponse

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
