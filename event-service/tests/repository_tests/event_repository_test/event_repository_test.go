package event_repository_test

import (
	"context"

	// "encoding/json"
	// "fmt"
	"main/internal/models/event"
	"main/internal/repositories/event_repository"
	"main/internal/storage/orm"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

func setupTestDB(t *testing.T) orm.ORM {
	// Создаем in-memory базу данных SQLite для тестирования
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	db.Exec("PRAGMA foreign_keys = ON")
	// create enum (types Events)
	err = db.Exec(`
			CREATE TABLE IF NOT EXISTS events (
				id UUID PRIMARY KEY,
				name VARCHAR(128) NOT NULL,
				start_date DATETIME NOT NULL,
				end_date DATETIME NOT NULL,
				event_type TEXT CHECK(event_type IN ('REGIONAL_STAGE', 'OLYMPIAD', 'STAGE', 'VIEW_WORKS', 'APPEAL')) NOT NULL,
				previous_event_id UUID REFERENCES events(id) ON DELETE RESTRICT ON UPDATE CASCADE,
				subject VARCHAR(128),
				additional_info TEXT
			)
		`).Error
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
	return &orm.Gorm{DB: db}
}

func TestCreateEvent(t *testing.T) {
	db := setupTestDB(t)
	repo := event_repository.EventRepository{}

	var Subject string = ""

	mainContext := context.Background()

	mainEvent := event.Event{
		ID:        uuid.New(),
		Name:      "MainEvent",
		StartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
		EndDate:   time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
		EventType: event.RegionalStage,
	}
	if err := db.Create(mainContext, &mainEvent); err != nil {
		return
	}

	testCases := []struct {
		name           string
		event          event.Event
		timer          time.Duration
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:  "Created correct event RegionalStage",
			timer: 10 * time.Millisecond,
			event: event.Event{
				Name:      "RegionalStage",
				StartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				EndDate:   time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
				EventType: event.RegionalStage,
			},
			expectError: false,
		},
		{
			name:  "Create event with unnecessary fields",
			timer: 10 * time.Millisecond,
			event: event.Event{
				Name:            "RegionalStage2",
				StartDate:       time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				EndDate:         time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
				EventType:       event.RegionalStage,
				Subject:         Subject,
				PreviousEventID: &mainEvent.ID,
			},
			expectError: false,
		},
		{
			name:  "Create event type Olympiad with correct foreignKey",
			timer: 10 * time.Millisecond,
			event: event.Event{
				Name:            "Olympiad",
				StartDate:       time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				EndDate:         time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
				EventType:       event.Olympiad,
				Subject:         Subject,
				PreviousEventID: &mainEvent.ID,
			},
			expectError: false,
		},
		{
			name:  "Create event type Olympiad with incorrect Date",
			timer: 10 * time.Millisecond,
			event: event.Event{
				Name:            "Olympiad",
				StartDate:       time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
				EndDate:         time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				EventType:       event.Olympiad,
				Subject:         Subject,
				PreviousEventID: &mainEvent.ID,
			},
			expectError: false,
		},
		{
			name:  "Create event type stage without fields SubjectID and PreviousID",
			timer: 10 * time.Millisecond,
			event: event.Event{
				Name:      "Stage",
				StartDate: time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
				EndDate:   time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				EventType: event.Stage,
			},
			expectError: false,
		},
		{
			name:  "Create event type appeal",
			timer: 10 * time.Millisecond,
			event: event.Event{
				Name:      "Appeal",
				StartDate: time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
				EndDate:   time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				EventType: event.Appeal,
			},
			expectError: false,
		},
		{
			name:           "Create event without fields",
			timer:          10 * time.Millisecond,
			event:          event.Event{},
			expectError:    true,
			expectedErrMsg: "CHECK constraint failed",
		},
		{
			name:  "Create event with incorrect PreviousID",
			timer: 10 * time.Millisecond,

			event: event.Event{
				Name:            "Appeal",
				StartDate:       time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
				EndDate:         time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
				EventType:       event.Appeal,
				PreviousEventID: &uuid.Max,
				Subject:         Subject,
			},
			expectError:    true,
			expectedErrMsg: "FOREIGN KEY constraint failed",
		},
		{
			name:  "Stop execution use context by timer",
			timer: time.Nanosecond,
			event: event.Event{
				EventType:       event.RegionalStage,
				PreviousEventID: nil,
			},
			expectError:    true,
			expectedErrMsg: "context deadline exceeded",
		},
		{
			name:  "Stop execution use context by cancel",
			timer: 100 * time.Millisecond,
			event: event.Event{
				EventType:       event.RegionalStage,
				PreviousEventID: nil,
			},
			expectError:    true,
			expectedErrMsg: "context canceled",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(mainContext, tc.timer)
			if tc.timer > 10*time.Millisecond {
				cancel()
			} else {
				defer cancel()
			}
			// create expected ID
			expectedID := uuid.New()
			tc.event.ID = expectedID
			// create event
			id, err := repo.CreateEvent(ctx, db, tc.event)
			if tc.expectError {
				assert.Error(t, err, "expected an error for case: %s", tc.name)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err, "no error expected for case: %s", tc.name)
				assert.Equal(t, expectedID, id, "expected event id to match for case: %s", tc.name)
			}
		})
	}
}

