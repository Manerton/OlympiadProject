package subject

import (
	"strings"
	"sync"
)

type Subject struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(100);not null;unique"`
}

type SubjectStorage struct {
	subjects    map[int]string
	nameToID    map[string]int
	mu          sync.RWMutex
	initialized bool
}

const (
	SubjectEconomics = iota + 1
	SubjectBiology
	SubjectSpanish
	SubjectInformatics
	SubjectHistory
	SubjectChemistry
	SubjectRussian
	SubjectItalian
	SubjectChinese
	SubjectSocialScience
	SubjectPhysics
	SubjectMath
	SubjectGerman
	SubjectEcology
	SubjectLiterature
	SubjectLaw
	SubjectLabor
	SubjectGeography
	SubjectFrench
	SubjectEnglish
	SubjectPhysicalEducation
	SubjectAstronomy
	SubjectLifeSafety
	SubjectArt
)

func NewSubjectsStorage() *SubjectStorage {
	SubjectIDToName := map[int]string{
		SubjectEconomics:         "Экономика",
		SubjectBiology:           "Биология",
		SubjectSpanish:           "Испанский язык",
		SubjectInformatics:       "Информатика",
		SubjectHistory:           "История",
		SubjectChemistry:         "Химия",
		SubjectRussian:           "Русский язык",
		SubjectItalian:           "Итальянский язык",
		SubjectChinese:           "Китайский язык",
		SubjectSocialScience:     "Обществознание",
		SubjectPhysics:           "Физика",
		SubjectMath:              "Математика",
		SubjectGerman:            "Немецкий язык",
		SubjectEcology:           "Экология",
		SubjectLiterature:        "Литература",
		SubjectLaw:               "Право",
		SubjectLabor:             "Труд",
		SubjectGeography:         "География",
		SubjectFrench:            "Французский язык",
		SubjectEnglish:           "Английский язык",
		SubjectPhysicalEducation: "Физическая культура",
		SubjectAstronomy:         "Астрономия",
		SubjectLifeSafety:        "Основы безопасности и защиты Родины",
		SubjectArt:               "Искусство",
	}

	storage := &SubjectStorage{
		subjects: SubjectIDToName,
		nameToID: make(map[string]int, len(SubjectIDToName)),
	}

	// Строим обратный маппинг
	storage.buildNameToIDMap()
	storage.initialized = true

	return storage
}

// buildNameToIDMap строит обратный маппинг name -> id
func (s *SubjectStorage) buildNameToIDMap() {
	for id, name := range s.subjects {
		// Добавляем точное совпадение
		s.nameToID[strings.ToLower(name)] = id

		// Добавляем сокращенные варианты
		shortName := s.getShortName(name)
		if shortName != "" {
			s.nameToID[strings.ToLower(shortName)] = id
		}
	}
}

// getShortName возвращает короткое название предмета
func (s *SubjectStorage) getShortName(fullName string) string {
	// Убираем " язык" из названий языков
	shortName := strings.TrimSuffix(fullName, " язык")
	if shortName != fullName {
		return shortName
	}

	// Другие сокращения
	switch fullName {
	case "Физическая культура":
		return "физкультура"
	case "Основы безопасности и защиты Родины":
		return "обж"
	case "Искусство":
		return "мхк"
	}

	return ""
}

// GetSubjectIDByName возвращает ID предмета по его названию
func (s *SubjectStorage) GetSubjectIDByName(name string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Очищаем и нормализуем входное имя
	cleanName := strings.TrimSpace(strings.ToLower(name))

	// 1. Пробуем точное совпадение
	if id, exists := s.nameToID[cleanName]; exists {
		return id, nil
	}

	// 2. Пробуем найти по частичному совпадению
	for storedName, id := range s.nameToID {
		if strings.Contains(storedName, cleanName) ||
			strings.Contains(cleanName, storedName) {
			return id, nil
		}
	}

	// 3. Специальные случаи и синонимы
	synonyms := s.getSynonyms()
	if id, exists := synonyms[cleanName]; exists {
		return id, nil
	}

	return 0, &SubjectNotFoundError{Name: name}
}

