package subject

type Subject struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(100);not null;unique"`
}

type SubjectStorage struct {
	subjects []string
}

func NewSubjectsStorage() *SubjectStorage {
	subjects := []string{
		"Матиматика",
		"Русский язык",
		"Английский язык",
		"История",
	}
	return &SubjectStorage{subjects: subjects}
}

func (s *SubjectStorage) GetAllSubject() []string {
	return s.subjects
}
