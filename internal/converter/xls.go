package converter

import (
	"fmt"
	"strings"

	"json2xls/internal/models"

	"github.com/xuri/excelize/v2"
)

// normalizeColor нормализует цвет из hex формата (#RRGGBB) в формат для excelize
func normalizeColor(color string) string {
	if color == "" {
		return ""
	}
	// Убираем # если есть
	color = strings.TrimPrefix(color, "#")
	// Excelize ожидает формат без #
	return strings.ToUpper(color)
}

// applyFormat применяет форматирование к ячейке
func applyFormat(f *excelize.File, sheetName string, cell string, format *models.Format) error {
	if format == nil {
		return nil
	}

	style := &excelize.Style{}

	// Настройка шрифта
	font := &excelize.Font{}
	hasFontFormat := false
	if format.Bold != nil {
		font.Bold = *format.Bold
		hasFontFormat = true
	}
	if format.Italic != nil {
		font.Italic = *format.Italic
		hasFontFormat = true
	}
	if format.FontName != nil && *format.FontName != "" {
		font.Family = *format.FontName
		hasFontFormat = true
	}
	if format.FontSize != nil && *format.FontSize > 0 {
		font.Size = *format.FontSize
		hasFontFormat = true
	}
	if format.TextColor != nil && *format.TextColor != "" {
		font.Color = normalizeColor(*format.TextColor)
		hasFontFormat = true
	}
	if hasFontFormat {
		style.Font = font
	}

	// Настройка фона
	if format.BgColor != nil && *format.BgColor != "" {
		style.Fill = excelize.Fill{
			Type:    "pattern",
			Color:   []string{normalizeColor(*format.BgColor)},
			Pattern: 1,
		}
	}

	// Создаем стиль только если есть какое-либо форматирование
	if style.Font == nil && style.Fill.Type == "" {
		return nil
	}

	styleID, err := f.NewStyle(style)
	if err != nil {
		return err
	}

	return f.SetCellStyle(sheetName, cell, cell, styleID)
}

// ConvertToXLS конвертирует JSON запрос в XLS файл
func ConvertToXLS(req *models.Request) (*excelize.File, error) {
	f := excelize.NewFile()

	// Если нет листов, используем дефолтный
	if len(req.Sheets) == 0 {
		return f, nil
	}

	// Удаляем дефолтный лист если есть листы в запросе
	sheetList := f.GetSheetList()
	if len(sheetList) > 0 {
		// Удаляем Sheet1 только если он существует и мы создаем новые листы
		for _, sheetName := range sheetList {
			if sheetName == "Sheet1" {
				if err := f.DeleteSheet("Sheet1"); err != nil {
					return nil, fmt.Errorf("error deleting default sheet: %w", err)
				}
				break
			}
		}
	}

	// Создаем листы
	for sheetIndex, sheet := range req.Sheets {
		sheetName := sheet.Title
		if sheetName == "" {
			sheetName = fmt.Sprintf("Sheet%d", sheetIndex+1)
		}

		// Создаем новый лист
		index, err := f.NewSheet(sheetName)
		if err != nil {
			return nil, fmt.Errorf("error creating sheet %s: %w", sheetName, err)
		}

		// Заполняем данные
		for rowIndex, row := range sheet.Rows {
			for colIndex, cell := range row {
				// Пропускаем ячейку, если нет value
				if cell.Value == nil {
					continue
				}

				cellName, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
				if err != nil {
					return nil, fmt.Errorf("error converting coordinates: %w", err)
				}

				// Устанавливаем значение ячейки
				if err := f.SetCellValue(sheetName, cellName, cell.Value); err != nil {
					return nil, fmt.Errorf("error setting cell value: %w", err)
				}

				// Применяем форматирование (если указано)
				if cell.Format != nil {
					if err := applyFormat(f, sheetName, cellName, cell.Format); err != nil {
						return nil, fmt.Errorf("error applying format: %w", err)
					}
				}
			}
		}

		// Устанавливаем активный лист
		if sheetIndex == 0 {
			f.SetActiveSheet(index)
		}
	}

	return f, nil
}

