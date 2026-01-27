#!/bin/bash

# Скрипт для сборки проекта json2xls

set -e

# Включаем поддержку Go модулей
export GO111MODULE=on

echo "🔧 Настройка окружения..."
go env -w GO111MODULE=on

echo "📦 Загрузка зависимостей..."
go mod download
go mod tidy

echo "🔨 Сборка приложения..."
go build -o ./bin/json2xls ./cmd/json2xls

echo "✅ Сборка завершена успешно!"
echo "Запустите приложение: ./bin/json2xls"

