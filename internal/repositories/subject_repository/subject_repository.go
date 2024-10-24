package subjectrepository

import (
	"main/internal/dto/subject_dto"
	"main/internal/models/subject"

	"gorm.io/gorm"
)

type SubjectRepository struct{}

// Get subject by id
func (r *SubjectRepository) GetSubjectById(db *gorm.DB, id uint) (subject.Subject, error) {
	subject_res := subject.Subject{ID: id}
	if err := db.First(&subject_res).Error; err != nil {
		return subject.Subject{}, err
	}
	return subject_res, nil
}

// Get all subjects
func (r *SubjectRepository) GetAllSubjects(db *gorm.DB) ([]subject.Subject, error) {
	subjects_res := []subject.Subject{}
	if err := db.Find(&subjects_res).Error; err != nil {
		return nil, err
	}
	return subjects_res, nil
}

// Add new subject in DB
func (r *SubjectRepository) CreateSubject(db *gorm.DB, subject_dto subject_dto.SubjectDTO) (uint, error) {
	new_subject := subject.Subject{Name: subject_dto.Name}
	if err := db.Create(&new_subject).Error; err != nil {
		return 0, err
	}
	return new_subject.ID, nil
}

// Update subject
func (r *SubjectRepository) UpdateSubject(db *gorm.DB, subject_dto subject_dto.SubjectDTO) (uint, error) {
	update_subject := subject.Subject{
		ID:   subject_dto.ID,
		Name: subject_dto.Name,
	}
	if err := db.Updates(&update_subject).Error; err != nil {
		return 0, err
	}
	return update_subject.ID, nil
}

// Delete subject
func (r *SubjectRepository) DeleteSubject(db *gorm.DB, id uint) error {
	if err := db.Delete(subject.Subject{}, id).Error; err != nil {
		return err
	}
	return nil
}
