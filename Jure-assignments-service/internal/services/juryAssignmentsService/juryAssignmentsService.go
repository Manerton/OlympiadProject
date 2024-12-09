package juryAssignmentsService

import (
	"fmt"
	"main/internal/dto/juryAssignmentsDto"
	"main/internal/lib/converter/dtoConverter"
	"main/internal/lib/supportRequest"
	"main/internal/models/juryAssignments"
	"main/internal/storage/orm"

	"golang.org/x/sync/errgroup"
)

type juryAssignmentsRepositoryInterface interface {
	GetJuryAssignmentsByFilter(orm.ORM, juryAssignments.JuryAssignments) (juryAssignments.JuryAssignments, error)
	GetAllJuryAssignments(orm.ORM, ...interface{}) ([]juryAssignments.JuryAssignments, error)
	GetAllJuryAssignmentsByFilter(orm.ORM, juryAssignments.JuryAssignments) ([]juryAssignments.JuryAssignments, error)
	GetPartOfAllJuryAssignmentsByFilter(
		orm.ORM, []string, juryAssignments.JuryAssignments) ([]juryAssignments.JuryAssignments, error)
	CreateJuryAssignments(
		orm.ORM, juryAssignments.JuryAssignments) (uint, error)
	UpdateJuryAssignments(
		orm.ORM, juryAssignments.JuryAssignments) (uint, error)
	DeleteJuryAssignments(orm.ORM, uint) error
}

type JuryAssignmentsService struct {
	orm        orm.ORM
	repository juryAssignmentsRepositoryInterface
}

func NewJuryAssignmentsService(orm orm.ORM, jr juryAssignmentsRepositoryInterface) *JuryAssignmentsService {
	return &JuryAssignmentsService{orm: orm, repository: jr}
}

func (s *JuryAssignmentsService) GetAllJuryAssignments() ([]juryAssignmentsDto.JuryAssignmentsDTO, error) {
	const op = "services.juryAssignmentsService.GetAllJuryAssignments"
	results, err := s.repository.GetAllJuryAssignments(s.orm)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return dtoConverter.ConvertManyJuryAssignmentsToDTO(results), nil
}

func (s *JuryAssignmentsService) GetAllJuryAssignmentsByFilter(
	filter juryAssignmentsDto.JuryAssignmentsDTO) ([]juryAssignmentsDto.JuryAssignmentsDTO, error) {
	const op = "services.juryAssignmentsService.GetAllEventsByFilter"
	modelFilter := dtoConverter.ConvertDTOtoJuryAssignments(filter)
	results, err := s.repository.GetAllJuryAssignmentsByFilter(s.orm, modelFilter)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return dtoConverter.ConvertManyJuryAssignmentsToDTO(results), nil
}

func (s *JuryAssignmentsService) GetPartOfAllJuryAssignmentsByFilter(
	fields []string, filter juryAssignmentsDto.JuryAssignmentsDTO) ([]juryAssignmentsDto.JuryAssignmentsDTO, error) {
	const op = "services.juryAssignmentsService.GetPartOfAllJuryAssignmentsByFilter"
	modelFilter := dtoConverter.ConvertDTOtoJuryAssignments(filter)
	result, err := s.repository.GetPartOfAllJuryAssignmentsByFilter(s.orm, fields, modelFilter)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return dtoConverter.ConvertManyJuryAssignmentsToDTO(result), nil
}

func (s *JuryAssignmentsService) GetJuryAssignmentsByFilter(filter juryAssignmentsDto.JuryAssignmentsDTO) (juryAssignmentsDto.JuryAssignmentsDTO, error) {
	const op = "services.juryAssignmentsService.GetJuryAssignmentsByFilter"

	modelFilter := dtoConverter.ConvertDTOtoJuryAssignments(filter)
	result, err := s.repository.GetJuryAssignmentsByFilter(s.orm, modelFilter)
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

func (s *JuryAssignmentsService) CreateManyAssignmentsByOneJury(dto juryAssignmentsDto.OneJuryManyAssignments) ([]uint, error) {
	const op = "services.juryAssignmentsService.CreateManyAssignmentsByOneJury"

	errGroup := errgroup.Group{}

	ids := make(chan uint, len(dto.EventsID))

	for _, eventID := range dto.EventsID {
		tempDto := juryAssignmentsDto.JuryAssignmentsDTO{JuryID: dto.JuryID, EventID: eventID}
		errGroup.Go(func() error {
			id, err := s.CreateJuryAssignments(tempDto)
			if err != nil {
				return err
			}
			ids <- id
			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		close(ids)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	close(ids)

	sliceIds := []uint{}
	for id := range ids {
		sliceIds = append(sliceIds, id)
	}

	return sliceIds, nil
}

// func (s *JuryAssignmentsService) CreateJuryAssignmentsWithTransactionSupport(db *gorm.DB, dto juryAssignmentsDto.JuryAssignmentsDTO) (uint, error) {
// 	const op = "services.juryAssignmentsService.CreateJuryAssignments"
// 	if !isExistingUserID(dto.JuryID) {
// 		return 0, fmt.Errorf("%s: Jury is not exist", op)
// 	}
// 	if !isExistingEventID(dto.EventID) {
// 		return 0, fmt.Errorf("%s: Event is not exist", op)
// 	}

// 	model := dtoConverter.ConvertDTOtoJuryAssignments(dto)
// 	id, err := s.repository.CreateJuryAssignments(db, model)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}
// 	return id, nil
// }

func (s *JuryAssignmentsService) CreateJuryAssignments(dto juryAssignmentsDto.JuryAssignmentsDTO) (uint, error) {
	const op = "services.juryAssignmentsService.CreateJuryAssignments"
	if !isExistingUserID(dto.JuryID) {
		return 0, fmt.Errorf("%s: Jury is not exist", op)
	}
	if !isExistingEventID(dto.EventID) {
		return 0, fmt.Errorf("%s: Event is not exist", op)
	}

	model := dtoConverter.ConvertDTOtoJuryAssignments(dto)
	id, err := s.repository.CreateJuryAssignments(s.orm, model)
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
	id, err := s.repository.UpdateJuryAssignments(s.orm, model)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *JuryAssignmentsService) DeleteJuryAssignments(id uint) error {
	const op = "services.juryAssignmentsService.CreateJuryAssignments"
	err := s.repository.DeleteJuryAssignments(s.orm, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
