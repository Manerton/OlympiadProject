package subject

type Subject struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(100);not null"`
}
