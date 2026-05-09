package main

import (
	"log"
	"net/http"

	"db-designer-vkr/internal/handler"
)

func main() {
	http.HandleFunc("/analyze", handler.AnalyzeHandler)

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
