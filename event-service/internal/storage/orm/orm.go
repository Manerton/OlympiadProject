package orm

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type ORM interface {
	// Find by parametrs
	//
	// model - your db model
	//
	// fields - []string filelds model for select
	//
	// offset, limit - just ofset and limit)
	//
	// order - ORDER BY
	Find(ctx context.Context, model interface{}, fields *[]string, offset, limit *int, order *string, dest interface{}, conds ...interface{}) error
	First(ctx context.Context, model interface{}, fields *[]string, dest interface{}, conds ...interface{}) error

	Count(ctx context.Context, model interface{}, count *int64, query interface{}, args ...interface{}) error
	Create(ctx context.Context, dest interface{}) error
	Updates(ctx context.Context, dest interface{}) error
	Delete(ctx context.Context, dest interface{}, conds ...interface{}) error
	TransactionBegin() (ORM, error)
	TransactionCommit() error
	TransactionRollback() error
}

type Gorm struct {
	DB *gorm.DB
}

func NewGormORM(db *gorm.DB) ORM {
	return &Gorm{DB: db}
}

func (g *Gorm) Count(ctx context.Context, model interface{}, count *int64, query interface{}, args ...interface{}) error {
	const op = "storage.orm.Count"
	err := g.DB.WithContext(ctx).Model(model).Where(query, args).Count(count).Error
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (g *Gorm) Create(ctx context.Context, dest interface{}) error {
	const op = "storage.orm.Create"
	if err := g.DB.WithContext(ctx).Create(dest).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (g *Gorm) Updates(ctx context.Context, dest interface{}) error {
	const op = "storage.orm.Update"
	if err := g.DB.WithContext(ctx).Updates(dest).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (g *Gorm) Delete(ctx context.Context, dest interface{}, conds ...interface{}) error {
	const op = "storage.orm.Delete"
	if err := g.DB.WithContext(ctx).Delete(dest, conds...).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (g *Gorm) First(ctx context.Context, model interface{}, fields *[]string, dest interface{}, conds ...interface{}) error {
	const op = "storage.orm.First"
	query := g.DB.WithContext(ctx).Model(model)
	if fields != nil {
		query.Select(*fields)
	}

	if err := query.First(dest, conds...).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Find by parametrs
func (g *Gorm) Find(ctx context.Context, model interface{}, fields *[]string, offset, limit *int, order *string, dest interface{}, conds ...interface{}) error {
	const op = "storage.orm.Find"
	query := g.DB.WithContext(ctx).Model(model)
	if order != nil {
		query.Order(*order)
	}
	if fields != nil {
		query.Select(*fields)
	}
	if offset != nil {
		query = query.Offset(*offset)
	}
	if limit != nil {
		query = query.Limit(*limit)
	}

	if err := query.Find(dest, conds...).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (g *Gorm) TransactionBegin() (ORM, error) {
	tx := g.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &Gorm{DB: tx}, nil
}

func (g *Gorm) TransactionCommit() error {
	if err := g.DB.Commit().Error; err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}
	return nil
}

func (g *Gorm) TransactionRollback() error {
	if err := g.DB.Rollback().Error; err != nil {
		return fmt.Errorf("transaction rollback failed: %w", err)
	}
	return nil
}
