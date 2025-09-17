package jury_assignments_service

import (
	"context"
	"fmt"
	"log/slog"
	"main/internal/dto/juryAssignmentsDto"
	"main/internal/lib/converter/dtoConverter"
	"main/internal/lib/errs"
	"main/internal/lib/liblogger"
	"main/internal/lib/supportRequest"
	"main/internal/models/jury_assignments"
	"main/internal/storage/orm"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type juryAssignmentsRepositoryInterface interface {
	GetJuryAssignmentsByFilter(context.Context, orm.ORM, jury_assignments.JuryAssignments) (jury_assignments.JuryAssignments, error)
	GetAllJuryAssignments(context.Context, orm.ORM) ([]jury_assignments.JuryAssignments, error)
	GetAllJuryAssignmentsByFilter(context.Context, orm.ORM, jury_assignments.JuryAssignments) ([]jury_assignments.JuryAssignments, error)
	CreateJuryAssignments(context.Context,
		orm.ORM, jury_assignments.JuryAssignments) (uuid.UUID, error)

	UpdateJuryAssignments(context.Context,
		orm.ORM, jury_assignments.JuryAssignments) error
	DeleteJuryAssignments(context.Context, orm.ORM, uuid.UUID) error
	DeleteByFields(context.Context, orm.ORM, jury_assignments.JuryAssignments) error
}

type JuryAssignmentsService struct {
	orm        orm.ORM
	supportReq *supportRequest.SupportRequest
	repository juryAssignmentsRepositoryInterface
	log        *slog.Logger
}

func NewJuryAssignmentsService(log *slog.Logger, orm orm.ORM, supportReq *supportRequest.SupportRequest, jr juryAssignmentsRepositoryInterface) *JuryAssignmentsService {
	jlog := log.With("owner", "JuryAssignmentsService")

	return &JuryAssignmentsService{
		supportReq: supportReq,
		orm:        orm,
		repository: jr,
		log:        jlog,
	}
}

func (s *JuryAssignmentsService) GetAllJuryAssignments(ctx context.Context) ([]juryAssignmentsDto.JuryAssignmentsResponseDTO, error) {
	const op = "services.juryAssignmentsService.GetAllJuryAssignments"

	log := s.log.With(
		slog.String("op", op),
	)

	results, err := s.repository.GetAllJuryAssignments(ctx, s.orm)
	if err != nil {
		log.Error("failed to get all jury assignments", liblogger.Err(err))
		return nil, errs.ErrInternalError.Wrap("failed get all jury-assignments")
	}

	return dtoConverter.ConvertManyJuryAssignmentsToDTO(results), nil
}

func (s *JuryAssignmentsService) GetJuryAssignmentsByID(ctx context.Context, id string) (juryAssignmentsDto.JuryAssignmentsResponseDTO, error) {
	const op = "services.juryAssignmentsService.GetJuryAssignmentsByID"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id to uuid", slog.String("id", id), liblogger.Err(err))
		return juryAssignmentsDto.JuryAssignmentsResponseDTO{}, errs.ErrBadRequest.Wrap("failed parse id")
	}

	filter := jury_assignments.JuryAssignments{ID: uid}

	juryResult, err := s.repository.GetJuryAssignmentsByFilter(ctx, s.orm, filter)
	if s.orm.IsNotFound(err) {
		log.Warn("jury-assignments not found", liblogger.Err(err))
		return juryAssignmentsDto.JuryAssignmentsResponseDTO{}, errs.ErrNotFound.Wrap("jury-assignments not found")
	}

	if err != nil {
		log.Error("failed get jury by id", liblogger.Err(err))
		return juryAssignmentsDto.JuryAssignmentsResponseDTO{}, errs.ErrInternalError.Wrap("failed get jury-assignments by id")
	}

	return dtoConverter.ConvertJuryAssignmentsToDTO(juryResult), nil
}

func (s *JuryAssignmentsService) GetAllByEventId(ctx context.Context, eventId string) ([]juryAssignmentsDto.JuryAssignmentsResponseDTO, error) {
	const op = "services.juryAssignmentsService.GetAllByEventId"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(eventId)
	if err != nil {
		log.Error("failed parse id to uuid", slog.String("id", eventId), liblogger.Err(err))
		return nil, errs.ErrBadRequest.Wrap("failed parse id")
	}

	model := jury_assignments.JuryAssignments{EventID: uid}
	result, err := s.repository.GetAllJuryAssignmentsByFilter(ctx, s.orm, model)
	if err != nil {
		log.Error("failed get all jury by event id", liblogger.Err(err))
		return nil, errs.ErrInternalError.Wrap("failed get jury-assignments by event-id")
	}

	return dtoConverter.ConvertManyJuryAssignmentsToDTO(result), nil
}

