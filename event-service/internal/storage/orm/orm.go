package orm

import (
	"fmt"

	"gorm.io/gorm"
)

type ORM interface {
	Find(model interface{}, fields *[]string, offset, limit *int, dest interface{}, conds ...interface{}) error
	First(model interface{}, fields *[]string, dest interface{}, conds ...interface{}) error

	Count(model interface{}, count *int64, query interface{}, args ...interface{}) error
	Create(dest interface{}) error
	Updates(dest interface{}) error
	Delete(dest interface{}, conds ...interface{}) error
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

func (g *Gorm) Count(model interface{}, count *int64, query interface{}, args ...interface{}) error {
	const op = "storage.orm.Count"
	err := g.DB.Model(model).Where(query, args).Count(count).Error
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (g *Gorm) Create(dest interface{}) error {
	const op = "storage.orm.Create"
	if err := g.DB.Create(dest).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (g *Gorm) Updates(dest interface{}) error {
	const op = "storage.orm.Update"
	if err := g.DB.Updates(dest).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (g *Gorm) Delete(dest interface{}, conds ...interface{}) error {
	const op = "storage.orm.Delete"
	if err := g.DB.Delete(dest, conds...).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (g *Gorm) First(model interface{}, fields *[]string, dest interface{}, conds ...interface{}) error {
	const op = "storage.orm.First"
	query := g.DB.Model(model)
	if fields != nil {
		query.Select(*fields)
	}

	if err := query.First(dest, conds...).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (g *Gorm) Find(model interface{}, fields *[]string, offset, limit *int, dest interface{}, conds ...interface{}) error {
	const op = "storage.orm.Find"
	query := g.DB.Model(model)
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
