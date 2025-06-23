package refresh_repository

import (
	"context"
	"main/internal/storage/orm"

	"github.com/google/uuid"
)

type RefreshRepository struct{}

func (r *RefreshRepository) Create(ctx context.Context, orm orm.ORM) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}