func (s *JuryAssignmentsService) GetAllByJuryId(ctx context.Context, juryId string) ([]juryAssignmentsDto.JuryAssignmentsResponseDTO, error) {
	const op = "services.juryAssignmentsService.GetAllByJuryId"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(juryId)
	if err != nil {
		log.Error("failed parse id to uuid", slog.String("id", juryId), liblogger.Err(err))
		return nil, errs.ErrBadRequest.Wrap("failed parse id")
	}

	model := jury_assignments.JuryAssignments{UserID: uid}
	result, err := s.repository.GetAllJuryAssignmentsByFilter(ctx, s.orm, model)
	if err != nil {
		log.Error("failed get all by jury id", liblogger.Err(err))
		return nil, errs.ErrInternalError.Wrap("failed get jury-assignments by user-id")
	}

	return dtoConverter.ConvertManyJuryAssignmentsToDTO(result), nil
}

func (s *JuryAssignmentsService) GetAllJuryAssignmentsByFilter(ctx context.Context,
	filter juryAssignmentsDto.JuryAssignmentsResponseDTO) ([]juryAssignmentsDto.JuryAssignmentsResponseDTO, error) {
	const op = "services.juryAssignmentsService.GetAllEventsByFilter"

	log := s.log.With(
		slog.String("op", op),
	)

	modelFilter := dtoConverter.ConvertDTOtoJuryAssignments(filter)
	results, err := s.repository.GetAllJuryAssignmentsByFilter(ctx, s.orm, modelFilter)
	if err != nil {
		log.Error("failed get all jury-assignmetns by jury-id", liblogger.Err(err))
		return nil, errs.ErrInternalError.Wrap("failed get all jury-assigments by filter")
	}
	return dtoConverter.ConvertManyJuryAssignmentsToDTO(results), nil
}

