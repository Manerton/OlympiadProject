package event_repository_test

import (
	"fmt"
	"main/internal/models/event"
	"main/internal/models/subject"
	"main/internal/repositories/event_repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Создаем in-memory базу данных SQLite для тестирования
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	db.Exec("PRAGMA foreign_keys = ON")
	// create enum (types Events)
	create_enum_type_str := fmt.Sprintf(
		`DO $$ BEGIN IF NOT EXISTS
		(SELECT 1 FROM pg_type WHERE typname = 'event_type') 
		THEN CREATE TYPE event_type AS ENUM 
		('%s', '%s', '%s', '%s');
		END IF; END $$;`,
		event.RegionalStage,
		event.Olympiad,
		event.Stage,
		event.Appeal,
	)
	db.Exec(create_enum_type_str)

	// Миграция таблиц для тестирования
	if err := db.AutoMigrate(&subject.Subject{}, &event.Event{}); err != nil {
		t.Fatalf("failed to migrate tables: %v", err)
	}

	// Создаём тестовые данные для вспомогатльной таблицы subject
	testSubjects := []subject.Subject{
		{
			ID:   1,
			Name: "Test 1",
		},
		{
			ID:   2,
			Name: "Test 2",
		},
		{
			ID:   3,
			Name: "Test 3",
		},
	}
	if err = db.Create(&testSubjects).Error; err != nil {
		t.Fatalf("failed to create support subjects: %v", err)
	}

	return db
}

func TestCreateEvent(t *testing.T) {
	db := setupTestDB(t)
	repo := event_repository.EventRepository{}

	var SubjectID uint = 1
	var UnCorrectPreviousEventID uint = 99
	var UnCorrectSubjectID uint = 99

	main_event := event.Event{
		Model:     gorm.Model{ID: 1},
		Name:      "MainEvent",
		StartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
		EndDate:   time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
		EventType: event.RegionalStage,
	}
	if err := db.Create(&main_event).Error; err != nil {
		return
	}

	testCases := []struct {
		name           string
		event          event.Event
		expectError    bool
		expectedErrMsg string
		expectedId     uint
	}{
		{
			name: "Created correct event RegionalStage",
			event: event.Event{
				Model:     gorm.Model{ID: 2},
				Name:      "RegionalStage",
				StartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				EndDate:   time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
				EventType: event.RegionalStage,
			},
			expectedId:  2,
			expectError: false,
		},
		{
			name: "Create event with unnecessary fields",
			event: event.Event{
				Model:           gorm.Model{ID: 3},
				Name:            "RegionalStage2",
				StartDate:       time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				EndDate:         time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
				EventType:       event.RegionalStage,
				SubjectID:       &SubjectID,
				PreviousEventID: &main_event.ID,
			},
			expectedId:  3,
			expectError: false,
		},
		{
			name: "Create event type Olympiad with correct foreignKey",
			event: event.Event{
				Model:           gorm.Model{ID: 4},
				Name:            "Olympiad",
				StartDate:       time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				EndDate:         time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
				EventType:       event.Olympiad,
				SubjectID:       &SubjectID,
				PreviousEventID: &main_event.ID,
			},
			expectedId:  4,
			expectError: false,
		},
		{
			name: "Create event type Olympiad with uncorrect Date",
			event: event.Event{
				Model:           gorm.Model{ID: 5},
				Name:            "Olympiad",
				StartDate:       time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
				EndDate:         time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				EventType:       event.Olympiad,
				SubjectID:       &SubjectID,
				PreviousEventID: &main_event.ID,
			},
			expectedId:  5,
			expectError: false,
		},
		{
			name: "Create event type stage without fields SubjectID/PreviousID",
			event: event.Event{
				Model:     gorm.Model{ID: 6},
				Name:      "Stage",
				StartDate: time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
				EndDate:   time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				EventType: event.Stage,
			},
			expectedId:  6,
			expectError: false,
		},
		{
			name: "Create event type appeal",
			event: event.Event{
				Model:     gorm.Model{ID: 7},
				Name:      "Appeal",
				StartDate: time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
				EndDate:   time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				EventType: event.Appeal,
			},
			expectedId:  7,
			expectError: false,
		},
		{
			name: "Create event without fields",
			event: event.Event{
				Model: gorm.Model{ID: 8},
			},
			expectedId:  8,
			expectError: false,
		},
		{
			name: "Create event with uncorrect PreviousID",
			event: event.Event{
				Model:           gorm.Model{ID: 9},
				Name:            "Appeal",
				StartDate:       time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
				EndDate:         time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				EventType:       event.Appeal,
				PreviousEventID: &UnCorrectPreviousEventID,
				SubjectID:       &SubjectID,
			},
			expectError:    true,
			expectedErrMsg: "FOREIGN KEY constraint failed",
		},
		{
			name: "Create event with uncorrect SubjectID",
			event: event.Event{
				Model:           gorm.Model{ID: 10},
				Name:            "Appeal",
				StartDate:       time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
				EndDate:         time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				EventType:       event.Appeal,
				PreviousEventID: &main_event.ID,
				SubjectID:       &UnCorrectSubjectID,
			},
			expectError:    true,
			expectedErrMsg: "FOREIGN KEY constraint failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := repo.CreateEvent(db, tc.event)
			if tc.expectError {
				assert.Error(t, err, "expected an error for case: %s", tc.name)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err, "expected no error for case: %s", tc.name)
				assert.Equal(t, tc.expectedId, id, "expected event id to match for case: %s", tc.name)

			}
		})
	}
}

