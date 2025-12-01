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
		1:  "Экономика",
		2:  "Биология",
		3:  "Испанский язык",
		4:  "Информатика",
		5:  "История",
		6:  "Химия",
		7:  "Русский язык",
		8:  "Итальянский язык",
		9:  "Китайский язык",
		10: "Обществознание",
		11: "Физика",
		12: "Математика",
		13: "Немецкий язык",
		14: "Экология",
		15: "Литература",
		16: "Право",
		17: "Труд",
		18: "География",
		19: "Французский язык",
		20: "Английский язык",
		21: "Физическая культура",
		22: "Астрономия",
		23: "Основы безопастности и защиты Родины",
		24: "Искусство",
	}
	return &SubjectStorage{subjects: SubjectIDToName}
}

func (s *SubjectStorage) GetAllSubject() map[int]string {
	return s.subjects
}
