package juryAssignmentsRepository

import (
	"fmt"
	"main/internal/models/juryAssignments"

	"gorm.io/gorm"
)

type JuryAssignmentsRepository struct{}

func NewJuryAssignmentsRepository() *JuryAssignmentsRepository {
	return &JuryAssignmentsRepository{}
}

func (r *JuryAssignmentsRepository) GetJuryAssignmentsByFilter(
	db *gorm.DB, filter juryAssignments.JuryAssignments) (juryAssignments.JuryAssignments, error) {
	const op = "repositories.juryAssignmentsRepository.GetJuryAssignmentsByID"
	if err := db.First(&filter).Error; err != nil {
		return juryAssignments.JuryAssignments{}, fmt.Errorf("%s: %w", op, err)
	}
	return filter, nil
}

func (r *JuryAssignmentsRepository) GetPartOfAllJuryAssignmentsByFilter(
	db *gorm.DB, fields []string, filter juryAssignments.JuryAssignments) ([]*juryAssignments.JuryAssignments, error) {
	const op = "repositories.juryAssignmentsRepository.GetPartOfGetJuryAssignmentsByFilter"
	partOfJuryAssignmentsRes := []*juryAssignments.JuryAssignments{}
	if err := db.Select(fields).First(&partOfJuryAssignmentsRes, filter).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return partOfJuryAssignmentsRes, nil
}

func (r *JuryAssignmentsRepository) GetAllJuryAssignments(db *gorm.DB) ([]*juryAssignments.JuryAssignments, error) {
	const op = "repositories.juryAssignmentsRepository.GetAllJuryAssignments"
	juryAssignmentsRes := []*juryAssignments.JuryAssignments{}
	if err := db.Find(&juryAssignmentsRes).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return juryAssignmentsRes, nil
}

func (r *JuryAssignmentsRepository) GetAllJuryAssignmentsByFilter(
	db *gorm.DB, filter juryAssignments.JuryAssignments) ([]*juryAssignments.JuryAssignments, error) {
	const op = "repositories.juryAssignmentsRepository.GetAllJuryAssignmentsByFilter"
	juryAssignmentsRes := []*juryAssignments.JuryAssignments{}
	err := db.Find(&juryAssignmentsRes, filter).Error
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return juryAssignmentsRes, nil
}

func (r *JuryAssignmentsRepository) CreateJuryAssignments(
	db *gorm.DB, juryAssignments juryAssignments.JuryAssignments) (uint, error) {
	const op = "repositories.juryAssignmentsRepository.CreateJuryAssignments"

	if err := db.Create(&juryAssignments).Error; err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return juryAssignments.ID, nil
}

func (r *JuryAssignmentsRepository) UpdateJuryAssignments(
	db *gorm.DB, juryAssignments juryAssignments.JuryAssignments) (uint, error) {
	const op = "repositories.juryAssignmentsRepository.UpdateJuryAssignments"

	if err := db.Updates(&juryAssignments).Error; err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return juryAssignments.ID, nil
}

func (r *JuryAssignmentsRepository) DeleteJuryAssignments(db *gorm.DB, id uint) error {
	const op = "repositories.juryAssignmentsRepository.DeleteJuryAssignments"
	if err := db.Delete(&juryAssignments.JuryAssignments{ID: id}).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