func (s *JuryAssignmentsService) CreateMany(ctx context.Context, dto juryAssignmentsDto.CreateOneAssigmentsManyJury) ([]uuid.UUID, error) {
	const op = "services.juryAssignmentsService.CreateMany"

	errGroup := errgroup.Group{}

	ids := make(chan uuid.UUID, len(dto.JuryIDs))

	for _, juryId := range dto.JuryIDs {
		tempDto := juryAssignmentsDto.CreateJuryAssignmentsDTO{UserID: juryId, EventID: dto.EventID}
		errGroup.Go(func() error {
			id, err := s.Create(ctx, tempDto)
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

// func isExistingID(id uuid.UUID, service string) (bool, error) {
// 	if id == uuid.Nil {
// 		return false, fmt.Errorf("id is nil: %v", id)
// 	}
// 	response, err := supportRequest.SupportRequest(id, service)
// 	if err != nil {
// 		return false, err
// 	}
// 	if response.Data == nil {
// 		return false, fmt.Errorf("data is nil")
// 	}
// 	return true, nil
// }

// func (s *JuryAssignmentsService) CreateManyAssignmentsByOneJury(ctx context.Context, dto juryAssignmentsDto.OneJuryManyAssignments) ([]uuid.UUID, error) {
// 	const op = "services.juryAssignmentsService.CreateManyAssignmentsByOneJury"

// 	errGroup := errgroup.Group{}

// 	ids := make(chan uuid.UUID, len(dto.EventIDs))

// 	for _, eventID := range dto.EventIDs {
// 		tempDto := juryAssignmentsDto.CreateJuryAssignmentsDTO{JuryID: dto.JuryID, EventID: eventID}
// 		errGroup.Go(func() error {
// 			id, err := s.Create(ctx, tempDto)
// 			if err != nil {
// 				return err
// 			}
// 			ids <- id
// 			return nil
// 		})
// 	}

// 	if err := errGroup.Wait(); err != nil {
// 		close(ids)
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	close(ids)

// 	sliceIds := []uuid.UUID{}
// 	for id := range ids {
// 		sliceIds = append(sliceIds, id)
// 	}

// 	return sliceIds, nil
// }

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

func (s *JuryAssignmentsService) Create(ctx context.Context, dto juryAssignmentsDto.CreateJuryAssignmentsDTO) (uuid.UUID, error) {
	const op = "services.juryAssignmentsService.Create"

	log := s.log.With(
		slog.String("op", op),
	)

	model, err := dtoConverter.FromCreateDTOtoModel(dto)
	if err != nil {
		log.Error("failed convert dto to model", liblogger.Err(err))
		return uuid.Nil, errs.ErrInternalError.Wrap("failed convert dto to model")
	}

	// ok, err := s.supportReq.PrepareRequest(model.JuryID, supportRequest.EventService)
	// if err != nil {
	// 	log.Error("failed existing jury by id", liblogger.Err(err))
	// 	return uuid.Nil, fmt.Errorf("%s: failed exist jury by", errMsg)
	// }
	// if !ok {
	// 	log.Error("failed find jury by id")
	// 	return uuid.Nil, fmt.Errorf("%s: Jury is not exist", errMsg)
	// }

	// ok, err = s.supportReq.PrepareRequest(model.EventID, supportRequest.UserService)
	// if err != nil {
	// 	log.Error("failed existing event by id", liblogger.Err(err))
	// 	return uuid.Nil, fmt.Errorf("%s: failed exist event by", errMsg)
	// }
	// if ok {
	// 	log.Error("failed find jury by id")
	// 	return uuid.Nil, fmt.Errorf("%s: Event is not exist", errMsg)
	// }

	id, err := s.repository.CreateJuryAssignments(ctx, s.orm, model)
	if err != nil {
		log.Error("failed create jury-assignmetns")
		return uuid.Nil, errs.ErrInternalError.Wrap("failed craete jury-assignmetns")
	}
	return id, nil
}

func (s *JuryAssignmentsService) Update(ctx context.Context, id string, dto juryAssignmentsDto.UpdateJuryAssignmentsDTO) error {
	const op = "services.juryAssignmentsService.Update"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id to uuid", slog.String("id", id), liblogger.Err(err))
		return errs.ErrBadRequest.Wrap("failed parse id")
	}

	model, err := dtoConverter.FromUpdateDTOtoModel(dto, uid)
	if err != nil {
		log.Error("failed convert dto to model", liblogger.Err(err))
		return errs.ErrInternalError.Wrap("failed dto to model")
	}

	// if dto.JuryID != nil {
	// 	ok, err := s.supportReq.PrepareRequest(model.JuryID, supportRequest.UserService)
	// 	if err != nil {
	// 		log.Error("failed existing jury user by id", liblogger.Err(err))
	// 		return fmt.Errorf("%s: Failed existing jury by id", errMsg)
	// 	}
	// 	if !ok {
	// 		log.Error("failed find user by id")
	// 		return fmt.Errorf("%s: Jury is not exist", errMsg)
	// 	}
	// }

	// if dto.EventID != nil {
	// 	ok, err := s.supportReq.PrepareRequest(model.EventID, supportRequest.EventService)
	// 	if err != nil {
	// 		log.Error("failed existing event by id", liblogger.Err(err))
	// 		return fmt.Errorf("%s: failed exist event by", errMsg)
	// 	}
	// 	if ok {
	// 		log.Error("failed find jury by id")
	// 		return fmt.Errorf("%s: Event is not exist", errMsg)
	// 	}
	// }

	err = s.repository.UpdateJuryAssignments(ctx, s.orm, model)
	if err != nil {
		log.Error("failed update jury-assigment", liblogger.Err(err))
		return errs.ErrInternalError.Wrap("failed update jury-assigment")
	}
	return nil
}

func (s *JuryAssignmentsService) Delete(ctx context.Context, id string) error {
	const op = "services.juryAssignmentsService.Delete"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id to uuid", slog.String("id", id), liblogger.Err(err))
		return errs.ErrBadRequest.Wrap("failed parse id")
	}

	err = s.repository.DeleteJuryAssignments(ctx, s.orm, uid)
	if err != nil {
		log.Error("failed delete jury-assinment", liblogger.Err(err))
		return errs.ErrInternalError.Wrap("failed delete jury-assinment")
	}
	return nil
}

func (s *JuryAssignmentsService) DeleteByFields(ctx context.Context, juryAssignment juryAssignmentsDto.JuryAssignmentsResponseDTO) error {
	const op = "services.juryAssignmentsService.DeleteByFields"

	log := s.log.With(
		slog.String("op", op),
	)

	juryAssignmentModel := dtoConverter.ConvertDTOtoJuryAssignments(juryAssignment)
	err := s.repository.DeleteByFields(ctx, s.orm, juryAssignmentModel)
	if err != nil {
		log.Error("failed delete jury-assingment", liblogger.Err(err))
		return errs.ErrInternalError.Wrap("failed delete jury-assingment by fields")
	}

	return nil
}
