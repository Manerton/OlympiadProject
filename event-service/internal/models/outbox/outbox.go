package outbox

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Outbox struct {
	ID          uuid.UUID
	Payload     datatypes.JSON
	Processed   bool
	QueueName   string
	CreatedAt   time.Time
	ProcessedAt *time.Time
}