func TestGetEventByID(t *testing.T) {
	db := setupTestDB(t)
	repo := event_repository.EventRepository{}
	// create uuid
	uid := uuid.New()
	// test event
	tempEvent := event.Event{
		ID:        uid,
		Name:      "RegionalStage",
		StartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
		EndDate:   time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
		EventType: event.RegionalStage,
	}
	// create test data
	if err := db.Create(context.Background(), &tempEvent); err != nil {
		return
	}

	testCases := []struct {
		name              string
		searched_id       uuid.UUID
		expectError       bool
		expectedStartDate time.Time
		expectedEndDate   time.Time
		expectedEventType event.EventType
		expectedErrMsg    string
	}{
		{
			name:              "Get existing event",
			searched_id:       uid,
			expectedStartDate: time.Date(2024, time.January, 1, 1, 1, 1, 1, time.UTC),
			expectedEndDate:   time.Date(2024, time.February, 1, 1, 1, 1, 1, time.UTC),
			expectedEventType: event.RegionalStage,
			expectError:       false,
		},
		{
			name:           "Non-existing event",
			searched_id:    uuid.Max,
			expectError:    true,
			expectedErrMsg: "record not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event, err := repo.GetEventByID(context.Background(), db, tc.searched_id)
			if tc.expectError {
				assert.Error(t, err, "expected an error for case: %s", tc.name)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err, "no error expected for case: %s", tc.name)
				assert.Equal(t, tc.expectedStartDate, event.StartDate, "expected start date event to match for case: %s", tc.name)
				assert.Equal(t, tc.expectedEndDate, event.EndDate, "expected end date event to match for case: %s", tc.name)
				assert.Equal(t, tc.expectedEventType, event.EventType, "expected event type to match for case: %s", tc.name)
			}
		})
	}
}

