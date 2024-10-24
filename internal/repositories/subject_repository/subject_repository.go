package subjectrepository

import (
	"fmt"
	"main/internal/dto/subject_dto"
	"main/internal/models/subject"

	"gorm.io/gorm"
)

type SubjectRepository struct{}

// Get subject by id
func (r *SubjectRepository) GetSubjectById(db *gorm.DB, id uint) (subject.Subject, error) {
	const op = "repositories.subject_repository.GetSubjectById"
	subject_res := subject.Subject{ID: id}
	if err := db.First(&subject_res).Error; err != nil {
		return subject.Subject{}, fmt.Errorf("%s: %w", op, err)
	}
	return subject_res, nil
}

// Get all subjects
func (r *SubjectRepository) GetAllSubjects(db *gorm.DB) ([]subject.Subject, error) {
	const op = "repositories.subject_repository.GetAllSubjects"
	subjects_res := []subject.Subject{}
	if err := db.Find(&subjects_res).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return subjects_res, nil
}

// Add new subject in DB
func (r *SubjectRepository) CreateSubject(db *gorm.DB, subject_dto subject_dto.SubjectDTO) (uint, error) {
	const op = "repositories.subject_repository.CreateSubject"
	new_subject := subject.Subject{Name: subject_dto.Name}
	if err := db.Create(&new_subject).Error; err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return new_subject.ID, nil
}

// Update subject
func (r *SubjectRepository) UpdateSubject(db *gorm.DB, subject_dto subject_dto.SubjectDTO) (uint, error) {
	const op = "repositories.subject_repository.UpdateSubject"
	update_subject := subject.Subject{
		ID:   subject_dto.ID,
		Name: subject_dto.Name,
	}
	if err := db.Updates(&update_subject).Error; err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return update_subject.ID, nil
}

// Delete subject
func (r *SubjectRepository) DeleteSubject(db *gorm.DB, id uint) error {
	const op = "repositories.subject_repository.DeleteSubject"
	if err := db.Delete(subject.Subject{}, id).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
