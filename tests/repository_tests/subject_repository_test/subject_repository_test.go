package subject_repository_test

import (
	"main/internal/models/subject"
	"main/internal/repositories/subject_repository"
	"testing"

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

	// Миграция таблиц для тестирования
	if err := db.AutoMigrate(&subject.Subject{}); err != nil {
		t.Fatalf("failed to migrate tables: %v", err)
	}

	return db
}

func TestCreateSubject(t *testing.T) {
	db := setupTestDB(t)
	repo := subject_repository.SubjectRepository{}

	testCases := []struct {
		name           string
		name_subject   string
		id             uint
		expectError    bool
		expectedName   string
		expectedErrMsg string
		expectedID     uint
	}{
		{
			name:           "Auto created ID",
			name_subject:   "History",
			id:             0,
			expectError:    false,
			expectedName:   "",
			expectedErrMsg: "",
			expectedID:     1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testSubject := subject.Subject{ID: tc.id, Name: tc.name_subject}
			id, err := repo.CreateSubject(db, testSubject)
			if tc.expectError {

			} else {
				assert.NoError(t, err, "expected no error for case: %s", tc.name)
				assert.Equal(t, tc.expectedID, id, "expected subject id to match for case: %s", tc.name)
			}
		})
	}
}

func TestGetSubjectByID(t *testing.T) {
	db := setupTestDB(t)
	repo := subject_repository.SubjectRepository{}

	// Добавляем тестовые данные
	testSubject := subject.Subject{
		Name: "Mathematics",
	}
	if err := db.Create(&testSubject).Error; err != nil {
		t.Fatalf("failed to create test subject: %v", err)
	}

	// Определяем тестовые случаи
	testCases := []struct {
		name           string
		id             uint
		expectedName   string
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:         "Existing subject",
			id:           testSubject.ID,
			expectedName: "Mathematics",
			expectError:  false,
		},
		{
			name:           "Non-existing subject",
			id:             9999,
			expectError:    true,
			expectedErrMsg: "record not found",
		},
		{
			name:           "Invalid ID (0)",
			id:             0,
			expectError:    true,
			expectedErrMsg: "repositories.subject_repository.GetSubjectById: invalid ID 0",
		},
	}

	// Выполняем тесты
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := repo.GetSubjectByID(db, tc.id)

			// Проверяем наличие или отсутствие ошибки
			if tc.expectError {
				assert.Error(t, err, "expected an error for case: %s", tc.name)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err, "expected no error for case: %s", tc.name)
				assert.Equal(t, tc.expectedName, result.Name, "expected subject name to match for case: %s", tc.name)
			}
		})
	}
}
