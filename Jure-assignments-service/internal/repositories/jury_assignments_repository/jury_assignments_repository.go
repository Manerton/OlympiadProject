package jury_assignments_repository

import (
	"context"
	"fmt"
	"main/internal/models/jury_assignments"
	"main/internal/storage/orm"

	"github.com/google/uuid"
)

type JuryAssignmentsRepository struct{}

func NewJuryAssignmentsRepository() *JuryAssignmentsRepository {
	return &JuryAssignmentsRepository{}
}

func (r *JuryAssignmentsRepository) GetJuryAssignmentsByFilter(ctx context.Context,
	myOrm orm.ORM, filter jury_assignments.JuryAssignments) (jury_assignments.JuryAssignments, error) {
	const op = "repositories.juryAssignmentsRepository.GetJuryAssignmentsByID"

	result := jury_assignments.JuryAssignments{}
	if err := myOrm.First(ctx, jury_assignments.JuryAssignments{}, nil, &result, filter); err != nil {
		return jury_assignments.JuryAssignments{}, fmt.Errorf("%s: %w", op, err)
	}
	return result, nil
}

func (r *JuryAssignmentsRepository) GetAllJuryAssignments(ctx context.Context, myOrm orm.ORM) ([]jury_assignments.JuryAssignments, error) {
	const op = "repositories.juryAssignmentsRepository.GetAllJuryAssignments"

	juryAssignmentsRes := []jury_assignments.JuryAssignments{}
	if err := myOrm.Find(ctx, jury_assignments.JuryAssignments{}, nil, nil, nil, nil, &juryAssignmentsRes); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return juryAssignmentsRes, nil
}

func (r *JuryAssignmentsRepository) GetAllJuryAssignmentsByFilter(ctx context.Context,
	myOrm orm.ORM, filter jury_assignments.JuryAssignments) ([]jury_assignments.JuryAssignments, error) {
	const op = "repositories.juryAssignmentsRepository.GetAllJuryAssignmentsByFilter"

	juryAssignmentsRes := []jury_assignments.JuryAssignments{}
	if err := myOrm.Find(
		ctx,
		jury_assignments.JuryAssignments{}, nil, nil, nil, nil,
		&juryAssignmentsRes,
		filter); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return juryAssignmentsRes, nil
}

func (r *JuryAssignmentsRepository) CreateJuryAssignments(ctx context.Context,
	myOrm orm.ORM, juryAssignments jury_assignments.JuryAssignments) (uuid.UUID, error) {
	const op = "repositories.juryAssignmentsRepository.CreateJuryAssignments"

	if err := myOrm.Create(ctx, &juryAssignments); err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return juryAssignments.ID, nil
}

func (r *JuryAssignmentsRepository) UpdateJuryAssignments(ctx context.Context,
	myOrm orm.ORM, juryAssignments jury_assignments.JuryAssignments) error {
	const op = "repositories.juryAssignmentsRepository.UpdateJuryAssignments"

	if err := myOrm.Updates(ctx, nil, &juryAssignments); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *JuryAssignmentsRepository) DeleteJuryAssignments(ctx context.Context, myOrm orm.ORM, id uuid.UUID) error {
	const op = "repositories.juryAssignmentsRepository.DeleteJuryAssignments"
	if err := myOrm.Delete(ctx, &jury_assignments.JuryAssignments{ID: id}); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
