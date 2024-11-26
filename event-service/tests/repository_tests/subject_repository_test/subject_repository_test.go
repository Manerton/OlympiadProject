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
			name:         "Auto created ID",
			name_subject: "History",
			id:           0,
			expectError:  false,
			expectedID:   1,
		},
		{
			name:           "Created with my ID",
			name_subject:   "Mathematics",
			id:             2,
			expectError:    false,
			expectedErrMsg: "",
			expectedID:     2,
		},
		{
			name:           "Created existing id",
			name_subject:   "test",
			id:             2,
			expectError:    true,
			expectedErrMsg: "UNIQUE constraint failed: subjects.id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testSubject := subject.Subject{ID: tc.id, Name: tc.name_subject}
			id, err := repo.CreateSubject(db, testSubject)
			if tc.expectError {
				assert.Error(t, err, "expected an error for case: %s", tc.name)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
				}
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

func TestUpdatedSubjectID(t *testing.T) {
	db := setupTestDB(t)
	repo := subject_repository.SubjectRepository{}

	testSubjects := []subject.Subject{
		{
			ID:   1,
			Name: "Mathematics",
		},
		{
			ID:   2,
			Name: "History",
		},
		{
			ID:   3,
			Name: "TEST",
		},
	}
	if err := db.Create(&testSubjects).Error; err != nil {
		t.Fatalf("failed to create test subject: %v", err)
	}

	testCases := []struct {
		name           string
		new_name       string
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:        "Update name",
			new_name:    "NEW NAME",
			expectError: false,
		},
		{
			name:           "Update on existing name",
			new_name:       "TEST",
			expectError:    true,
			expectedErrMsg: "UNIQUE constraint failed",
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := repo.UpdateSubject(db, subject.Subject{ID: testSubjects[i].ID, Name: tc.new_name})
			subject, _ := repo.GetSubjectByID(db, id)

			// Проверяем наличие или отсутствие ошибки
			if tc.expectError {
				assert.Error(t, err, "expected an error for case: %s", tc.name)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err, "expected no error for case: %s", tc.name)
				assert.Equal(t, tc.new_name, subject.Name, "expected subject name to match for case: %s", tc.name)
			}
		})
	}
}

func TestDeleteSubject(t *testing.T) {
	db := setupTestDB(t)
	repo := subject_repository.SubjectRepository{}

	testSubject := subject.Subject{
		ID:   1,
		Name: "Mathematics",
	}
	if err := db.Create(&testSubject).Error; err != nil {
		t.Fatalf("failed to create test subject: %v", err)
	}

	testCases := []struct {
		name           string
		deleted_id     uint
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:        "deleted subject",
			deleted_id:  1,
			expectError: false,
		},
		{
			name:        "deleted non-existing subject",
			deleted_id:  1,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.DeleteSubject(db, tc.deleted_id)

			// Проверяем наличие или отсутствие ошибки
			if tc.expectError {
				assert.Error(t, err, "expected an error for case: %s", tc.name)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrMsg, "expected error message to contain: %s", tc.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err, "expected no error for case: %s", tc.name)
			}
		})
	}
}