func TestGetEventByListID(t *testing.T) {
	db := setupTestDB(t)
	repo := event_repository.EventRepository{}
	// data for test
	var testData []event.Event = []event.Event{
		{
			Name:      "A",
			EventType: event.RegionalStage,
		},
		{
			Name:      "B",
			EventType: event.RegionalStage,
		},
		{
			Name:      "C",
			EventType: event.RegionalStage,
		},
	}

	listID := make([]uuid.UUID, 0, len(testData))
	for i := range len(testData) {
		tempID := uuid.New()
		testData[i].ID = tempID
		listID = append(listID, tempID)
	}

	// create test data
	db.Create(context.Background(), testData)

	testCases := []struct {
		name           string
		serchedID      []uuid.UUID
		expectedData   []event.Event
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:         "Successful search",
			serchedID:    listID,
			expectedData: testData,
			expectError:  false,
		},
		{
			name:         "Search through partially existing IDs",
			serchedID:    []uuid.UUID{uuid.Max, uuid.Nil},
			expectedData: []event.Event{},
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			events, err := repo.GetEventsByListID(context.Background(), db, tc.serchedID)
			if tc.expectError {
				assert.Error(t, err, "expected an error for case: %s", tc.name)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err, "no error expected for case: %s", tc.name)
				for _, expectEvent := range tc.expectedData {
					expect := false
					for _, actualEvent := range events {
						if expectEvent.ID == actualEvent.ID {
							assert.Equal(t, expectEvent, actualEvent)
							expect = true
						}
					}
					if !expect {
						t.Error("expected data not equal actual data")
					}
				}
			}
		})
	}
}

func TestGetEventsByType(t *testing.T) {
	db := setupTestDB(t)
	repo := event_repository.EventRepository{}

	// Create TestData
	IdWithChild := uuid.New()
	testData := []event.Event{
		{
			ID:        uuid.New(),
			EventType: event.RegionalStage,
		},
		{
			ID:        uuid.New(),
			EventType: event.RegionalStage,
		},
		{
			ID:        IdWithChild,
			EventType: event.RegionalStage,
		},
		{
			ID:              uuid.New(),
			EventType:       event.Olympiad,
			PreviousEventID: &IdWithChild,
		},
	}
	db.Create(context.Background(), testData)
	// Create test cases
	testCases := []struct {
		name           string
		searchedType   string
		expectedCount  int
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:          "Success search existing events",
			searchedType:  string(event.RegionalStage),
			expectedCount: 3,
			expectError:   false,
		},
		{
			name:          "Success search existing single event",
			searchedType:  string(event.Olympiad),
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Success search non-existing events",
			searchedType:  "TEST",
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			events, err := repo.GetEventsByType(context.Background(), db, event.EventType(tc.searchedType), nil, nil, nil)
			if tc.expectError {
				assert.Error(t, err, "expected an error for case: %s", tc.name)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err, "no error expected for case: %s", tc.name)
				assert.Equal(t, len(events), tc.expectedCount, "expected count events: %v", tc.expectedCount)
			}
		})
	}
}

func TestGetEventsByPreviousID(t *testing.T) {
	db := setupTestDB(t)
	repo := event_repository.EventRepository{}

	mainEventID := uuid.New()
	testData := []event.Event{
		{
			ID:        mainEventID,
			EventType: event.RegionalStage,
		},
		{
			ID:              uuid.New(),
			EventType:       event.Olympiad,
			PreviousEventID: &mainEventID,
		},
		{
			ID:              uuid.New(),
			EventType:       event.Olympiad,
			PreviousEventID: &mainEventID,
		},
	}

	// create test data
	db.Create(context.Background(), testData)

	testCases := []struct {
		name              string
		searcedPreviousId uuid.UUID
		expectedCount     int
		expectError       bool
		expectedErrMsg    string
	}{
		{
			name:              "Success search existings events",
			searcedPreviousId: mainEventID,
			expectedCount:     2,
			expectError:       false,
		},
		{
			name:              "Search by not existing id",
			searcedPreviousId: uuid.Max,
			expectedCount:     0,
			expectError:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			events, err := repo.GetEventsByPreviousID(context.Background(), db, tc.searcedPreviousId, nil, nil, nil)
			if tc.expectError {
				assert.Error(t, err, "expected an error for case: %s", tc.name)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err, "no error expected for case: %s", tc.name)
				assert.Equal(t, len(events), tc.expectedCount, "expected count events: %v", tc.expectedCount)
			}
		})
	}
}

