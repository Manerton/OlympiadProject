package orm

import (
	"fmt"

	"gorm.io/gorm"
)

type ORM interface {
	Create(interface{}) error
	Updates(interface{}) error
	Delete(interface{}, ...interface{}) error
	Find(interface{}, ...interface{}) error
	FindWithSelect([]string, interface{}, ...interface{}) error
	First(interface{}, ...interface{}) error
}

type Gorm struct {
	DB *gorm.DB
}

func NewGormORM(db *gorm.DB) ORM {
	return &Gorm{DB: db}
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

func (g *Gorm) Find(dest interface{}, conds ...interface{}) error {
	const op = "storage.orm.Find"
	if err := g.DB.Find(dest, conds...).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (g *Gorm) FindWithSelect(fields []string, dest interface{}, conds ...interface{}) error {
	const op = "storage.orm.Find"
	if err := g.DB.Select(fields).Find(dest, conds...).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (g *Gorm) First(dest interface{}, conds ...interface{}) error {
	const op = "storage.orm.First"
	if err := g.DB.First(dest, conds...).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
