// internal/service/excel_parser.go

package excel_parser

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExcelRow представляет строку из Excel файла
type ExcelRow struct {
	Subject  string
	Classes  []string
	Dates    []time.Time
	Profiles []string
}

// ExcelParser парсит Excel файл
type ExcelParser struct{}

func NewExcelParser() *ExcelParser {
	return &ExcelParser{}
}

// Parse парсит Excel файл и возвращает слайс строк
func (p *ExcelParser) Parse(file multipart.File, year int) ([]ExcelRow, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open excel file: %w", err)
	}
	defer f.Close()

	// Получаем первую таблицу
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("excel file must have header and at least one data row")
	}

	// Пропускаем заголовок
	var result []ExcelRow
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 4 {
			continue // Пропускаем пустые строки
		}

		excelRow, err := p.parseRow(row, year)
		if err != nil {
			return nil, fmt.Errorf("error parsing row %d: %w", i+1, err)
		}

		result = append(result, excelRow)
	}

	return result, nil
}

// parseRow парсит одну строку Excel
func (p *ExcelParser) parseRow(row []string, year int) (ExcelRow, error) {
	var excelRow ExcelRow

	// Парсим subject
	excelRow.Subject = strings.TrimSpace(row[0])
	if excelRow.Subject == "" {
		return excelRow, fmt.Errorf("subject is empty")
	}

	// Парсим classes
	classesStr := strings.TrimSpace(row[1])
	if classesStr == "" {
		return excelRow, fmt.Errorf("classes is empty")
	}
	classes := strings.Split(classesStr, ",")
	for _, class := range classes {
		class = strings.TrimSpace(class)
		if class != "" {
			excelRow.Classes = append(excelRow.Classes, class)
		}
	}

	// Парсим dates
	datesStr := strings.TrimSpace(row[2])
	if datesStr == "" {
		return excelRow, fmt.Errorf("dates is empty")
	}
	dates := strings.Split(datesStr, ",")
	for _, dateStr := range dates {
		dateStr = strings.TrimSpace(dateStr)
		if dateStr == "" {
			continue
		}

		// Пробуем разные форматы дат
		var date time.Time
		var err error

		formats := []string{
			"02.01.2006",
			"2.1.2006",
			"2006-01-02",
			"01/02/2006",
		}

		for _, format := range formats {
			date, err = time.Parse(format, dateStr)
			if err == nil {
				break
			}
		}

		if err != nil {
			return excelRow, fmt.Errorf("invalid date format: %s", dateStr)
		}

		// Устанавливаем год
		date = time.Date(year, date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		excelRow.Dates = append(excelRow.Dates, date)
	}

	// Парсим profiles (опционально)
	profilesStr := strings.TrimSpace(row[3])
	if profilesStr != "" && profilesStr != "-" {
		profiles := strings.Split(profilesStr, ",")
		for _, profile := range profiles {
			profile = strings.TrimSpace(profile)
			if profile != "" {
				excelRow.Profiles = append(excelRow.Profiles, profile)
			}
		}
	}

	return excelRow, nil
}
