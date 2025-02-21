package event_repository_test

// import (
// 	"fmt"
// 	"main/internal/models/event"
// 	"main/internal/repositories/event_repository"
// 	"main/internal/storage/orm"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// func setupTestDB(t *testing.T) orm.ORM {
// 	// Создаем in-memory базу данных SQLite для тестирования
// 	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
// 	if err != nil {
// 		t.Fatalf("failed to connect database: %v", err)
// 	}
// 	db.Exec("PRAGMA foreign_keys = ON")
// 	// create enum (types Events)
// 	create_enum_type_str := fmt.Sprintf(
// 		`DO $$ BEGIN IF NOT EXISTS
// 		(SELECT 1 FROM pg_type WHERE typname = 'event_type')
// 		THEN CREATE TYPE event_type AS ENUM
// 		('%s', '%s', '%s', '%s', '%s');
// 		END IF; END $$;`,
// 		event.RegionalStage,
// 		event.Olympiad,
// 		event.Stage,
// 		event.ViewWorks,
// 		event.Appeal,
// 	)
// 	db.Exec(create_enum_type_str)

// 	// Миграция таблиц для тестирования
// 	if err := db.AutoMigrate(&event.Event{}); err != nil {
// 		t.Fatalf("failed to migrate tables: %v", err)
// 	}

// 	// Создаём тестовые данные для вспомогатльной таблицы subject
// 	// testSubjects := []subject.Subject{
// 	// 	{
// 	// 		ID:   1,
// 	// 		Name: "Test 1",
// 	// 	},
// 	// 	{
// 	// 		ID:   2,
// 	// 		Name: "Test 2",
// 	// 	},
// 	// 	{
// 	// 		ID:   3,
// 	// 		Name: "Test 3",
// 	// 	},
// 	// }
// 	// if err = db.Create(&testSubjects).Error; err != nil {
// 	// 	t.Fatalf("failed to create support subjects: %v", err)
// 	// }

// 	return &orm.Gorm{DB: db}
// }

// func TestCreateEvent(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := event_repository.EventRepository{}

// 	var Subject string = ""
// 	var UnCorrectPreviousEventID uint = 99

// 	main_event := event.Event{
// 		ID:        1,
// 		Name:      "MainEvent",
// 		StartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 		EndDate:   time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 		EventType: event.RegionalStage,
// 	}
// 	if err := db.Create(&main_event).Error; err != nil {
// 		return
// 	}

// 	testCases := []struct {
// 		name           string
// 		event          event.Event
// 		expectError    bool
// 		expectedErrMsg string
// 		expectedId     uint
// 	}{
// 		{
// 			name: "Created correct event RegionalStage",
// 			event: event.Event{
// 				ID:        2,
// 				Name:      "RegionalStage",
// 				StartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 				EndDate:   time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 				EventType: event.RegionalStage,
// 			},
// 			expectedId:  2,
// 			expectError: false,
// 		},
// 		{
// 			name: "Create event with unnecessary fields",
// 			event: event.Event{
// 				ID:              3,
// 				Name:            "RegionalStage2",
// 				StartDate:       time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 				EndDate:         time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 				EventType:       event.RegionalStage,
// 				Subject:         Subject,
// 				PreviousEventID: &main_event.ID,
// 			},
// 			expectedId:  3,
// 			expectError: false,
// 		},
// 		{
// 			name: "Create event type Olympiad with correct foreignKey",
// 			event: event.Event{
// 				ID:              4,
// 				Name:            "Olympiad",
// 				StartDate:       time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 				EndDate:         time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 				EventType:       event.Olympiad,
// 				Subject:         Subject,
// 				PreviousEventID: &main_event.ID,
// 			},
// 			expectedId:  4,
// 			expectError: false,
// 		},
// 		{
// 			name: "Create event type Olympiad with incorrect Date",
// 			event: event.Event{
// 				ID:              5,
// 				Name:            "Olympiad",
// 				StartDate:       time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 				EndDate:         time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 				EventType:       event.Olympiad,
// 				Subject:         Subject,
// 				PreviousEventID: &main_event.ID,
// 			},
// 			expectedId:  5,
// 			expectError: false,
// 		},
// 		{
// 			name: "Create event type stage without fields SubjectID/PreviousID",
// 			event: event.Event{
// 				ID:        6,
// 				Name:      "Stage",
// 				StartDate: time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 				EndDate:   time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 				EventType: event.Stage,
// 			},
// 			expectedId:  6,
// 			expectError: false,
// 		},
// 		{
// 			name: "Create event type appeal",
// 			event: event.Event{
// 				ID:        7,
// 				Name:      "Appeal",
// 				StartDate: time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 				EndDate:   time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 				EventType: event.Appeal,
// 			},
// 			expectedId:  7,
// 			expectError: false,
// 		},
// 		{
// 			name: "Create event without fields",
// 			event: event.Event{
// 				ID: 8,
// 			},
// 			expectedId:  8,
// 			expectError: false,
// 		},
// 		{
// 			name: "Create event with incorrect PreviousID",
// 			event: event.Event{
// 				ID:              9,
// 				Name:            "Appeal",
// 				StartDate:       time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 				EndDate:         time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 				EventType:       event.Appeal,
// 				PreviousEventID: &UnCorrectPreviousEventID,
// 				Subject:         Subject,
// 			},
// 			expectError:    true,
// 			expectedErrMsg: "FOREIGN KEY constraint failed",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			id, err := repo.CreateEvent(db, tc.event)
// 			if tc.expectError {
// 				assert.Error(t, err, "expected an error for case: %s", tc.name)
// 				if err != nil {
// 					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
// 				}
// 			} else {
// 				assert.NoError(t, err, "expected no error for case: %s", tc.name)
// 				assert.Equal(t, tc.expectedId, id, "expected event id to match for case: %s", tc.name)

