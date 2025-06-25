package refresh_reptository

import (
	"context"
	"fmt"
	"main/internal/models/refresh_token"
	"main/internal/storage/orm"

	"github.com/google/uuid"
)

type RefreshRepository struct{}

func (r *RefreshRepository) GetById(ctx context.Context, orm orm.ORM, id uuid.UUID) (refresh_token.RefreshToken, error) {
	const op = "repository.refresh_reposiotry.GetById"

	resultToken := refresh_token.RefreshToken{}
	err := orm.First(ctx, refresh_token.RefreshToken{}, nil, &resultToken)
	if err != nil {
		return refresh_token.RefreshToken{}, fmt.Errorf("%s: %w", op, err)
	}

	return resultToken, nil
}

func (r *RefreshRepository) GetByUserId(ctx context.Context, orm orm.ORM, userId uuid.UUID) ([]refresh_token.RefreshToken, error) {
	const op = "repository.refresh_reposiotry.GetByUserId"

	resultsToken := []refresh_token.RefreshToken{}
	err := orm.Find(ctx, refresh_token.RefreshToken{}, nil, nil, nil, nil, &resultsToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resultsToken, nil
}

func (r *RefreshRepository) Create(ctx context.Context, orm orm.ORM, token refresh_token.RefreshToken) (uuid.UUID, error) {
	const op = "repository.refresh_reposiotry.Create"

	err := orm.Create(ctx, &token)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return token.ID, nil
}

func (r *RefreshRepository) Update(ctx context.Context, orm orm.ORM, token refresh_token.RefreshToken) error {
	const op = "repository.refresh_reposiotry.Update"

	err := orm.Updates(ctx, &token)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *RefreshRepository) Delete(ctx context.Context, orm orm.ORM, id uuid.UUID) error {
	const op = "repository.refresh_reposiotry.Update"

	modelDel := refresh_token.RefreshToken{ID: id}
	err := orm.Delete(ctx, modelDel)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
