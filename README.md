# JSON to XLS Converter

REST API сервис на Go для конвертации JSON данных в Excel файлы (XLSX) с поддержкой форматирования ячеек.

## Описание

Приложение принимает JSON данные через REST API и возвращает Excel файл (.xlsx) с поддержкой:

- Множественных листов
- Форматирования ячеек (цвета, шрифты, стили)
- Различных типов данных (числа, строки, булевы значения)

## Требования

- Go 1.21 или выше
- Доступ в интернет для загрузки зависимостей

## Быстрый запуск через Docker (без клонирования репозитория)

### Вариант A (рекомендуется): скачать `docker-compose.yml` и запустить

1. Скачайте `docker-compose.yml` из репозитория:

```bash
curl -fsSL -o docker-compose.yml https://raw.githubusercontent.com/tribesman/go-json-xls-json/main/docker-compose.yml
```

2. (Опционально) Настройте переменные окружения перед запуском:

- **PORT**: порт, на котором сервис слушает внутри контейнера (по умолчанию `8080`)
- **SERVICE_PORT**: порт, на который пробрасываем наружу (по умолчанию `8080`)
- **BASIC_AUTH_USER / BASIC_AUTH_PASS**: если оба заданы — включается Basic Auth для `/json2xls`, `/xls2json`, `/openapi.json`

Пример:

```bash
export PORT=8080
export SERVICE_PORT=9090
export BASIC_AUTH_USER=admin
export BASIC_AUTH_PASS=secret
```

3. Запустите:

```bash
docker compose up -d
```

Образ будет скачан из GHCR: `ghcr.io/tribesman/go-json-xls-json:latest`

### Вариант B: одной командой (без сохранения файла)

```bash
curl -fsSL https://raw.githubusercontent.com/tribesman/go-json-xls-json/main/docker-compose.yml | docker compose -f - up -d
```

## Сборка проекта

### Настройка окружения

Если вы получаете ошибку `go: modules disabled by GO111MODULE=off`, включите поддержку модулей:

```bash
# Для текущей сессии
export GO111MODULE=on

# Или для постоянной настройки (добавьте в ~/.zshrc или ~/.bashrc)
echo 'export GO111MODULE=on' >> ~/.zshrc
source ~/.zshrc
```

Проверьте настройку:

```bash
go env GO111MODULE
```

### Для текущей платформы

**Быстрая сборка (используя скрипт):**

```bash
# Клонируйте или перейдите в директорию проекта
cd json2xls

# Запустите скрипт сборки
./build.sh

# Запустите приложение
./json2xls
```

**Ручная сборка:**

```bash
# Клонируйте или перейдите в директорию проекта
cd json2xls

# Убедитесь, что модули включены
export GO111MODULE=on

# Загрузите зависимости
go mod download

# Или используйте go mod tidy для синхронизации зависимостей
go mod tidy

# Соберите приложение
go build -o json2xls

# Запустите приложение
./json2xls
```

### Для других платформ

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o json2xls-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o json2xls.exe

# macOS (ARM)
GOOS=darwin GOARCH=arm64 go build -o json2xls-macos-arm64

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o json2xls-macos-amd64
```

## Запуск

После сборки запустите приложение:

```bash
./json2xls
```

Сервер запустится на порту `8080` по умолчанию.

## API Endpoints

### POST /json2xls

Конвертирует JSON в Excel файл.

**Запрос:**

- Method: `POST`
- Content-Type: `application/json`
- Body: JSON объект с данными для конвертации

**Ответ:**

- Content-Type: `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
- Body: Excel файл (.xlsx)
- Headers: `Content-Disposition: attachment; filename=output.xlsx`

**Пример запроса:**

```bash
curl -X POST http://localhost:8080/json2xls \
  -H "Content-Type: application/json" \
  -d @request.json \
  --output output.xlsx
```

### POST /xls2json

Конвертирует Excel файл (XLS/XLSX) в JSON.

**Запрос:**

- Method: `POST`
- Content-Type: `multipart/form-data`
- Body: Форма с полем `file`, содержащим Excel файл

**Ответ:**

- Content-Type: `application/json`
- Body: JSON объект в том же формате, что используется для `/json2xls`

**Пример запроса:**

```bash
curl -X POST http://localhost:8080/xls2json \
  -F "file=@input.xlsx" \
  -o output.json
```

**Пример запроса с использованием JavaScript (fetch):**

```javascript
const formData = new FormData();
formData.append("file", fileInput.files[0]);

fetch("http://localhost:8080/xls2json", {
  method: "POST",
  body: formData,
})
  .then((response) => response.json())
  .then((data) => console.log(data));
```

**Примечания:**

- Поддерживаются файлы формата `.xlsx` и `.xls`
- Максимальный размер файла: 10MB
- Извлекаются все листы, ячейки и их форматирование
- Типы данных определяются автоматически (числа, строки, булевы значения)

### GET /doc

Интерактивная документация API с использованием Redoc.

**Запрос:**

- Method: `GET`

**Ответ:**

- Content-Type: `text/html`
- Body: HTML страница с интерактивной документацией

**Использование:**
Откройте в браузере: `http://localhost:8080/doc`

### GET /openapi.json

OpenAPI спецификация API в формате JSON.

**Запрос:**

- Method: `GET`

**Ответ:**

- Content-Type: `application/json`
- Body: OpenAPI 3.0 спецификация

