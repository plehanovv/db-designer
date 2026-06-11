package main

import (
	"log"
	"net/http"
	"os"

	"db-designer-vkr/internal/handler"
	"db-designer-vkr/internal/storage"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fileServer := http.FileServer(http.Dir("web"))
	http.Handle("/", noCache(fileServer))
	http.HandleFunc("/analyze", handler.AnalyzeHandler)
	http.HandleFunc("/generate-sql", handler.GenerateSQLHandler)

	if store, enabled := storage.NewFromEnv(); enabled {
		if err := store.Init(); err != nil {
			log.Printf("PostgreSQL result storage is disabled: %v", err)
		} else {
			handler.SetStore(store)
			log.Println("PostgreSQL result storage is enabled")
		}
	}

	log.Println("Server started on :" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}
