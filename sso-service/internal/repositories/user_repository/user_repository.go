package user_repository

import (
	"context"
	"fmt"
	"main/internal/models/user"
	"main/internal/storage/orm"

	"github.com/google/uuid"
)

type UserRepository struct {
	model user.User
}

func (r *UserRepository) GetByEmail(ctx context.Context, orm orm.ORM, email string) (user.User, error) {
	const op = "repositories.UserRepository.GetByEmail"

	userResult := user.User{}
	err := orm.First(ctx, r.model, nil, &userResult, user.User{Email: email})
	if err != nil {
		return user.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return userResult, nil
}

func (r *UserRepository) GetById(ctx context.Context, orm orm.ORM, id uuid.UUID) (user.User, error) {
	const op = "repositories.UserRepository.GetById"
	return user.User{}, nil
}

func (r *UserRepository) GetByListId(ctx context.Context, orm orm.ORM, ids []uuid.UUID) ([]user.User, error) {
	const op = "repositories.UserRepository.GetByListId"
	return []user.User{}, nil
}

func (r *UserRepository) Create(ctx context.Context, orm orm.ORM, user user.User) (uuid.UUID, error) {
	const op = "repositories.UserRepository.Create"

	err := orm.Create(ctx, &user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return user.ID, nil
}
