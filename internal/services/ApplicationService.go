package ApplicationService

import (
	"fmt"
	"main/internal/dto"
	"main/internal/models"
	application_repository "main/internal/repositories"

	"gorm.io/gorm"
)

type ApplicationService struct {
	db         *gorm.DB
	repository *application_repository.ApplicationRepository
}

func NewApplicationService(db *gorm.DB, repo *application_repository.ApplicationRepository) *ApplicationService {
	return &ApplicationService{
		db:         db,
		repository: repo,
	}
}

// Получение всех заявок
func (s *ApplicationService) GetAllApplications() ([]dto.ApplicationResponseDTO, error) {
	const op = "services.application_service.GetAllApplications"
	applications, err := s.repository.GetAllApplications(s.db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertManyApplicationsToDTO(applications), nil
}

// Получение заявки по ID
func (s *ApplicationService) GetApplicationByID(id uint) (dto.ApplicationResponseDTO, error) {
	const op = "services.application_service.GetApplicationByID"
	application, err := s.repository.GetApplicationByID(s.db, id)
	if err != nil {
		return dto.ApplicationResponseDTO{}, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertApplicationToDTO(application), nil
}

// Создание новой заявки
func (s *ApplicationService) CreateApplication(applicationDTO dto.CreateApplicationDTO) (uint, error) {
	const op = "services.application_service.CreateApplication"
	application := ConvertDTOtoApplication(applicationDTO)
	if err := s.repository.CreateApplication(s.db, &application); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return application.ApplicationID, nil
}

// Обновление статуса заявки
func (s *ApplicationService) UpdateApplicationStatus(id uint, statusDTO dto.UpdateApplicationStatusDTO) error {
	const op = "services.application_service.UpdateApplicationStatus"
	if err := s.repository.UpdateApplicationStatus(s.db, id, statusDTO.Status); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Удаление заявки по ID
func (s *ApplicationService) DeleteApplication(id uint) error {
	const op = "services.application_service.DeleteApplication"
	if err := s.repository.DeleteApplicationByID(s.db, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Функции для преобразования между DTO и моделью

func ConvertDTOtoApplication(dto dto.CreateApplicationDTO) models.Application {
	return models.Application{
		UserID:  dto.UserID,
		EventID: dto.EventID,
	}
}

func ConvertApplicationToDTO(application models.Application) dto.ApplicationResponseDTO {
	return dto.ApplicationResponseDTO{
		ApplicationID: application.ApplicationID,
		UserID:        application.UserID,
		EventID:       application.EventID,
		Status:        application.Status,
		SubmittedAt:   application.SubmittedAt,
		UpdatedAt:     application.UpdatedAt,
	}
}

func ConvertManyApplicationsToDTO(applications []models.Application) []dto.ApplicationResponseDTO {
	var applicationsDTO []dto.ApplicationResponseDTO
	for _, application := range applications {
		applicationsDTO = append(applicationsDTO, ConvertApplicationToDTO(application))
	}
	return applicationsDTO
}