**Пример запроса:**

```bash
curl http://localhost:8080/openapi.json
```

### GET /health

Проверка работоспособности сервиса.

**Запрос:**

- Method: `GET`

**Ответ:**

- Status: `200 OK`
- Content-Type: `application/json`
- Body: `{"status":"ok"}`

**Пример запроса:**

```bash
curl http://localhost:8080/health
```

## Формат данных

### Структура JSON запроса

```json
{
  "sheets": [
    {
      "title": "Sheet1",
      "rows": [
        [
          {
            "value": 1,
            "format": {
              "bold": true,
              "bg_color": "#FFFF00"
            }
          },
          {
            "value": "Product A",
            "format": {
              "font_name": "Arial",
              "font_size": 12,
              "bold": true
            }
          },
          {
            "value": 100.5,
            "format": {
              "text_color": "green",
              "italic": true
            }
          },
          {
            "value": true,
            "format": {
              "bold": true,
              "text_color": "blue"
            }
          }
        ],
        [
          {
            "value": 2,
            "format": {
              "bg_color": "#F0F0F0"
            }
          },
          {
            "value": "Product B"
          },
          {
            "value": 200.75
          },
          {
            "value": false
          }
        ]
      ]
    }
  ]
}
```

### Описание полей

#### Корневой объект

- `sheets` (массив, обязательный) - массив листов Excel

#### Sheet (лист)

- `title` (строка, опционально) - название листа. Если не указано, будет использовано имя по умолчанию (Sheet1, Sheet2, и т.д.)
- `rows` (массив, обязательный) - массив строк

#### Row (строка)

- Массив ячеек (`Cell`)

#### Cell (ячейка)

- `value` (любой тип, опционально) - значение ячейки. Поддерживаются типы:
  - Числа (int, float)
  - Строки (string)
  - Булевы значения (true/false)
  - Если не указано, но есть `format`, создается пустая ячейка с примененным форматированием
- `format` (объект, опционально) - форматирование ячейки

#### Format (форматирование)

Все поля опциональны:

- `bold` (boolean) - жирный шрифт
- `italic` (boolean) - курсив
- `bg_color` (string) - цвет фона ячейки в формате hex (#RRGGBB). Пример: `"#FFFF00"` (желтый)
- `text_color` (string) - цвет текста. Может быть:
  - Hex формат: `"#FF0000"` (красный)
  - Название цвета: `"red"`, `"green"`, `"blue"` и т.д.
- `font_name` (string) - название шрифта. Пример: `"Arial"`, `"Times New Roman"`
- `font_size` (number) - размер шрифта. Пример: `12`, `14.5`

### Примеры цветов

**Hex формат:**

- `"#FF0000"` - красный
- `"#00FF00"` - зеленый
- `"#0000FF"` - синий
- `"#FFFF00"` - желтый
- `"#FF00FF"` - пурпурный
- `"#00FFFF"` - голубой
- `"#FFFFFF"` - белый
- `"#000000"` - черный
- `"#F0F0F0"` - светло-серый
- `"#C0C0C0"` - серый

## Примеры использования

### Пример 1: Простой запрос без форматирования

```json
{
  "sheets": [
    {
      "title": "Products",
      "rows": [
        [{ "value": "ID" }, { "value": "Name" }, { "value": "Price" }],
        [{ "value": 1 }, { "value": "Product A" }, { "value": 100.5 }],
        [{ "value": 2 }, { "value": "Product B" }, { "value": 200.75 }]
      ]
    }
  ]
}
```

### Пример 1.1: Ячейка без значения, но с форматированием

```json
{
  "sheets": [
    {
      "title": "Sheet1",
      "rows": [
        [
          { "value": "Header" },
          { "format": { "bg_color": "#F0F0F0" } },
          { "value": "Data" }
        ]
      ]
    }
  ]
}
```

В этом примере вторая ячейка будет пустой, но с серым фоном.

### Пример 2: Запрос с форматированием

```json
{
  "sheets": [
    {
      "title": "Report",
      "rows": [
        [
          {
            "value": "Header",
            "format": {
              "bold": true,
              "bg_color": "#4472C4",
              "text_color": "#FFFFFF",
              "font_size": 14
            }
          }
        ],
        [
          {
            "value": "Data",
            "format": {
              "font_name": "Arial",
              "font_size": 12
            }
          }
        ]
      ]
    }
  ]
}
```

### Пример 3: Множественные листы

```json
{
  "sheets": [
    {
      "title": "Sheet 1",
      "rows": [[{ "value": "Data from Sheet 1" }]]
    },
    {
      "title": "Sheet 2",
      "rows": [[{ "value": "Data from Sheet 2" }]]
    }
  ]
}
```

## Обработка ошибок

Приложение возвращает следующие HTTP статусы:

- `200 OK` - успешная конвертация
- `400 Bad Request` - неверный формат данных:
  - Неверный формат JSON (для `/json2xls`)
  - Отсутствует файл в запросе (для `/xls2json`)
  - Неподдерживаемый формат файла (для `/xls2json`)
  - Ошибка при открытии Excel файла (для `/xls2json`)
- `405 Method Not Allowed` - неверный HTTP метод
- `500 Internal Server Error` - ошибка при конвертации

## Зависимости

- [excelize/v2](https://github.com/xuri/excelize) - библиотека для работы с Excel файлами

## Лицензия

MIT
