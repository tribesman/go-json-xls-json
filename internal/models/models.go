package models

// Request структура для входящего JSON
type Request struct {
	Sheets []Sheet `json:"sheets"`
}

// Sheet представляет лист в Excel
type Sheet struct {
	Title string `json:"title"`
	Rows  []Row  `json:"rows"`
}

// Row представляет строку в Excel
type Row []Cell

// Cell представляет ячейку в Excel
type Cell struct {
	Value  interface{} `json:"value,omitempty"`
	Format *Format     `json:"format,omitempty"`
}

// Format представляет форматирование ячейки
type Format struct {
	Bold      *bool   `json:"bold,omitempty"`
	Italic    *bool   `json:"italic,omitempty"`
	BgColor   *string `json:"bg_color,omitempty"`
	TextColor *string `json:"text_color,omitempty"`
	FontName  *string `json:"font_name,omitempty"`
	FontSize  *float64 `json:"font_size,omitempty"`
}

