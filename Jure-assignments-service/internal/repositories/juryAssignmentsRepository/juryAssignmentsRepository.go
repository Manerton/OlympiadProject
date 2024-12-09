package juryAssignmentsRepository

import (
	"fmt"
	"main/internal/models/juryAssignments"
	"main/internal/storage/orm"
)

type JuryAssignmentsRepository struct{}

func NewJuryAssignmentsRepository() *JuryAssignmentsRepository {
	return &JuryAssignmentsRepository{}
}

func (r *JuryAssignmentsRepository) GetJuryAssignmentsByFilter(
	myOrm orm.ORM, filter juryAssignments.JuryAssignments) (juryAssignments.JuryAssignments, error) {
	const op = "repositories.juryAssignmentsRepository.GetJuryAssignmentsByID"
	if err := myOrm.First(&filter); err != nil {
		return juryAssignments.JuryAssignments{}, fmt.Errorf("%s: %w", op, err)
	}
	return filter, nil
}

func (r *JuryAssignmentsRepository) GetPartOfAllJuryAssignmentsByFilter(
	myOrm orm.ORM, fields []string, filter juryAssignments.JuryAssignments) ([]juryAssignments.JuryAssignments, error) {
	const op = "repositories.juryAssignmentsRepository.GetPartOfAllJuryAssignmentsByFilter"
	partOfJuryAssignmentsRes := []juryAssignments.JuryAssignments{}
	if err := myOrm.FindWithSelect(fields, &partOfJuryAssignmentsRes, filter); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return partOfJuryAssignmentsRes, nil
}

func (r *JuryAssignmentsRepository) GetAllJuryAssignments(myOrm orm.ORM, conds ...interface{}) ([]juryAssignments.JuryAssignments, error) {
	const op = "repositories.juryAssignmentsRepository.GetAllJuryAssignments"
	juryAssignmentsRes := []juryAssignments.JuryAssignments{}
	if err := myOrm.Find(&juryAssignmentsRes, conds...); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return juryAssignmentsRes, nil
}

func (r *JuryAssignmentsRepository) GetAllJuryAssignmentsByFilter(
	myOrm orm.ORM, filter juryAssignments.JuryAssignments) ([]juryAssignments.JuryAssignments, error) {
	const op = "repositories.juryAssignmentsRepository.GetAllJuryAssignmentsByFilter"
	juryAssignmentsRes := []juryAssignments.JuryAssignments{}
	if err := myOrm.Find(&juryAssignmentsRes, filter); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return juryAssignmentsRes, nil
}

func (r *JuryAssignmentsRepository) CreateJuryAssignments(
	myOrm orm.ORM, juryAssignments juryAssignments.JuryAssignments) (uint, error) {
	const op = "repositories.juryAssignmentsRepository.CreateJuryAssignments"

	if err := myOrm.Create(&juryAssignments); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return juryAssignments.ID, nil
}

func (r *JuryAssignmentsRepository) UpdateJuryAssignments(
	myOrm orm.ORM, juryAssignments juryAssignments.JuryAssignments) (uint, error) {
	const op = "repositories.juryAssignmentsRepository.UpdateJuryAssignments"

	if err := myOrm.Updates(&juryAssignments); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return juryAssignments.ID, nil
}

func (r *JuryAssignmentsRepository) DeleteJuryAssignments(myOrm orm.ORM, id uint) error {
	const op = "repositories.juryAssignmentsRepository.DeleteJuryAssignments"
	if err := myOrm.Delete(&juryAssignments.JuryAssignments{ID: id}); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