func TestGetEventByID(t *testing.T) {
	db := setupTestDB(t)
	repo := event_repository.EventRepository{}

	tempEvent := event.Event{
		Name:      "RegionalStage",
		StartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
		EndDate:   time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
		EventType: event.RegionalStage,
	}
	if err := db.Create(&tempEvent).Error; err != nil {
		return
	}

	testCases := []struct {
		name              string
		searched_id       uint
		expectError       bool
		expectedStartDate time.Time
		expectedEndDate   time.Time
		expectedEventType event.EventType
		expectedErrMsg    string
	}{
		{
			name:              "Get existing event",
			searched_id:       1,
			expectedStartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
			expectedEndDate:   time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
			expectedEventType: event.RegionalStage,
			expectError:       false,
		},
		{
			name:              "Get first event on default id = 0",
			searched_id:       0,
			expectedStartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
			expectedEndDate:   time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
			expectedEventType: event.RegionalStage,
			expectError:       false,
		},
		{
			name:           "Non-existing event",
			searched_id:    9999,
			expectError:    true,
			expectedErrMsg: "record not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event, err := repo.GetEventByID(db, tc.searched_id)
			if tc.expectError {
				assert.Error(t, err, "expected an error for case: %s", tc.name)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err, "expected no error for case: %s", tc.name)
				assert.Equal(t, tc.expectedStartDate, event.StartDate, "expected start date event to match for case: %s", tc.name)
				assert.Equal(t, tc.expectedEndDate, event.EndDate, "expected end date event to match for case: %s", tc.name)
				assert.Equal(t, tc.expectedEventType, event.EventType, "expected event type  to match for case: %s", tc.name)
			}
		})
	}
}

func TestGetEventsByPreviousID(t *testing.T) {
	db := setupTestDB(t)
	repo := event_repository.EventRepository{}

	var SubjectID uint = 1
	mainEvent := event.Event{
		Model:     gorm.Model{ID: 1},
		Name:      "main_event",
		EventType: event.RegionalStage,
		StartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
		EndDate:   time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
	}
	if err := db.Create(&mainEvent).Error; err != nil {
		return
	}

	tempEventsOlympiad := []event.Event{
		{
			Model: gorm.Model{
				ID:        2,
				CreatedAt: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				UpdatedAt: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
			},
			Name:            "Olympiad 1",
			EventType:       event.Olympiad,
			StartDate:       time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
			EndDate:         time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
			PreviousEventID: &mainEvent.ID,
			SubjectID:       &SubjectID,
		},
	}
	if err := db.Create(&tempEventsOlympiad).Error; err != nil {
		return
	}

	testCases := []struct {
		name           string
		searched_id    uint
		expectError    bool
		expectedEvents []event.Event
		expectedErrMsg string
	}{
		{
			name:           "Get events on existing PreviousEventID",
			searched_id:    mainEvent.ID,
			expectError:    false,
			expectedEvents: tempEventsOlympiad,
		},
		{
			name:           "Get events on not-existing PreviousEventID",
			searched_id:    99,
			expectError:    false,
			expectedEvents: []event.Event{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			events, err := repo.GetEventsByPreviousID(db, tc.searched_id)
			if tc.expectError {
				assert.Error(t, err, "expected an error for case: %s", tc.name)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err, "expected no error for case: %s", tc.name)
				assert.Equal(t, tc.expectedEvents, events, "expected start date event to match for case: %s", tc.name)
			}
		})
	}
}

func TestGetEventsByType(t *testing.T) {
	db := setupTestDB(t)
	repo := event_repository.EventRepository{}

	var SubjectID uint = 1
	var MainRegionalEventID uint = 1
	tempEventReg := []event.Event{
		{
			Model:     gorm.Model{ID: MainRegionalEventID},
			Name:      "RegionalStageMain",
			StartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
			EndDate:   time.Date(2024, time.December, 1, 1, 1, 1, 1, time.UTC),
			EventType: event.RegionalStage,
		}}
	tempEventsOly := []event.Event{

		{
			Model: gorm.Model{
				ID:        14,
				CreatedAt: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				UpdatedAt: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
			},
			Name:            "Olympiad1",
			StartDate:       time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
			EndDate:         time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
			EventType:       event.Olympiad,
			PreviousEventID: &MainRegionalEventID,
			SubjectID:       &SubjectID,
		},
		{
			Model: gorm.Model{
				ID:        15,
				CreatedAt: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				UpdatedAt: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
			},
			Name:            "Olympiad2",
			StartDate:       time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
			EndDate:         time.Date(2024, time.July, 1, 1, 1, 1, 1, time.UTC),
			EventType:       event.Olympiad,
			PreviousEventID: &MainRegionalEventID,
			SubjectID:       &SubjectID,
		},
	}
	if err := db.Create(&tempEventReg).Error; err != nil {
		return
	}
	if err := db.Create(&tempEventsOly).Error; err != nil {
		return
	}

	testCases := []struct {
		name           string
		search_type    event.EventType
		expectData     []event.Event
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:        "Get all olympiad",
			search_type: event.Olympiad,
			expectError: false,
			expectData:  tempEventsOly,
		},
		{
			name:        "get non-existing stage",
			search_type: event.Stage,
			expectData:  []event.Event{},
			expectError: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			events, err := repo.GetEventsByType(db, tc.search_type)
			if tc.expectError {
				assert.Error(t, err, "expected an error for case: %s", tc.name)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err, "expected no error for case: %s", tc.name)
				assert.Equal(t, tc.expectData, events, "expected start date event to match for case: %s", tc.name)
			}
		})
	}

}
