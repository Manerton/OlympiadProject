package school

import "github.com/google/uuid"

type School struct {
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name   string
	Region string
}