func TestGetAllEvents(t *testing.T) {
	db := setupTestDB(t)
	repo := event_repository.EventRepository{}

	testData := []event.Event{
		{
			ID:        uuid.New(),
			EventType: event.RegionalStage,
		},
		{
			ID:        uuid.New(),
			EventType: event.Olympiad,
		},
		{
			ID:        uuid.New(),
			EventType: event.Stage,
		},
		{
			ID:        uuid.New(),
			EventType: event.Appeal,
		},
		{
			ID:        uuid.New(),
			EventType: event.ViewWorks,
		},
	}

	db.Create(context.Background(), testData)

	events, err := repo.GetAllEvents(context.Background(), db, nil, nil)
	assert.NoError(t, err, "no error expected")
	assert.Equal(t, len(testData), len(events), "Expected len test data equal real len events")

}

func TestGetEventByFilterAndFields(t *testing.T) {
	db := setupTestDB(t)
	repo := event_repository.EventRepository{}

	specId := uuid.New()
	specStartDate := time.Now()
	specEndDate := time.Now().Add(1 * time.Hour)

	testData := []event.Event{
		{
			ID:        uuid.New(),
			Name:      "Test1",
			EventType: event.RegionalStage,
		},
		{
			ID:        uuid.New(),
			Name:      "Test2",
			EventType: event.RegionalStage,
			StartDate: specStartDate,
			EndDate:   specEndDate,
		},
		{
			ID:        specId,
			Name:      "Test3",
			EventType: event.Olympiad,
			Subject:   "math",
		},
	}

	db.Create(context.Background(), testData)

	testCases := []struct {
		name           string
		searchedFilter event.Event
		expectEvent    event.Event
		searchedFields []string
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:           "Search event by filter special name",
			searchedFilter: event.Event{Name: "Test1"},
			expectEvent:    testData[0],
			searchedFields: nil,
			expectError:    false,
		},
		{
			name:           "Search event by filter startDate",
			expectEvent:    testData[1],
			searchedFilter: event.Event{StartDate: specStartDate},
			searchedFields: nil,
			expectError:    false,
		},
		{
			name:           "Search event by filter id",
			expectEvent:    testData[2],
			searchedFilter: event.Event{ID: specId},
			searchedFields: nil,
			expectError:    false,
		},
		{
			name:           "Search event by filter and view specila fields",
			searchedFilter: event.Event{StartDate: specStartDate},
			expectEvent:    event.Event{Name: "Test2", StartDate: specStartDate},
			searchedFields: []string{"name", "start_date"},
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event, err := repo.GetEventByFilterAndFields(context.Background(), db, tc.searchedFilter, &tc.searchedFields)
			if tc.expectError {
				assert.Error(t, err, "expected an error for case: %s", tc.name)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err, "no error expected for case %s", tc.name)
				assert.True(t, eventsEqual(event, tc.expectEvent))
			}
		})
	}

}

func eventsEqual(e1, e2 event.Event) bool {
	// Сравниваем ID, Name, EventType, Subject, AdditionalInfo
	if e1.ID != e2.ID || e1.Name != e2.Name || e1.EventType != e2.EventType || e1.Subject != e2.Subject || e1.AdditionalInfo != e2.AdditionalInfo {
		return false
	}
	// Сравниваем StartDate и EndDate, игнорируя Location
	if !e1.StartDate.UTC().Equal(e2.StartDate.UTC()) || !e1.EndDate.UTC().Equal(e2.EndDate.UTC()) {
		return false
	}
	return true
}

