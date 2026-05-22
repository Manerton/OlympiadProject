// internal/service/event_excel_service.go

package excel_service

import (
	"context"
	"fmt"
	"main/internal/models/event"
	"main/internal/models/subject"
	"main/internal/repositories/event_repository"
	"main/internal/storage/orm"
	"main/support/excel_parser"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type EventExcelService struct {
	eventRepo      event_repository.EventRepository
	parser         *excel_parser.ExcelParser
	db             orm.ORM
	subjectStorage *subject.SubjectStorage
}

func NewEventExcelService(eventRepo event_repository.EventRepository, db orm.ORM) *EventExcelService {
	return &EventExcelService{
		eventRepo:      eventRepo,
		db:             db,
		subjectStorage: subject.NewSubjectsStorage(),
		parser:         excel_parser.NewExcelParser(),
	}
}

// CreateEventsFromExcel создает события из Excel файла
func (s *EventExcelService) CreateEventsFromExcel(
	ctx context.Context,
	file multipart.File,
	year int,
) ([]uuid.UUID, error) {
	// Парсим Excel
	rows, err := s.parser.Parse(file, year)
	if err != nil {
		return nil, fmt.Errorf("failed to parse excel: %w", err)
	}

	startDate := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)

	endDate := time.Date(year, time.December, 31, 23, 59, 59, 0, time.UTC)

	regionalStage := event.Event{
		Name:      fmt.Sprintf("ВСОШ %d", year),
		StartDate: startDate,
		EndDate:   endDate,
		EventType: event.RegionalStage,
	}

	regId, err := s.eventRepo.CreateEvent(ctx, s.db, regionalStage)

	var createdIDs []uuid.UUID

	// Для каждой строки создаем иерархию событий
	for _, row := range rows {
		id, err := s.createEventHierarchy(ctx, row, regId)
		if err != nil {
			return nil, fmt.Errorf("failed to create event hierarchy for subject %s: %w", row.Subject, err)
		}
		createdIDs = append(createdIDs, id)
	}

	return createdIDs, nil
}

// createEventHierarchy создает иерархию событий для одного предмета
func (s *EventExcelService) createEventHierarchy(
	ctx context.Context,
	row excel_parser.ExcelRow,
	regId uuid.UUID,
) (uuid.UUID, error) {
	// 1. Создаем OLYMPIAD
	olympiadID, err := s.createOlympiad(ctx, row, regId)
	if err != nil {
		return uuid.Nil, err
	}

	// 2. Создаем STAGE для каждой даты
	var stageIDs []uuid.UUID
	for i, date := range row.Dates {
		stageID, err := s.createStage(ctx, olympiadID, row, date, i+1)
		if err != nil {
			return uuid.Nil, err
		}
		stageIDs = append(stageIDs, stageID)
	}

	// 3. Создаем CLASS для каждого класса
	for _, classNum := range row.Classes {
		classID, err := s.createClass(ctx, olympiadID, row, classNum)
		if err != nil {
			return uuid.Nil, err
		}

		// Привязываем класс к первому стейджу (или ко всем - зависит от бизнес-логики)
		if len(stageIDs) > 0 {
			err = s.linkClassToStage(ctx, classID, stageIDs[0])
			if err != nil {
				return uuid.Nil, err
			}
		}
	}

	return olympiadID, nil
}

// createOlympiad создает событие типа OLYMPIAD
// internal/service/event_excel_service.go

func (s *EventExcelService) createOlympiad(
	ctx context.Context,
	row excel_parser.ExcelRow,
	regId uuid.UUID,
) (uuid.UUID, error) {
	// Получаем ID предмета по названию
	subjectID, err := s.subjectStorage.GetSubjectIDByName(row.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("unknown subject '%s': %w", row.Subject, err)
	}

	name := fmt.Sprintf("Олимпиада по %s", row.Subject)

	// Определяем start и end date (берем первую и последнюю дату)
	startDate := row.Dates[0]
	endDate := row.Dates[len(row.Dates)-1]

	eventDTO := event.Event{
		Name:            name,
		PreviousEventID: &regId,
		StartDate:       startDate,
		EndDate:         endDate,
		EventType:       event.Olympiad,
		Subject:         subjectID, // Используем числовой ID
		Profiles:        row.Profiles,
	}

	return s.eventRepo.CreateEvent(ctx, s.db, eventDTO)
}

// createStage создает событие типа STAGE
func (s *EventExcelService) createStage(
	ctx context.Context,
	olympiadID uuid.UUID,
	row excel_parser.ExcelRow,
	date time.Time,
	stageNumber int,
) (uuid.UUID, error) {
	name := fmt.Sprintf("%s - Этап %d", row.Subject, stageNumber)

	eventDTO := event.Event{
		Name:            name,
		StartDate:       date,
		EndDate:         date,
		EventType:       event.Stage,
		PreviousEventID: &olympiadID,
	}

	return s.eventRepo.CreateEvent(ctx, s.db, eventDTO)
}

// createClass создает событие типа CLASS
func (s *EventExcelService) createClass(
	ctx context.Context,
	olympiadID uuid.UUID,
	row excel_parser.ExcelRow,
	classNum string,
) (uuid.UUID, error) {
	name := fmt.Sprintf("%s - %s класс", row.Subject, classNum)

	// Определяем категорию класса
	classCategory := s.determineClassCategory(row.Classes)

	eventDTO := event.Event{
		Name:            name,
		StartDate:       row.Dates[0],
		EndDate:         row.Dates[len(row.Dates)-1],
		EventType:       event.Class,
		ClassCategory:   classCategory,
		PreviousEventID: &olympiadID,
		Profiles:        row.Profiles,
	}

	return s.eventRepo.CreateEvent(ctx, s.db, eventDTO)
}

// linkClassToStage привязывает класс к стейджу (устанавливает previous_event_id)
func (s *EventExcelService) linkClassToStage(
	ctx context.Context,
	classID uuid.UUID,
	stageID uuid.UUID,
) error {
	// Здесь может быть дополнительная логика привязки
	// Например, обновление previous_event_id у класса
	return nil
}

// determineClassCategory определяет категорию класса на основе списка классов
func (s *EventExcelService) determineClassCategory(classes []string) event.ClassCategoryType {
	if len(classes) == 1 {
		switch classes[0] {
		case "9":
			return event.Class9
		case "10":
			return event.Class10
		case "11":
			return event.Class11
		}
	}

	// Проверяем комбинации
	has9 := contains(classes, "9")
	has10 := contains(classes, "10")
	has11 := contains(classes, "11")

	if has9 && has10 && has11 {
		return event.Class9_11
	}
	if has9 && has10 {
		return event.Class9_10
	}
	if has10 && has11 {
		return event.Class10_11
	}

	return event.Class9_11 // По умолчанию
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
