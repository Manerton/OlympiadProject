package application_repository

import (
	"fmt"
	"main/internal/models"

	"gorm.io/gorm"
)

// Структура репозитория
type ApplicationRepository struct{}

// 1. Создание новой заявки
func (r *ApplicationRepository) CreateApplication(db *gorm.DB, application *models.Application) error {
	const op = "repositories.application_repository.CreateApplication"
	if err := db.Create(application).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// 2. Получение заявки по ID
func (r *ApplicationRepository) GetApplicationByID(db *gorm.DB, id uint) (models.Application, error) {
	const op = "repositories.application_repository.GetApplicationByID"
	if id == 0 {
		return models.Application{}, fmt.Errorf("%s: invalid ID %d", op, id)
	}
	var application models.Application
	if err := db.First(&application, id).Error; err != nil {
		return models.Application{}, fmt.Errorf("%s: %w", op, err)
	}
	return application, nil
}

// 3. Получение всех заявок
func (r *ApplicationRepository) GetAllApplications(db *gorm.DB) ([]models.Application, error) {
	const op = "repositories.application_repository.GetAllApplications"
	var applications []models.Application
	if err := db.Find(&applications).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return applications, nil
}

// 4. Обновление статуса заявки
func (r *ApplicationRepository) UpdateApplicationStatus(db *gorm.DB, id uint, status *bool) error {
	const op = "repositories.application_repository.UpdateApplicationStatus"
	if id == 0 {
		return fmt.Errorf("%s: invalid ID %d", op, id)
	}

	// Обновление статуса заявки
	if err := db.Model(&models.Application{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// 5. Удаление заявки по ID
func (r *ApplicationRepository) DeleteApplicationByID(db *gorm.DB, id uint) error {
	const op = "repositories.application_repository.DeleteApplicationByID"
	if id == 0 {
		return fmt.Errorf("%s: invalid ID %d", op, id)
	}

	if err := db.Delete(&models.Application{}, id).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
