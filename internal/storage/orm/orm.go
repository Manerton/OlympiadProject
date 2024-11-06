package orm

import "gorm.io/gorm"

type ORM interface {
	Create(interface{})
	Updates(interface{})
	Delete(interface{})
	Find(interface{})
	First(interface{})
}

type Gorm struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) ORM {
	return &Gorm{db: db}
}

func (g *Gorm) Create(interFace interface{}) {
	g.db.Create(interFace)
}

func (g *Gorm) Updates(interFace interface{}) {
	g.db.Updates(interFace)
}

func (g *Gorm) Delete(interFace interface{}) {
	g.db.Delete(interFace)
}

func (g *Gorm) Find(interFace interface{}) {
	g.db.Find(interFace)
}

func (g *Gorm) First(interFace interface{}) {
	g.db.First(interFace)
}