// 			}
// 		})
// 	}
// }

// func TestGetEventByID(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := event_repository.EventRepository{}

// 	tempEvent := event.Event{
// 		Name:      "RegionalStage",
// 		StartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 		EndDate:   time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 		EventType: event.RegionalStage,
// 	}
// 	if err := db.Create(&tempEvent).Error; err != nil {
// 		return
// 	}

// 	testCases := []struct {
// 		name              string
// 		searched_id       uint
// 		expectError       bool
// 		expectedStartDate time.Time
// 		expectedEndDate   time.Time
// 		expectedEventType event.EventType
// 		expectedErrMsg    string
// 	}{
// 		{
// 			name:              "Get existing event",
// 			searched_id:       1,
// 			expectedStartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 			expectedEndDate:   time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 			expectedEventType: event.RegionalStage,
// 			expectError:       false,
// 		},
// 		{
// 			name:              "Get first event on default id = 0",
// 			searched_id:       0,
// 			expectedStartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 			expectedEndDate:   time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 			expectedEventType: event.RegionalStage,
// 			expectError:       false,
// 		},
// 		{
// 			name:           "Non-existing event",
// 			searched_id:    9999,
// 			expectError:    true,
// 			expectedErrMsg: "record not found",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			event, err := repo.GetEventByID(db, tc.searched_id)
// 			if tc.expectError {
// 				assert.Error(t, err, "expected an error for case: %s", tc.name)
// 				if err != nil {
// 					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
// 				}
// 			} else {
// 				assert.NoError(t, err, "expected no error for case: %s", tc.name)
// 				assert.Equal(t, tc.expectedStartDate, event.StartDate, "expected start date event to match for case: %s", tc.name)
// 				assert.Equal(t, tc.expectedEndDate, event.EndDate, "expected end date event to match for case: %s", tc.name)
// 				assert.Equal(t, tc.expectedEventType, event.EventType, "expected event type  to match for case: %s", tc.name)
// 			}
// 		})
// 	}
// }

// func TestGetEventsByPreviousID(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := event_repository.EventRepository{}

// 	var Subject string = "Математика"
// 	mainEvent := event.Event{
// 		ID:        1,
// 		Name:      "main_event",
// 		EventType: event.RegionalStage,
// 		StartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 		EndDate:   time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 	}
// 	if err := db.Create(&mainEvent).Error; err != nil {
// 		return
// 	}

// 	tempEventsOlympiad := []event.Event{
// 		{
// 			ID:              2,
// 			Name:            "Olympiad 1",
// 			EventType:       event.Olympiad,
// 			PreviousEventID: &mainEvent.ID,
// 			Subject:         Subject,
// 		},
// 	}
// 	if err := db.Create(&tempEventsOlympiad).Error; err != nil {
// 		return
// 	}

