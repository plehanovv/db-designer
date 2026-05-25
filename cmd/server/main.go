package main

import (
	"log"
	"net/http"
	"os"

	"db-designer-vkr/internal/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fileServer := http.FileServer(http.Dir("web"))
	http.Handle("/", fileServer)
	http.HandleFunc("/analyze", handler.AnalyzeHandler)

	log.Println("Server started on :" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
