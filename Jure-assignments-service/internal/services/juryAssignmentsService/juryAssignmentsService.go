package juryAssignmentsService

import (
	"fmt"
	"main/internal/dto/juryAssignmentsDto"
	"main/internal/lib/converter/dtoConverter"
	"main/internal/lib/supportRequest"
	"main/internal/models/juryAssignments"
	"sync"

	"gorm.io/gorm"
)

type juryAssignmentsRepositoryInterface interface {
	GetJuryAssignmentsByFilter(db *gorm.DB, filter juryAssignments.JuryAssignments) (juryAssignments.JuryAssignments, error)
	GetAllJuryAssignments(*gorm.DB) ([]juryAssignments.JuryAssignments, error)
	GetAllJuryAssignmentsByFilter(db *gorm.DB, filter juryAssignments.JuryAssignments) ([]juryAssignments.JuryAssignments, error)
	GetPartOfAllJuryAssignmentsByFilter(
		db *gorm.DB, fields []string, filter juryAssignments.JuryAssignments) ([]juryAssignments.JuryAssignments, error)
	CreateJuryAssignments(
		db *gorm.DB, juryAssignments juryAssignments.JuryAssignments) (uint, error)
	UpdateJuryAssignments(
		db *gorm.DB, juryAssignments juryAssignments.JuryAssignments) (uint, error)
	DeleteJuryAssignments(db *gorm.DB, id uint) error
}

type JuryAssignmentsService struct {
	db         *gorm.DB
	repository juryAssignmentsRepositoryInterface
}

func NewJuryAssignmentsService(db *gorm.DB, jr juryAssignmentsRepositoryInterface) *JuryAssignmentsService {
	return &JuryAssignmentsService{db: db, repository: jr}
}

func (s *JuryAssignmentsService) GetAllJuryAssignments() ([]juryAssignmentsDto.JuryAssignmentsDTO, error) {
	const op = "services.juryAssignmentsService.GetAllJuryAssignments"
	results, err := s.repository.GetAllJuryAssignments(s.db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return dtoConverter.ConvertManyJuryAssignmentsToDTO(results), nil
}

func (s *JuryAssignmentsService) GetAllJuryAssignmentsByFilter(
	filter juryAssignmentsDto.JuryAssignmentsDTO) ([]juryAssignmentsDto.JuryAssignmentsDTO, error) {
	const op = "services.juryAssignmentsService.GetAllEventsByFilter"
	modelFilter := dtoConverter.ConvertDTOtoJuryAssignments(filter)
	results, err := s.repository.GetAllJuryAssignmentsByFilter(s.db, modelFilter)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return dtoConverter.ConvertManyJuryAssignmentsToDTO(results), nil
}

func (s *JuryAssignmentsService) GetPartOfAllJuryAssignmentsByFilter(
	fields []string, filter juryAssignmentsDto.JuryAssignmentsDTO) ([]juryAssignmentsDto.JuryAssignmentsDTO, error) {
	const op = "services.juryAssignmentsService.GetPartOfAllJuryAssignmentsByFilter"
	modelFilter := dtoConverter.ConvertDTOtoJuryAssignments(filter)
	result, err := s.repository.GetPartOfAllJuryAssignmentsByFilter(s.db, fields, modelFilter)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return dtoConverter.ConvertManyJuryAssignmentsToDTO(result), nil
}

func (s *JuryAssignmentsService) GetJuryAssignmentsByFilter(filter juryAssignmentsDto.JuryAssignmentsDTO) (juryAssignmentsDto.JuryAssignmentsDTO, error) {
	const op = "services.juryAssignmentsService.GetJuryAssignmentsByFilter"

	modelFilter := dtoConverter.ConvertDTOtoJuryAssignments(filter)
	result, err := s.repository.GetJuryAssignmentsByFilter(s.db, modelFilter)
	if err != nil {
		return juryAssignmentsDto.JuryAssignmentsDTO{}, fmt.Errorf("%s: %w", op, err)
	}
	return dtoConverter.ConvertJuryAssignmentsToDTO(result), nil
}

func isExistingUserID(id uint) bool {
	return true
	if id == 0 {
		return false
	}
	path := ""
	response, err := supportRequest.SupportRequest("GET", path)
	if err != nil {
		return false
	}
	if response.Data != nil {
		return true
	}
	return false
}

func isExistingEventID(id uint) bool {
	return true

	if id == 0 {
		return false
	}
	path := ""
	response, err := supportRequest.SupportRequest("GET", path)
	if err != nil {
		return false
	}
	if response.Data != nil {
		return true
	}
	return false
}

func (s *JuryAssignmentsService) CreateManyAssignmentsByOneJury(dto juryAssignmentsDto.OneJuryManyAssignments) ([]uint, []error) {
	wg := sync.WaitGroup{}
	errors := make(chan error, len(dto.EventsID))
	ids := make(chan uint, len(dto.EventsID))

	for _, eventID := range dto.EventsID {
		wg.Add(1)
		tempDto := juryAssignmentsDto.JuryAssignmentsDTO{JuryID: dto.JuryID, EventID: eventID}
		go func(goDto juryAssignmentsDto.JuryAssignmentsDTO) {
			defer wg.Done()
			id, err := s.CreateJuryAssignments(goDto)
			if err != nil {
				errors <- err
				return
			}
			ids <- id
		}(tempDto)
		fmt.Println()
	}
	wg.Wait()
	close(errors)
	close(ids)

	sliceErrors := []error{}
	wg.Add(1)
	go func(sliceErrors *[]error) {
		defer wg.Done()
		for err := range errors {
			*sliceErrors = append(*sliceErrors, err)
		}
	}(&sliceErrors)

	sliceIds := []uint{}
	wg.Add(1)
	go func(sliceIds *[]uint) {
		defer wg.Done()
		for id := range ids {
			*sliceIds = append(*sliceIds, id)
		}
	}(&sliceIds)
	wg.Wait()

	return sliceIds, sliceErrors
}

func (s *JuryAssignmentsService) CreateJuryAssignments(dto juryAssignmentsDto.JuryAssignmentsDTO) (uint, error) {
	const op = "services.juryAssignmentsService.CreateJuryAssignments"
	if !isExistingUserID(dto.JuryID) {
		return 0, fmt.Errorf("%s: Jury is not exist", op)
	}
	if !isExistingEventID(dto.EventID) {
		return 0, fmt.Errorf("%s: Event is not exist", op)
	}

	model := dtoConverter.ConvertDTOtoJuryAssignments(dto)
	id, err := s.repository.CreateJuryAssignments(s.db, model)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *JuryAssignmentsService) UpdateJuryAssignments(dto juryAssignmentsDto.JuryAssignmentsDTO) (uint, error) {
	const op = "services.juryAssignmentsService.CreateJuryAssignments"

	if !isExistingUserID(dto.JuryID) {
		return 0, fmt.Errorf("%s: Jury is not exist", op)
	}
	if !isExistingEventID(dto.EventID) {
		return 0, fmt.Errorf("%s: Event is not exist", op)
	}

	model := dtoConverter.ConvertDTOtoJuryAssignments(dto)
	id, err := s.repository.UpdateJuryAssignments(s.db, model)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *JuryAssignmentsService) DeleteJuryAssignments(id uint) error {
	const op = "services.juryAssignmentsService.CreateJuryAssignments"
	err := s.repository.DeleteJuryAssignments(s.db, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