// 	testCases := []struct {
// 		name           string
// 		searched_id    uint
// 		expectError    bool
// 		expectedEvents []event.Event
// 		expectedErrMsg string
// 	}{
// 		{
// 			name:           "Get events on existing PreviousEventID",
// 			searched_id:    mainEvent.ID,
// 			expectError:    false,
// 			expectedEvents: tempEventsOlympiad,
// 		},
// 		{
// 			name:           "Get events on not-existing PreviousEventID",
// 			searched_id:    99,
// 			expectError:    false,
// 			expectedEvents: []event.Event{},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			events, err := repo.GetEventsByPreviousID(db, tc.searched_id, nil, nil)
// 			if tc.expectError {
// 				assert.Error(t, err, "expected an error for case: %s", tc.name)
// 				if err != nil {
// 					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
// 				}
// 			} else {
// 				assert.NoError(t, err, "expected no error for case: %s", tc.name)
// 				assert.Equal(t, tc.expectedEvents, events, "expected start date event to match for case: %s", tc.name)
// 			}
// 		})
// 	}
// }

// func TestGetEventsByType(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := event_repository.EventRepository{}

// 	var Subject string = "Математика"
// 	var MainRegionalEventID uint = 1
// 	tempEventReg := []event.Event{
// 		{
// 			ID:        MainRegionalEventID,
// 			Name:      "RegionalStageMain",
// 			StartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 			EndDate:   time.Date(2024, time.December, 1, 1, 1, 1, 1, time.UTC),
// 			EventType: event.RegionalStage,
// 		}}
// 	tempEventsOly := []event.Event{

// 		{
// 			ID:              14,
// 			Name:            "Olympiad1",
// 			StartDate:       time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 			EndDate:         time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 			EventType:       event.Olympiad,
// 			PreviousEventID: &MainRegionalEventID,
// 			Subject:         Subject,
// 		},
// 		{
// 			ID:              15,
// 			Name:            "Olympiad2",
// 			StartDate:       time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 			EndDate:         time.Date(2024, time.July, 1, 1, 1, 1, 1, time.UTC),
// 			EventType:       event.Olympiad,
// 			PreviousEventID: &MainRegionalEventID,
// 			Subject:         Subject,
// 		},
// 	}
// 	if err := db.Create(&tempEventReg).Error; err != nil {
// 		return
// 	}
// 	if err := db.Create(&tempEventsOly).Error; err != nil {
// 		return
// 	}

// 	testCases := []struct {
// 		name           string
// 		search_type    event.EventType
// 		expectData     []event.Event
// 		expectError    bool
// 		expectedErrMsg string
// 	}{
// 		{
// 			name:        "Get all olympiad",
// 			search_type: event.Olympiad,
// 			expectError: false,
// 			expectData:  tempEventsOly,
// 		},
// 		{
// 			name:        "get non-existing stage",
// 			search_type: event.Stage,
// 			expectData:  []event.Event{},
// 			expectError: false,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			events, err := repo.GetEventsByType(db, tc.search_type, nil, nil)
// 			if tc.expectError {
// 				assert.Error(t, err, "expected an error for case: %s", tc.name)
// 				if err != nil {
// 					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
// 				}
// 			} else {
// 				assert.NoError(t, err, "expected no error for case: %s", tc.name)
// 				assert.Equal(t, tc.expectData, events, "expected start date event to match for case: %s", tc.name)
// 			}
// 		})
// 	}

// }

// func TestGetAllEvents(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := event_repository.EventRepository{}

// 	var Subject string = "Математика"
// 	var MainRegionalEventID uint = 1
// 	tempEvents := []event.Event{
// 		{
// 			ID:        MainRegionalEventID,
// 			Name:      "RegionalStageMain",
// 			StartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 			EndDate:   time.Date(2024, time.December, 1, 1, 1, 1, 1, time.UTC),
// 			EventType: event.RegionalStage,
// 		},
// 		{
// 			ID:              14,
// 			Name:            "Olympiad1",
// 			StartDate:       time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 			EndDate:         time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 			EventType:       event.Olympiad,
// 			PreviousEventID: &MainRegionalEventID,
// 			Subject:         Subject,
// 		},
// 		{
// 			ID:              15,
// 			Name:            "Olympiad2",
// 			StartDate:       time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 			EndDate:         time.Date(2024, time.July, 1, 1, 1, 1, 1, time.UTC),
// 			EventType:       event.Olympiad,
// 			PreviousEventID: &MainRegionalEventID,
// 			Subject:         Subject,
// 		},
// 	}

// 	if err := db.Create(&tempEvents).Error; err != nil {
// 		return
// 	}
// 	events, err := repo.GetAllEvents(db, nil, nil)

// 	assert.NoError(t, err)
// 	assert.Equal(t, tempEvents, events)
// }

// func TestUpdateEvent(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := event_repository.EventRepository{}

