package eventrepository

import (
	"gorm.io/gorm"
)

type eventRepositoryInterfaca interface {
	GetEventByID()
	GetEventByType()
	GetAllEvents()

	CreateEvent()
	UpdateEvent()
	DeleteEvent()
}

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) eventRepositoryInterfaca {
	return &EventRepository{db: db}
}

func (r *EventRepository) GetEventByID() {

}

func (r *EventRepository) GetEventByType() {

}

func (r *EventRepository) GetAllEvents() {

}

func (r *EventRepository) CreateEvent() {

}

func (r *EventRepository) UpdateEvent() {

}

func (r *EventRepository) DeleteEvent() {

}