func TestGetEventsByFilterAndFields(t *testing.T) {
	db := setupTestDB(t)
	repo := event_repository.EventRepository{}

	//
	mainEventId := uuid.New()

	testData := []event.Event{
		{
			ID:        mainEventId,
			Name:      "Test1.1",
			EventType: event.RegionalStage,
		},
		{
			ID:              uuid.New(),
			Name:            "Test2.1",
			EventType:       event.Olympiad,
			PreviousEventID: &mainEventId,
		},
		{
			ID:              uuid.New(),
			Name:            "Test2.2",
			EventType:       event.Olympiad,
			Subject:         "Math",
			PreviousEventID: &mainEventId,
		},
		{
			ID:        uuid.New(),
			Name:      "Test1.2",
			EventType: event.RegionalStage,
		},
		{
			ID:        uuid.New(),
			Name:      "Test3.1",
			EventType: event.Stage,
		},
		{
			ID:        uuid.New(),
			Name:      "Test3.2",
			EventType: event.Olympiad,
			Subject:   "Math",
		},
	}

	db.Create(context.Background(), testData)

	testCases := []struct {
		name           string
		searchedFilter event.Event
		searchedFields []string
		searchedLimit  int
		searchedOffset int
		searchedOrder  string
		expectData     []event.Event
		expectError    bool
		expectedErrMsg string
	}{
		{
			name: "Search by event type with order",
			searchedFilter: event.Event{
				EventType: event.RegionalStage,
			},
			searchedOrder: "name DESC",
			expectData:    []event.Event{testData[3], testData[0]},
			expectError:   false,
		},
		{
			name: "Search by subject with limit",
			searchedFilter: event.Event{
				Subject: "Math",
			},
			searchedLimit: 1,
			expectData:    []event.Event{testData[2]},
			expectError:   false,
		},
		{
			name:           "Search all with offset and order",
			searchedFilter: event.Event{},
			searchedFields: nil,
			searchedOffset: 2,
			searchedOrder:  "name",
			expectData:     []event.Event{testData[1], testData[2], testData[4], testData[5]},
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			events, err := repo.GetEventsByFilterAndFields(context.Background(), db,
				tc.searchedFilter, &tc.searchedFields, &tc.searchedOffset, &tc.searchedLimit, &tc.searchedOrder)
			if tc.expectError {
				assert.Error(t, err, "expected an error for case: %s", tc.name)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err, "no error expected for case %s", tc.name)
				assert.Equal(t, tc.expectData, events, "expected equals between actual and expected data")
			}
		})
	}

}

func TestUpdateEvent(t *testing.T) {
	db := setupTestDB(t)
	repo := event_repository.EventRepository{}

	mainId := uuid.New()
	mainOldStartDate := time.Now()
	mainNewStartDate := time.Now().Add(2 * time.Hour)
	mainEndDate := time.Now().Add(10 * time.Hour)

	testData := []event.Event{
		{
			ID:        mainId,
			Name:      "Test",
			EventType: event.RegionalStage,
			StartDate: mainOldStartDate,
			EndDate:   mainEndDate,
			Subject:   "Math",
		},
		{
			ID:              uuid.New(),
			Name:            "Empty event",
			EventType:       event.Olympiad,
			PreviousEventID: &mainId,
		},
	}

	db.Create(context.Background(), testData)

	expectedEvent := testData[0]
	expectedEvent.StartDate = mainNewStartDate

	testCases := []struct {
		name           string
		newEvent       event.Event
		expectedEvent  event.Event
		expectError    bool
		expectedErrMsg string
	}{
		{
			name: "Success update existring event",
			newEvent: event.Event{
				ID:        mainId,
				StartDate: mainNewStartDate,
			},
			expectedEvent: expectedEvent,
			expectError:   false,
		},
		{
			name: "Failed attempt to update a non-existent event",
			newEvent: event.Event{
				ID:             uuid.New(),
				AdditionalInfo: "TEST",
			},
			expectError:    true,
			expectedErrMsg: "record not found for update",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := repo.UpdateEvent(context.Background(), db, tc.newEvent)
			if tc.expectError {
				assert.Error(t, err, "expected an error for case: %s", tc.name)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err, "")
				updatedEvent, err := repo.GetEventByID(context.Background(), db, id)
				if err != nil {
					t.Error("failed when get updated event")
				}
				assert.True(t, eventsEqual(updatedEvent, tc.expectedEvent))
			}

		})
	}

}

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
