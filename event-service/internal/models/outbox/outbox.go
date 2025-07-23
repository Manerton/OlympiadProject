package outbox

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Outbox struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Payload     datatypes.JSON
	Processed   bool
	QueueName   string
	CreatedAt   time.Time
	ProcessedAt *time.Time
}

func (Outbox) TableName() string {
	return "outbox"
}
