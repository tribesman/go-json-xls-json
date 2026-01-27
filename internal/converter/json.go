package converter

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"json2xls/internal/models"

	"github.com/xuri/excelize/v2"
)

// formatColorToHex конвертирует цвет из формата excelize в hex формат с #
func formatColorToHex(color string) string {
	if color == "" {
		return ""
	}
	// Если уже есть #, возвращаем как есть
	if strings.HasPrefix(color, "#") {
		return strings.ToUpper(color)
	}
	// Добавляем # если его нет
	return "#" + strings.ToUpper(color)
}

// extractFormatFromStyle извлекает форматирование из стиля Excel
func extractFormatFromStyle(f *excelize.File, styleID int) (*models.Format, error) {
	style, err := f.GetStyle(styleID)
	if err != nil {
		return nil, err
	}

	format := &models.Format{}

	// Извлекаем информацию о шрифте
	if style.Font != nil {
		if style.Font.Bold {
			format.Bold = &[]bool{true}[0]
		}
		if style.Font.Italic {
			format.Italic = &[]bool{true}[0]
		}
		if style.Font.Family != "" {
			format.FontName = &style.Font.Family
		}
		if style.Font.Size > 0 {
			format.FontSize = &style.Font.Size
		}
		if style.Font.Color != "" {
			color := formatColorToHex(style.Font.Color)
			format.TextColor = &color
		}
	}

	// Извлекаем информацию о фоне
	if style.Fill.Type == "pattern" && len(style.Fill.Color) > 0 && style.Fill.Color[0] != "" {
		color := formatColorToHex(style.Fill.Color[0])
		format.BgColor = &color
	}

	// Возвращаем nil если нет форматирования
	if format.Bold == nil && format.Italic == nil && format.BgColor == nil &&
		format.TextColor == nil && format.FontName == nil && format.FontSize == nil {
		return nil, nil
	}

	return format, nil
}

// ConvertXLSToJSON конвертирует Excel файл в JSON
func ConvertXLSToJSON(f *excelize.File) (*models.Request, error) {
	req := &models.Request{
		Sheets: []models.Sheet{},
	}

	// Получаем список всех листов
	sheetList := f.GetSheetList()

	for _, sheetName := range sheetList {
		sheet := models.Sheet{
			Title: sheetName,
			Rows:  []models.Row{},
		}

		// Получаем все строки листа
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return nil, fmt.Errorf("error reading rows from sheet %s: %w", sheetName, err)
		}

		// Обрабатываем каждую строку
		for rowIdx, row := range rows {
			excelRow := models.Row{}
			excelRowNum := rowIdx + 1 // Excel использует 1-based индексацию

			// Обрабатываем каждую ячейку в строке
			for colIdx := range row {
				excelColNum := colIdx + 1 // Excel использует 1-based индексацию

				cellName, err := excelize.CoordinatesToCellName(excelColNum, excelRowNum)
				if err != nil {
					return nil, fmt.Errorf("error converting coordinates: %w", err)
				}

				// Получаем значение ячейки
				cellValue, err := f.GetCellValue(sheetName, cellName)
				if err != nil {
					return nil, fmt.Errorf("error getting cell value: %w", err)
				}

				// Пытаемся определить тип значения
				var typedValue interface{} = cellValue

				// Пробуем преобразовать в число
				if cellValue != "" {
					if intVal, err := strconv.ParseInt(cellValue, 10, 64); err == nil {
						typedValue = intVal
					} else if floatVal, err := strconv.ParseFloat(cellValue, 64); err == nil {
						typedValue = floatVal
					} else if boolVal, err := strconv.ParseBool(cellValue); err == nil {
						typedValue = boolVal
					}
				}

				// Получаем стиль ячейки
				styleID, err := f.GetCellStyle(sheetName, cellName)
				if err != nil {
					// Если не удалось получить стиль, продолжаем без форматирования
					excelRow = append(excelRow, models.Cell{
						Value:  typedValue,
						Format: nil,
					})
					continue
				}

				// Извлекаем форматирование
				format, err := extractFormatFromStyle(f, styleID)
				if err != nil {
					log.Printf("Warning: error extracting format for cell %s: %v", cellName, err)
				}

				excelRow = append(excelRow, models.Cell{
					Value:  typedValue,
					Format: format,
				})
			}

			// Добавляем строку только если в ней есть ячейки
			if len(excelRow) > 0 {
				sheet.Rows = append(sheet.Rows, excelRow)
			}
		}

		req.Sheets = append(req.Sheets, sheet)
	}

	return req, nil
}