// getSynonyms возвращает карту синонимов для предметов
func (s *SubjectStorage) getSynonyms() map[string]int {
	return map[string]int{
		"английский":     SubjectEnglish,
		"english":        SubjectEnglish,
		"немецкий":       SubjectGerman,
		"deutsch":        SubjectGerman,
		"французский":    SubjectFrench,
		"français":       SubjectFrench,
		"испанский":      SubjectSpanish,
		"español":        SubjectSpanish,
		"итальянский":    SubjectItalian,
		"italiano":       SubjectItalian,
		"китайский":      SubjectChinese,
		"русский":        SubjectRussian,
		"математика":     SubjectMath,
		"math":           SubjectMath,
		"физика":         SubjectPhysics,
		"химия":          SubjectChemistry,
		"биология":       SubjectBiology,
		"информатика":    SubjectInformatics,
		"informatics":    SubjectInformatics,
		"икт":            SubjectInformatics,
		"обществознание": SubjectSocialScience,
		"общество":       SubjectSocialScience,
		"литература":     SubjectLiterature,
		"история":        SubjectHistory,
		"география":      SubjectGeography,
		"экономика":      SubjectEconomics,
		"право":          SubjectLaw,
		"экология":       SubjectEcology,
		"астрономия":     SubjectAstronomy,
		"физкультура":    SubjectPhysicalEducation,
		"физ-ра":         SubjectPhysicalEducation,
		"обж":            SubjectLifeSafety,
		"ож":             SubjectLifeSafety,
		"мхк":            SubjectArt,
		"искусство":      SubjectArt,
		"труд":           SubjectLabor,
		"технология":     SubjectLabor,
	}
}

// GetSubjectNameByID возвращает название предмета по его ID
func (s *SubjectStorage) GetSubjectNameByID(id int) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	name, exists := s.subjects[id]
	if !exists {
		return "", &SubjectNotFoundError{ID: id}
	}

	return name, nil
}

// GetAllSubjects возвращает все предметы
func (s *SubjectStorage) GetAllSubjects() map[int]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Возвращаем копию, чтобы избежать изменений извне
	subjects := make(map[int]string, len(s.subjects))
	for id, name := range s.subjects {
		subjects[id] = name
	}

	return subjects
}

// SubjectExists проверяет существование предмета по ID
func (s *SubjectStorage) SubjectExists(id int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.subjects[id]
	return exists
}

// SubjectNameExists проверяет существование предмета по имени
func (s *SubjectStorage) SubjectNameExists(name string) bool {
	_, err := s.GetSubjectIDByName(name)
	return err == nil
}

// SearchSubjects ищет предметы по подстроке
func (s *SubjectStorage) SearchSubjects(query string) map[int]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query = strings.ToLower(strings.TrimSpace(query))
	results := make(map[int]string)

	for id, name := range s.subjects {
		if strings.Contains(strings.ToLower(name), query) {
			results[id] = name
		}
	}

	return results
}

// GetSubjectIDByNameFuzzy ищет ID предмета с неточным совпадением
func (s *SubjectStorage) GetSubjectIDByNameFuzzy(name string) (int, float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Сначала пробуем точное совпадение
	if id, err := s.GetSubjectIDByName(name); err == nil {
		return id, 1.0, nil
	}

	cleanName := strings.ToLower(strings.TrimSpace(name))

	// Ищем лучшее частичное совпадение
	var bestMatch int
	bestScore := 0.0

	for id, subjectName := range s.subjects {
		score := calculateStringSimilarity(cleanName, strings.ToLower(subjectName))
		if score > bestScore {
			bestScore = score
			bestMatch = id
		}
	}

	if bestScore > 0.5 { // Порог схожести 50%
		return bestMatch, bestScore, nil
	}

	return 0, 0, &SubjectNotFoundError{Name: name}
}

// calculateStringSimilarity вычисляет примерную схожесть строк
func calculateStringSimilarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	// Простой алгоритм на основе общих подстрок
	r1 := []rune(s1)
	r2 := []rune(s2)

	if len(r1) == 0 || len(r2) == 0 {
		return 0.0
	}

	// Считаем количество совпадающих символов
	matches := 0
	for _, c1 := range r1 {
		for _, c2 := range r2 {
			if c1 == c2 {
				matches++
				break
			}
		}
	}

	// Вычисляем коэффициент схожести
	similarity := float64(matches) / float64(max(len(r1), len(r2)))

	// Увеличиваем вес, если одна строка содержится в другой
	if strings.Contains(s1, s2) || strings.Contains(s2, s1) {
		similarity = (similarity + 1.0) / 2.0
	}

	return similarity
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// SubjectNotFoundError ошибка, когда предмет не найден
type SubjectNotFoundError struct {
	Name string
	ID   int
}

func (e *SubjectNotFoundError) Error() string {
	if e.Name != "" {
		return "предмет не найден: " + e.Name
	}
	return "предмет с ID не найден: " + string(rune(e.ID))
}
