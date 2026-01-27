package main

import (
	"log"
	"net/http"

	"json2xls/internal/handlers"
)

func main() {
	http.HandleFunc("/json2xls", handlers.HandleJson2Xls)
	http.HandleFunc("/xls2json", handlers.HandleXls2Json)
	http.HandleFunc("/doc", handlers.HandleDoc)
	http.HandleFunc("/openapi.json", handlers.HandleOpenAPI)
	http.HandleFunc("/health", handlers.HandleHealth)

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Printf("POST /json2xls - Convert JSON to XLS")
	log.Printf("POST /xls2json - Convert XLS/XLSX to JSON")
	log.Printf("GET /doc - API Documentation (Redoc)")
	log.Printf("GET /openapi.json - OpenAPI specification")
	log.Printf("GET /health - Health check")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

