package subject

type Subject struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(100);not null;unique"`
}

type SubjectStorage struct {
	subjects map[int]string
}

const (
	SubjectHistory = iota + 1
	SubjectBiology
	SubjectGeography
	SubjectForeignLanguage
	SubjectInformatics
	SubjectLiterature
	SubjectMath
	SubjectSocialScience
	SubjectRussian
	SubjectPhysics
	SubjectChemistry
)

func NewSubjectsStorage() *SubjectStorage {
	var SubjectIDToName = map[int]string{
		1:  "История",
		2:  "Биология",
		3:  "География",
		4:  "Иностранный язык",
		5:  "Информатика и ИКТ",
		6:  "Литература",
		7:  "Математика",
		8:  "Обществознание",
		9:  "Русский язык",
		10: "Физика",
		11: "Химия",
	}
	return &SubjectStorage{subjects: SubjectIDToName}
}

func (s *SubjectStorage) GetAllSubject() map[int]string {
	return s.subjects
}
