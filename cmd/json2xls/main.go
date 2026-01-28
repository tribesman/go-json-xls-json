package main

import (
	"log"
	"net/http"
	"os"

	"json2xls/internal/handlers"
)

func main() {
	// Защищаем BasicAuth'ом API (если заданы BASIC_AUTH_USER/BASIC_AUTH_PASS).
	// Документацию и health оставляем открытыми.
	http.HandleFunc("/json2xls", handlers.WithBasicAuth(handlers.HandleJson2Xls))
	http.HandleFunc("/xls2json", handlers.WithBasicAuth(handlers.HandleXls2Json))
	http.HandleFunc("/openapi.json", handlers.WithBasicAuth(handlers.HandleOpenAPI))
	http.HandleFunc("/doc", handlers.HandleDoc)
	http.HandleFunc("/health", handlers.HandleHealth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("Server starting on port %s", addr)
	log.Printf("POST /json2xls - Convert JSON to XLS")
	log.Printf("POST /xls2json - Convert XLS/XLSX to JSON")
	log.Printf("GET /doc - API Documentation (Redoc)")
	log.Printf("GET /openapi.json - OpenAPI specification")
	log.Printf("GET /health - Health check")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
