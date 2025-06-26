package jury_assignments_service

import (
	"context"
	"fmt"
	"log/slog"
	"main/internal/dto/juryAssignmentsDto"
	"main/internal/lib/converter/dtoConverter"
	"main/internal/lib/liblogger"
	"main/internal/lib/supportRequest"
	"main/internal/models/jury_assignments"
	"main/internal/storage/orm"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type juryAssignmentsRepositoryInterface interface {
	GetJuryAssignmentsByFilter(context.Context, orm.ORM, jury_assignments.JuryAssignments) (jury_assignments.JuryAssignments, error)
	GetAllJuryAssignments(context.Context, orm.ORM, ...interface{}) ([]jury_assignments.JuryAssignments, error)
	GetAllJuryAssignmentsByFilter(context.Context, orm.ORM, jury_assignments.JuryAssignments) ([]jury_assignments.JuryAssignments, error)
	GetPartOfAllJuryAssignmentsByFilter(context.Context,
		orm.ORM, []string, jury_assignments.JuryAssignments) ([]jury_assignments.JuryAssignments, error)
	CreateJuryAssignments(context.Context,
		orm.ORM, jury_assignments.JuryAssignments) (uuid.UUID, error)
	UpdateJuryAssignments(context.Context,
		orm.ORM, jury_assignments.JuryAssignments) error
	DeleteJuryAssignments(context.Context, orm.ORM, uuid.UUID) error
}

type JuryAssignmentsService struct {
	orm        orm.ORM
	repository juryAssignmentsRepositoryInterface
	log        *slog.Logger
}

func NewJuryAssignmentsService(log *slog.Logger, orm orm.ORM, jr juryAssignmentsRepositoryInterface) *JuryAssignmentsService {
	return &JuryAssignmentsService{
		orm:        orm,
		repository: jr,
		log:        log,
	}
}

func (s *JuryAssignmentsService) GetAllJuryAssignments(ctx context.Context) ([]juryAssignmentsDto.JuryAssignmentsDTO, error) {
	const op = "services.juryAssignmentsService.GetAllJuryAssignments"

	log := s.log.With(
		slog.String("op", op),
	)

	results, err := s.repository.GetAllJuryAssignments(ctx, s.orm)
	if err != nil {
		log.Error("failed to get all jury assignments", liblogger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return dtoConverter.ConvertManyJuryAssignmentsToDTO(results), nil
}

func (s *JuryAssignmentsService) GetAllJuryAssignmentsByFilter(ctx context.Context,
	filter juryAssignmentsDto.JuryAssignmentsDTO) ([]juryAssignmentsDto.JuryAssignmentsDTO, error) {
	const op = "services.juryAssignmentsService.GetAllEventsByFilter"
	modelFilter := dtoConverter.ConvertDTOtoJuryAssignments(filter)
	results, err := s.repository.GetAllJuryAssignmentsByFilter(ctx, s.orm, modelFilter)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return dtoConverter.ConvertManyJuryAssignmentsToDTO(results), nil
}

func (s *JuryAssignmentsService) GetPartOfAllJuryAssignmentsByFilter(ctx context.Context,
	fields []string, filter juryAssignmentsDto.JuryAssignmentsDTO) ([]juryAssignmentsDto.JuryAssignmentsDTO, error) {
	const op = "services.juryAssignmentsService.GetPartOfAllJuryAssignmentsByFilter"
	modelFilter := dtoConverter.ConvertDTOtoJuryAssignments(filter)
	result, err := s.repository.GetPartOfAllJuryAssignmentsByFilter(ctx, s.orm, fields, modelFilter)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return dtoConverter.ConvertManyJuryAssignmentsToDTO(result), nil
}

func (s *JuryAssignmentsService) GetJuryAssignmentsByFilter(ctx context.Context, filter juryAssignmentsDto.JuryAssignmentsDTO) (juryAssignmentsDto.JuryAssignmentsDTO, error) {
	const op = "services.juryAssignmentsService.GetJuryAssignmentsByFilter"

	modelFilter := dtoConverter.ConvertDTOtoJuryAssignments(filter)
	result, err := s.repository.GetJuryAssignmentsByFilter(ctx, s.orm, modelFilter)
	if err != nil {
		return juryAssignmentsDto.JuryAssignmentsDTO{}, fmt.Errorf("%s: %w", op, err)
	}
	return dtoConverter.ConvertJuryAssignmentsToDTO(result), nil
}

func isExistingUserID(id uuid.UUID) (bool, error) {
	return true, nil
	if id == uuid.Nil {
		return false, fmt.Errorf("id is nil: %v", id)
	}
	path := ""
	response, err := supportRequest.SupportRequest("GET", path)
	if err != nil {
		return false, err
	}
	if response.Data == nil {
		return false, fmt.Errorf("data is nil")
	}
	return true, nil
}

func isExistingEventID(id uuid.UUID) (bool, error) {
	return true, nil

	if id == uuid.Nil {
		return false, fmt.Errorf("id is nil: %v", id)
	}
	path := ""
	response, err := supportRequest.SupportRequest("GET", path)
	if err != nil {
		return false, err
	}
	if response.Data == nil {
		return false, fmt.Errorf("data is nil")
	}
	return true, nil
}

func (s *JuryAssignmentsService) CreateManyAssignmentsByOneJury(ctx context.Context, dto juryAssignmentsDto.OneJuryManyAssignments) ([]uuid.UUID, error) {
	const op = "services.juryAssignmentsService.CreateManyAssignmentsByOneJury"

	errGroup := errgroup.Group{}

	ids := make(chan uuid.UUID, len(dto.EventsID))

	for _, eventID := range dto.EventsID {
		tempDto := juryAssignmentsDto.JuryAssignmentsDTO{JuryID: dto.JuryID, EventID: eventID}
		errGroup.Go(func() error {
			id, err := s.CreateJuryAssignments(ctx, tempDto)
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

	sliceIds := []uuid.UUID{}
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

func (s *JuryAssignmentsService) CreateJuryAssignments(ctx context.Context, dto juryAssignmentsDto.JuryAssignmentsDTO) (uuid.UUID, error) {
	const op = "services.juryAssignmentsService.CreateJuryAssignments"

	log := s.log.With(
		slog.String("op", op),
	)

	ok, err := isExistingUserID(dto.JuryID)
	if err != nil {

	}
	if !ok {
		log.Error("failed find user by id")
		return uuid.Nil, fmt.Errorf("Jury is not exist")
	}
	ok, err = isExistingEventID(dto.EventID)
	if err != nil {

	}
	if ok {
		log.Error("failed find jury by id")
		return uuid.Nil, fmt.Errorf("Event is not exist")
	}

	model := dtoConverter.ConvertDTOtoJuryAssignments(dto)
	id, err := s.repository.CreateJuryAssignments(ctx, s.orm, model)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *JuryAssignmentsService) UpdateJuryAssignments(ctx context.Context, id string, dto juryAssignmentsDto.JuryAssignmentsDTO) error {
	const op = "services.juryAssignmentsService.CreateJuryAssignments"

	log := s.log.With(
		slog.String("op", op),
	)

	// TODO!! ГДЕ проверять все uuid, validator? mapper?

	ok, err := isExistingUserID(dto.JuryID)
	if err != nil {

	}
	if !ok {
		log.Error("failed find user by id")
		return fmt.Errorf("Jury is not exist")
	}
	ok, err = isExistingEventID(dto.EventID)
	if err != nil {

	}
	if ok {
		log.Error("failed find jury by id")
		return fmt.Errorf("Event is not exist")
	}

	model := dtoConverter.ConvertDTOtoJuryAssignments(dto)
	err = s.repository.UpdateJuryAssignments(ctx, s.orm, model)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *JuryAssignmentsService) DeleteJuryAssignments(ctx context.Context, id string) error {
	const op = "services.juryAssignmentsService.CreateJuryAssignments"
	const errMsg = "failed delete JuryAssignments"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id to uuid", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	err = s.repository.DeleteJuryAssignments(ctx, s.orm, uid)
	if err != nil {
		log.Error("failed delete jury assinments", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}
	return nil
}
