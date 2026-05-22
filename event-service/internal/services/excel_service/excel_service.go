// internal/service/excel_service/event_excel_service.go

package excel_service

import (
	"context"
	"fmt"
	"log/slog"
	"main/internal/lib/errs"
	"main/internal/lib/liblogger"
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
	eventRepo      *event_repository.EventRepository
	parser         *excel_parser.ExcelParser
	db             orm.ORM
	subjectStorage *subject.SubjectStorage
	log            *slog.Logger
}

func NewEventExcelService(
	eventRepo *event_repository.EventRepository,
	db orm.ORM,
	log *slog.Logger,
) *EventExcelService {
	serviceLog := log.With(slog.String("owner", "EventExcelService"))
	return &EventExcelService{
		eventRepo:      eventRepo,
		db:             db,
		subjectStorage: subject.NewSubjectsStorage(),
		parser:         excel_parser.NewExcelParser(),
		log:            serviceLog,
	}
}

// CreateEventsFromExcel создает события из Excel файла
func (s *EventExcelService) CreateEventsFromExcel(
	ctx context.Context,
	file multipart.File,
	year int,
) ([]uuid.UUID, error) {
	const op = "services.excel_service.CreateEventsFromExcel"
	log := s.log.With(slog.String("op", op))

	// Парсим Excel
	rows, err := s.parser.Parse(file, year)
	if err != nil {
		log.Error("failed to parse excel", liblogger.Err(err))
		return nil, errs.ErrBadRequest.Wrap("invalid excel file")
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
	if err != nil {
		log.Error("failed to create regional stage", liblogger.Err(err))
		return nil, errs.ErrInternalError.Wrap("failed to create regional stage event")
	}
	log.Info("regional stage created", slog.String("regId", regId.String()))

	var createdIDs []uuid.UUID

	// Для каждой строки создаем иерархию событий
	for i, row := range rows {
		log.Debug("processing row", slog.Int("index", i), slog.String("subject", row.Subject))
		id, err := s.createEventHierarchy(ctx, row, regId)
		if err != nil {
			log.Error("failed to create event hierarchy",
				slog.String("subject", row.Subject),
				liblogger.Err(err),
			)
			return nil, fmt.Errorf("failed to create event hierarchy for subject %s: %w", row.Subject, err)
		}
		createdIDs = append(createdIDs, id)
	}

	log.Info("events successfully created from excel", slog.Int("count", len(createdIDs)))
	return createdIDs, nil
}

// createEventHierarchy создает иерархию событий для одного предмета
func (s *EventExcelService) createEventHierarchy(
	ctx context.Context,
	row excel_parser.ExcelRow,
	regId uuid.UUID,
) (uuid.UUID, error) {
	const op = "services.excel_service.createEventHierarchy"
	log := s.log.With(slog.String("op", op), slog.String("subject", row.Subject))

	// 1. Создаем OLYMPIAD
	olympiadID, err := s.createOlympiad(ctx, row, regId)
	if err != nil {
		log.Error("failed to create olympiad", liblogger.Err(err))
		return uuid.Nil, err
	}

	// 2. Создаем STAGE для каждой даты
	var stageIDs []uuid.UUID
	for i, date := range row.Dates {
		stageID, err := s.createStage(ctx, olympiadID, row, date, i+1)
		if err != nil {
			log.Error("failed to create stage", slog.Int("stage_number", i+1), liblogger.Err(err))
			return uuid.Nil, err
		}
		stageIDs = append(stageIDs, stageID)
	}

	// 3. Создаем CLASS для каждого класса
	for _, classNum := range row.Classes {
		_, err := s.createClass(ctx, olympiadID, row, classNum)
		if err != nil {
			log.Error("failed to create class", slog.String("class", classNum), liblogger.Err(err))
			return uuid.Nil, err
		}

		// Привязываем класс к первому стейджу (или ко всем - зависит от бизнес-логики)
		// if len(stageIDs) > 0 {
		// 	err = s.linkClassToStage(ctx, classID, stageIDs[0])
		// 	if err != nil {
		// 		log.Error("failed to link class to stage",
		// 			slog.String("classID", classID.String()),
		// 			slog.String("stageID", stageIDs[0].String()),
		// 			liblogger.Err(err),
		// 		)
		// 		return uuid.Nil, err
		// 	}
		// }
	}

	return olympiadID, nil
}

// createOlympiad создает событие типа OLYMPIAD
func (s *EventExcelService) createOlympiad(
	ctx context.Context,
	row excel_parser.ExcelRow,
	regId uuid.UUID,
) (uuid.UUID, error) {
	const op = "services.excel_service.createOlympiad"
	log := s.log.With(slog.String("op", op), slog.String("subject", row.Subject))

	// Получаем ID предмета по названию
	subjectID, err := s.subjectStorage.GetSubjectIDByName(row.Subject)
	if err != nil {
		log.Error("unknown subject", slog.String("subject_name", row.Subject), liblogger.Err(err))
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
		Subject:         subjectID,
		Profiles:        row.Profiles,
	}

	id, err := s.eventRepo.CreateEvent(ctx, s.db, eventDTO)
	if err != nil {
		log.Error("failed to create olympiad event", liblogger.Err(err))
		return uuid.Nil, errs.ErrInternalError.Wrap("failed to create olympiad event")
	}
	log.Info("olympiad created", slog.String("olympiadID", id.String()))
	return id, nil
}

// createStage создает событие типа STAGE
func (s *EventExcelService) createStage(
	ctx context.Context,
	olympiadID uuid.UUID,
	row excel_parser.ExcelRow,
	date time.Time,
	stageNumber int,
) (uuid.UUID, error) {
	const op = "services.excel_service.createStage"
	log := s.log.With(slog.String("op", op),
		slog.String("subject", row.Subject),
		slog.Int("stage_number", stageNumber),
	)

	name := fmt.Sprintf("%s - Этап %d", row.Subject, stageNumber)

	eventDTO := event.Event{
		Name:            name,
		StartDate:       date,
		EndDate:         date,
		EventType:       event.Stage,
		PreviousEventID: &olympiadID,
	}

	id, err := s.eventRepo.CreateEvent(ctx, s.db, eventDTO)
	if err != nil {
		log.Error("failed to create stage event", liblogger.Err(err))
		return uuid.Nil, errs.ErrInternalError.Wrap("failed to create stage event")
	}
	log.Debug("stage created", slog.String("stageID", id.String()))
	return id, nil
}

// createClass создает событие типа CLASS
func (s *EventExcelService) createClass(
	ctx context.Context,
	olympiadID uuid.UUID,
	row excel_parser.ExcelRow,
	classNum string,
) (uuid.UUID, error) {
	const op = "services.excel_service.createClass"
	log := s.log.With(slog.String("op", op),
		slog.String("subject", row.Subject),
		slog.String("class", classNum),
	)

	name := fmt.Sprintf("%s - %s класс", row.Subject, classNum)

	// Определяем категорию класса
	classCategory := s.determineClassCategory(row.Classes)

	eventDTO := event.Event{
		Name:            name,
		StartDate:       row.Dates[0],
		EndDate:         row.Dates[len(row.Dates)-1],
		EventType:       event.Class,
		ClassCategory:   &classCategory,
		PreviousEventID: &olympiadID,
		Profiles:        row.Profiles,
	}

	id, err := s.eventRepo.CreateEvent(ctx, s.db, eventDTO)
	if err != nil {
		log.Error("failed to create class event", liblogger.Err(err))
		return uuid.Nil, errs.ErrInternalError.Wrap("failed to create class event")
	}
	log.Debug("class created", slog.String("classID", id.String()))
	return id, nil
}

// linkClassToStage привязывает класс к стейджу (устанавливает previous_event_id)
func (s *EventExcelService) linkClassToStage(
	ctx context.Context,
	classID uuid.UUID,
	stageID uuid.UUID,
) error {
	// В реальной реализации здесь должна быть логика обновления записи в БД,
	// например: s.eventRepo.UpdateEvent(ctx, s.db, ...)
	// Пока возвращаем nil, логируя вызов.
	const op = "services.excel_service.linkClassToStage"
	log := s.log.With(slog.String("op", op),
		slog.String("classID", classID.String()),
		slog.String("stageID", stageID.String()),
	)
	log.Info("linking class to stage")
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
