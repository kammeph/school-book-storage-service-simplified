package schoolclasses

import (
	"context"
	"database/sql"
)

type SchoolClassesRepository interface {
	GetAll(ctx context.Context, schoolId string) ([]SchoolClass, error)
	GetById(ctx context.Context, id string) (SchoolClass, error)
	Insert(ctx context.Context, schoolClass SchoolClass, schoolId, createdBy string) error
	Update(ctx context.Context, schoolClass SchoolClass, updatedBy string) error
	Delete(ctx context.Context, id string) error
}

type SqlSchoolClassesRepository struct {
	db *sql.DB
}

func NewSqlSchoolsRepository(db *sql.DB) SchoolClassesRepository {
	return &SqlSchoolClassesRepository{db}
}

func (r *SqlSchoolClassesRepository) GetAll(ctx context.Context, schoolId string) ([]SchoolClass, error) {
	const getAllQuery = "SELECT id, grade, letter, number_of_pupils, date_from, date_to FROM school_classes WHERE school_id = ?"
	stmt, err := r.db.PrepareContext(ctx, getAllQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, schoolId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schoolClasses := []SchoolClass{}
	for rows.Next() {
		var schoolClass SchoolClass
		if err := rows.Scan(&schoolClass.ID, &schoolClass.Grade, &schoolClass.Letter, &schoolClass.NumberOfPupils, &schoolClass.DateFrom, &schoolClass.DateTo); err != nil {
			return nil, err
		}
		schoolClasses = append(schoolClasses, schoolClass)
	}
	return schoolClasses, nil
}

func (r *SqlSchoolClassesRepository) GetById(ctx context.Context, id string) (SchoolClass, error) {
	const getByIdQuery = "SELECT id, grade, letter, number_of_pupils, date_from, date_to FROM school_classes WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, getByIdQuery)
	if err != nil {
		return SchoolClass{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	if row.Err() != nil {
		return SchoolClass{}, err
	}

	var schoolClass SchoolClass
	err = row.Scan(&schoolClass.ID, &schoolClass.Grade, &schoolClass.Letter, &schoolClass.NumberOfPupils, &schoolClass.DateFrom, &schoolClass.DateTo)
	return schoolClass, err
}

func (r *SqlSchoolClassesRepository) Insert(ctx context.Context, schoolClass SchoolClass, schoolId, createdBy string) error {
	const createQuery = "INSERT INTO school_classes (id, school_id, grade, letter, number_of_pupils, date_from, date_to, created_at, created_by) VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP(), ?)"
	stmt, err := r.db.PrepareContext(ctx, createQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, schoolClass.ID, schoolId, schoolClass.Grade, schoolClass.Letter, schoolClass.NumberOfPupils, schoolClass.DateFrom, schoolClass.DateTo, createdBy)
	return err
}

func (r *SqlSchoolClassesRepository) Update(ctx context.Context, schoolClass SchoolClass, updatedBy string) error {
	const updateQuery = "UPDATE school_classes SET grade = ?, letter = ?, number_of_pupils = ?, date_from = ?, date_to = ?, updated_at = CURRENT_TIMESTAMP(), updated_by = ? WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, updateQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, schoolClass.Grade, schoolClass.Letter, schoolClass.NumberOfPupils, schoolClass.DateFrom, schoolClass.DateTo, updatedBy, schoolClass.ID)
	return err
}

func (r *SqlSchoolClassesRepository) Delete(ctx context.Context, id string) error {
	const deleteQuery = "DELETE FROM school_classes WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, deleteQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	return err
}
