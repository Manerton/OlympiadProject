package excel_parser

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelRow struct {
	Subject  string
	Classes  []string
	Dates    []time.Time
	Profiles []string
}

type ExcelParser struct{}

func NewExcelParser() *ExcelParser {
	return &ExcelParser{}
}

func (p *ExcelParser) Parse(file multipart.File, year int) ([]ExcelRow, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open excel file: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName) // нужно только чтобы узнать количество строк
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("excel file must have header and at least one data row")
	}

	var result []ExcelRow
	// читаем строки, начиная со 2-й (индекс i в GetRows – это номер строки Excel, 1-я – заголовок)
	for i := 2; i <= len(rows); i++ {
		excelRow, err := p.parseRow(f, sheetName, i, year)
		if err != nil {
			return nil, fmt.Errorf("error parsing row %d: %w", i, err)
		}
		result = append(result, excelRow)
	}
	return result, nil
}

func (p *ExcelParser) parseRow(f *excelize.File, sheetName string, rowNum int, year int) (ExcelRow, error) {
	var excelRow ExcelRow

	// --- Subject ---
	cell := fmt.Sprintf("A%d", rowNum)
	subjectVal, err := f.GetCellValue(sheetName, cell)
	if err != nil {
		return excelRow, fmt.Errorf("failed to read subject cell: %w", err)
	}
	excelRow.Subject = strings.TrimSpace(subjectVal)
	if excelRow.Subject == "" {
		return excelRow, fmt.Errorf("subject is empty")
	}

	// --- Classes ---
	cell = fmt.Sprintf("B%d", rowNum)
	classesVal, err := f.GetCellValue(sheetName, cell)
	if err != nil {
		return excelRow, fmt.Errorf("failed to read classes cell: %w", err)
	}
	classesStr := strings.TrimSpace(classesVal)
	if classesStr == "" {
		return excelRow, fmt.Errorf("classes is empty")
	}
	for _, class := range strings.Split(classesStr, ",") {
		class = strings.TrimSpace(class)
		if class != "" {
			excelRow.Classes = append(excelRow.Classes, class)
		}
	}

	// --- Dates ---
	cell = fmt.Sprintf("C%d", rowNum)
	datesVal, err := f.GetCellValue(sheetName, cell)
	if err != nil {
		return excelRow, fmt.Errorf("failed to read dates cell: %w", err)
	}
	datesStr := strings.TrimSpace(datesVal)
	if datesStr == "" {
		return excelRow, fmt.Errorf("dates is empty")
	}

	for _, part := range strings.Split(datesStr, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		var date time.Time
		// Пробуем как число (Excel serial date)
		if serial, convErr := strconv.ParseFloat(part, 64); convErr == nil {
			date, err = excelize.ExcelDateToTime(serial, false)
			if err != nil {
				return excelRow, fmt.Errorf("invalid excel serial date: %s", part)
			}
		} else {
			// Не число – пробуем строковые форматы (порядок важен)
			formats := []string{
				"02.01.2006", // DD.MM.YYYY (ваш исходный формат)
				"2.1.2006",   // D.M.YYYY
				"1-2-06",     // M-D-YY (на случай, если GetCellValue вернёт американский формат)
				"2006-01-02", // YYYY-MM-DD
				"01/02/2006", // MM/DD/YYYY
			}
			parsed := false
			for _, format := range formats {
				date, err = time.Parse(format, part)
				if err == nil {
					parsed = true
					break
				}
			}
			if !parsed {
				return excelRow, fmt.Errorf("invalid date format: %s", part)
			}
		}
		// Устанавливаем нужный год
		date = time.Date(year, date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		excelRow.Dates = append(excelRow.Dates, date)
	}

	// --- Profiles ---
	cell = fmt.Sprintf("D%d", rowNum)
	profilesVal, err := f.GetCellValue(sheetName, cell)
	if err != nil {
		return excelRow, fmt.Errorf("failed to read profiles cell: %w", err)
	}
	profilesStr := strings.TrimSpace(profilesVal)
	if profilesStr != "" && profilesStr != "-" {
		for _, profile := range strings.Split(profilesStr, ",") {
			profile = strings.TrimSpace(profile)
			if profile != "" {
				excelRow.Profiles = append(excelRow.Profiles, profile)
			}
		}
	}

	return excelRow, nil
}
