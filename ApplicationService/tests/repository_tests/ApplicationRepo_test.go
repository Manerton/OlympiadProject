package ApplicationHandlerTest

import (
	"main/internal/models"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/stretchr/testify/require"
)

// Инициализация тестовой базы данных в памяти
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Application{})
	require.NoError(t, err)
	return db
}

// // Тест на создание новой заявки
// func TestCreateApplication(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := applicationrepository.ApplicationRepository{}

// 	application := models.Application{UserID: uuid.New(), EventID: uuid.New()}
// 	err := repo.CreateApplication(db, &application)
// 	require.NoError(t, err)

// 	var result models.Application
// 	err = db.First(&result, application.ApplicationID).Error
// 	require.NoError(t, err)
// 	assert.Equal(t, application.UserID, result.UserID)
// 	assert.Equal(t, application.EventID, result.EventID)
// 	assert.Nil(t, result.Status) // Проверка, что статус по умолчанию nil
// }

// // Тест на получение заявки по ID
// func TestGetApplicationByID(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := applicationrepository.ApplicationRepository{}

// 	application := models.Application{ApplicationID: uuid.New(), UserID: uuid.New(), EventID: uuid.New(), Status: nil}
// 	db.Create(&application)

// 	result, err := repo.GetApplicationByID(db, application.ApplicationID)
// 	require.NoError(t, err)
// 	assert.Equal(t, application.ApplicationID, result.ApplicationID)
// 	assert.Equal(t, application.UserID, result.UserID)
// 	assert.Equal(t, application.EventID, result.EventID)
// 	assert.Nil(t, result.Status) // Статус должен быть nil
// }

// // Тест на получение всех заявок
// func TestGetAllApplications(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := applicationrepository.ApplicationRepository{}

// 	db.Create(&models.Application{UserID: uuid.New(), EventID: uuid.New()})
// 	db.Create(&models.Application{UserID: uuid.New(), EventID: uuid.New()})

// 	applications, err := repo.GetAllApplications(db)
// 	require.NoError(t, err)
// 	assert.Len(t, applications, 2)
// }

// // Тест на обновление статуса заявки
// func TestUpdateApplicationStatus(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := applicationrepository.ApplicationRepository{}

// 	application := models.Application{UserID: uuid.New(), EventID: uuid.New(), Status: nil}
// 	db.Create(&application)

// 	status := true
// 	err := repo.UpdateApplicationStatus(db, application.ApplicationID, &status)
// 	require.NoError(t, err)

// 	var result models.Application
// 	err = db.First(&result, application.ApplicationID).Error
// 	require.NoError(t, err)
// 	assert.NotNil(t, result.Status) // Статус обновлен
// 	assert.Equal(t, &status, result.Status)
// }

// // Тест на удаление заявки по ID
// func TestDeleteApplicationByID(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := applicationrepository.ApplicationRepository{}

// 	application := models.Application{UserID: uuid.New(), EventID: uuid.New()}
// 	db.Create(&application)

// 	err := repo.DeleteApplicationByID(db, application.ApplicationID)
// 	require.NoError(t, err)

// 	var result models.Application
// 	err = db.First(&result, application.ApplicationID).Error
// 	assert.Error(t, err) // Заявка должна быть удалена
// 	assert.Equal(t, gorm.ErrRecordNotFound, err)
// }
