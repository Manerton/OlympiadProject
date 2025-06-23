package school_repository

import (
	"context"
	"fmt"
	"main/internal/models/school"
	"main/internal/storage/orm"

	"github.com/google/uuid"
)

type SchoolRepository struct{}

func (r *SchoolRepository) GetById(ctx context.Context, orm orm.ORM, id uuid.UUID) (school.School, error) {
	const op = "repositories.school_repository.GetById"

	schoolResult := school.School{}
	err := orm.First(ctx, school.School{}, nil, &schoolResult, school.School{ID: id})
	if err != nil {
		return school.School{}, fmt.Errorf("%s: %w", op, err)
	}
	return schoolResult, nil
}

func (r *SchoolRepository) GetAll(ctx context.Context, orm orm.ORM, offset, limit *int) ([]school.School, error) {
	const op = "repositories.school_repository.GetAll"

	schoolsResult := []school.School{}
	err := orm.Find(ctx, school.School{}, nil, offset, limit, nil, &schoolsResult, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return schoolsResult, nil
}

func (r *SchoolRepository) Create(ctx context.Context, orm orm.ORM, school school.School) (uuid.UUID, error) {
	const op = "repositories.school_repository.Create"

	err := orm.Create(ctx, &school)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return school.ID, nil
}

func (r *SchoolRepository) Update(ctx context.Context, orm orm.ORM, school school.School) error {
	const op = "repositories.school_repository.Update"

	err := orm.Updates(ctx, &school)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *SchoolRepository) Delete(ctx context.Context, orm orm.ORM, id uuid.UUID) error {
	const op = "repositories.school_repository.Delete"

	err := orm.Delete(ctx, &school.School{ID: id})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
