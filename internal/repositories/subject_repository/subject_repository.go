package subject_repository

import (
	"fmt"
	"main/internal/models/subject"

	"gorm.io/gorm"
)

type SubjectRepository struct{}

// Get subject by id
func (r *SubjectRepository) GetSubjectByID(db *gorm.DB, id uint) (subject.Subject, error) {
	const op = "repositories.subject_repository.GetSubjectById"

	subject_res := subject.Subject{ID: id}
	// test_subject := subject.Subject{}
	// test_err := db.Where("id = ?", id).First(&test_subject).Error
	// _ = test_err
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
func (r *SubjectRepository) CreateSubject(db *gorm.DB, subject subject.Subject) (uint, error) {
	const op = "repositories.subject_repository.CreateSubject"
	if err := db.Create(&subject).Error; err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return subject.ID, nil
}

// Update subject
func (r *SubjectRepository) UpdateSubject(db *gorm.DB, subject subject.Subject) (uint, error) {
	const op = "repositories.subject_repository.UpdateSubject"
	if err := db.Updates(&subject).Error; err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return subject.ID, nil
}

// Delete subject
func (r *SubjectRepository) DeleteSubject(db *gorm.DB, id uint) error {
	const op = "repositories.subject_repository.DeleteSubject"
	if err := db.Delete(subject.Subject{}, id).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
