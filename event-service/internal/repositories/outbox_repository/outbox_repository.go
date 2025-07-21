package outbox_repository

import (
	"context"
	"fmt"
	"main/internal/models/outbox"
	"main/internal/storage/orm"

	"github.com/google/uuid"
)

type OutboxRepository struct{}

func (r *OutboxRepository) GetAllNotProcessed(ctx context.Context, orm orm.ORM) ([]outbox.Outbox, error) {
	const op = "repository.OutboxRepository.GetAllNotProcessed"
	outboxResult := []outbox.Outbox{}

	err := orm.Find(ctx, outbox.Outbox{}, nil, nil, nil, nil, &outboxResult, outbox.Outbox{Processed: false})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return outboxResult, nil
}

func (r *OutboxRepository) Create(ctx context.Context, orm orm.ORM, outbox outbox.Outbox) (uuid.UUID, error) {
	const op = "repository.OutboxRepository.Create"
	err := orm.Create(ctx, &outbox)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return outbox.ID, nil
}

func (r *OutboxRepository) Update(ctx context.Context, orm orm.ORM, outbox outbox.Outbox) error {
	const op = "repository.OutboxRepository.Update"
	err := orm.Updates(ctx, outbox)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *OutboxRepository) Delete(ctx context.Context, orm orm.ORM, id uuid.UUID) error {
	const op = "repository.OutboxRepository.Delete"
	err := orm.Delete(ctx, outbox.Outbox{ID: id})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
