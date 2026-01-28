package handlers

import (
	"bytes"
	"crypto/subtle"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"json2xls/internal/converter"
	"json2xls/internal/models"

	"github.com/xuri/excelize/v2"
)

//go:embed openapi.json
var openAPISpecBytes []byte

func basicAuthEnabled() (user string, pass string, enabled bool) {
	user = os.Getenv("BASIC_AUTH_USER")
	pass = os.Getenv("BASIC_AUTH_PASS")
	if user == "" || pass == "" {
		return "", "", false
	}
	return user, pass, true
}

// WithBasicAuth wraps handler with Basic Auth if BASIC_AUTH_USER and BASIC_AUTH_PASS are set.
func WithBasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, enabled := basicAuthEnabled()
		if !enabled {
			next(w, r)
			return
		}

		u, p, ok := r.BasicAuth()
		if !ok ||
			subtle.ConstantTimeCompare([]byte(u), []byte(user)) != 1 ||
			subtle.ConstantTimeCompare([]byte(p), []byte(pass)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="json2xls"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// HandleJson2Xls обрабатывает POST запрос для конвертации JSON в XLS
func HandleJson2Xls(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Конвертируем в XLS
	f, err := converter.ConvertToXLS(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting to XLS: %v", err), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки для скачивания файла (до записи в ResponseWriter)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=output.xlsx")

	// Записываем файл напрямую в ResponseWriter
	if err := f.Write(w); err != nil {
		http.Error(w, fmt.Sprintf("Error writing file: %v", err), http.StatusInternalServerError)
		return
	}
}

// HandleXls2Json обрабатывает POST запрос для конвертации XLS в JSON
func HandleXls2Json(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсим multipart form (максимальный размер файла 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	// Получаем файл из формы
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting file: %v. Please use 'file' as form field name.", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Проверяем расширение файла
	filename := header.Filename
	if !strings.HasSuffix(strings.ToLower(filename), ".xlsx") &&
		!strings.HasSuffix(strings.ToLower(filename), ".xls") {
		http.Error(w, "Invalid file format. Only .xlsx and .xls files are supported.", http.StatusBadRequest)
		return
	}

	// Читаем файл в память
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
		return
	}

	// Открываем Excel файл
	f, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening Excel file: %v", err), http.StatusBadRequest)
		return
	}
	defer f.Close()

	// Конвертируем в JSON
	req, err := converter.ConvertXLSToJSON(f)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовки для JSON ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодируем и отправляем JSON
	if err := json.NewEncoder(w).Encode(req); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

// HandleOpenAPI возвращает OpenAPI спецификацию
func HandleOpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(openAPISpecBytes)
}

// HandleDoc отображает документацию с использованием Redoc
func HandleDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	html := `<!DOCTYPE html>
<html>
  <head>
    <title>JSON to XLS Converter API - Documentation</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <redoc spec-url='/openapi.json'></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
  </body>
</html>`

	w.Write([]byte(html))
}

// HandleHealth проверка здоровья сервиса
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