// 	var Subject string = "Математика"
// 	var NewSubject string = "Руссикй язык"
// 	var MainRegionalEventID uint = 1
// 	tempEvents := []event.Event{
// 		{
// 			ID:        MainRegionalEventID,
// 			Name:      "RegionalStageMain",
// 			StartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 			EndDate:   time.Date(2024, time.December, 1, 1, 1, 1, 1, time.UTC),
// 			EventType: event.RegionalStage,
// 		},
// 		{
// 			ID:              14,
// 			Name:            "Olympiad1",
// 			StartDate:       time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 			EndDate:         time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 			EventType:       event.Olympiad,
// 			PreviousEventID: &MainRegionalEventID,
// 			Subject:         Subject,
// 		},
// 		{
// 			ID:              15,
// 			Name:            "Olympiad2",
// 			StartDate:       time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 			EndDate:         time.Date(2024, time.July, 1, 1, 1, 1, 1, time.UTC),
// 			EventType:       event.Olympiad,
// 			PreviousEventID: &MainRegionalEventID,
// 			Subject:         Subject,
// 		},
// 	}
// 	if err := db.Create(&tempEvents).Error; err != nil {
// 		return
// 	}
// 	testCases := []struct {
// 		name           string
// 		update_id      uint
// 		update_data    event.Event
// 		expectError    bool
// 		expectedErrMsg string
// 	}{
// 		{
// 			name:      "Update name",
// 			update_id: 14,
// 			update_data: event.Event{
// 				Name:    "New Olympiad",
// 				Subject: NewSubject,
// 			},
// 			expectError: false,
// 		},
// 		{
// 			name:      "Update subject",
// 			update_id: 15,
// 			update_data: event.Event{
// 				Name:    "Subject Test",
// 				Subject: NewSubject,
// 			},
// 			expectError: false,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			id, err := repo.UpdateEvent(db, event.Event{
// 				ID:      tc.update_id,
// 				Name:    tc.update_data.Name,
// 				Subject: tc.update_data.Subject,
// 			})
// 			if tc.expectError {
// 				assert.Error(t, err, "expected an error for case: %s", tc.name)
// 				if err != nil {
// 					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
// 				}
// 			} else {
// 				assert.NoError(t, err, "expected no error for case: %s", tc.name)
// 				assert.Equal(t, tc.update_id, id, "expected start date event to match for case: %s", tc.name)
// 			}
// 		})
// 	}
// }

// func TestDeleteEvent(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := event_repository.EventRepository{}

// 	var Subject string = "Математика"
// 	var MainRegionalEventID uint = 1
// 	tempEvents := []event.Event{
// 		{
// 			ID:        MainRegionalEventID,
// 			Name:      "RegionalStageMain",
// 			StartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 			EndDate:   time.Date(2024, time.December, 1, 1, 1, 1, 1, time.UTC),
// 			EventType: event.RegionalStage,
// 		},
// 		{
// 			ID:              14,
// 			Name:            "Olympiad1",
// 			StartDate:       time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 			EndDate:         time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 			EventType:       event.Olympiad,
// 			PreviousEventID: &MainRegionalEventID,
// 			Subject:         Subject,
// 		},
// 		{
// 			ID:              15,
// 			Name:            "Olympiad2",
// 			StartDate:       time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
// 			EndDate:         time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
// 			EventType:       event.Olympiad,
// 			PreviousEventID: &MainRegionalEventID,
// 			Subject:         Subject,
// 		},
// 	}
// 	if err := db.Create(&tempEvents).Error; err != nil {
// 		return
// 	}
// 	testCases := []struct {
// 		name           string
// 		delete_id      uint
// 		expectError    bool
// 		expectedErrMsg string
// 	}{
// 		{
// 			name:        "Delete existing event",
// 			delete_id:   14,
// 			expectError: false,
// 		},
// 		{
// 			name:           "Delete event with child",
// 			delete_id:      MainRegionalEventID,
// 			expectError:    true,
// 			expectedErrMsg: "",
// 		},
// 		{
// 			name:           "Delete non-existing event",
// 			delete_id:      99,
// 			expectError:    true,
// 			expectedErrMsg: "",
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			err := repo.DeleteEvent(db, tc.delete_id)
// 			events, err2 := repo.GetAllEvents(db, nil, nil)
// 			_, _ = events, err2
// 			if tc.expectError {
// 				assert.Error(t, err, "expected an error for case: %s", tc.name)
// 				if err != nil {
// 					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
// 				}
// 			} else {
// 				assert.NoError(t, err, "expected no error for case: %s", tc.name)
// 			}
// 		})
// 	}
// }
