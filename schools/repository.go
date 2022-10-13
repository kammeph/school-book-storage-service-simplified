package schools

import (
	"context"
	"database/sql"
)

type SchoolsRepository interface {
	GetAll(ctx context.Context) ([]SchoolDto, error)
	GetById(ctx context.Context, id string) (SchoolDto, error)
	Insert(ctx context.Context, school SchoolModel) error
	Update(ctx context.Context, school SchoolModel) error
	Delete(ctx context.Context, id, updatedBy string) error
}

type SqlSchoolsRepository struct {
	db *sql.DB
}

func NewSqlSchoolsRepository(db *sql.DB) SchoolsRepository {
	return &SqlSchoolsRepository{db}
}

func (r *SqlSchoolsRepository) GetAll(ctx context.Context) ([]SchoolDto, error) {
	const getAllQuery = "SELECT s.id, s.name FROM schools s WHERE s.active = TRUE"
	stmt, err := r.db.PrepareContext(ctx, getAllQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schools := []SchoolDto{}
	for rows.Next() {
		var school SchoolDto
		if err := rows.Scan(&school.ID, &school.Name); err != nil {
			return nil, err
		}
		schools = append(schools, school)
	}
	return schools, nil
}

func (r *SqlSchoolsRepository) GetById(ctx context.Context, id string) (SchoolDto, error) {
	const getByIdQuery = "SELECT s.id, s.name FROM schools s WHERE s.active = TRUE and s.id = ?"
	stmt, err := r.db.PrepareContext(ctx, getByIdQuery)
	if err != nil {
		return SchoolDto{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	if row.Err() != nil {
		return SchoolDto{}, err
	}

	var school SchoolDto
	err = row.Scan(&school.ID, &school.Name)
	return school, err
}

func (r *SqlSchoolsRepository) Insert(ctx context.Context, school SchoolModel) error {
	const createQuery = "INSERT INTO schools (id, name, active, created_at, created_by) VALUES (?, ?, ?, ?, ?)"
	stmt, err := r.db.PrepareContext(ctx, createQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, school.ID, school.Name, school.Active, school.CreatedAt, school.CreatedBy)
	return err
}

func (r *SqlSchoolsRepository) Update(ctx context.Context, school SchoolModel) error {
	const updateQuery = "UPDATE schools SET name = ?, updated_at = CURRENT_TIMESTAMP(), updated_by = ? WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, updateQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, school.Name, school.UpdatedBy, school.ID)
	return err
}

func (r *SqlSchoolsRepository) Delete(ctx context.Context, id, updatedBy string) error {
	const deleteQuery = "UPDATE schools SET active = FALSE, updated_at = CURRENT_TIMESTAMP(), updated_by = ? WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, deleteQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, updatedBy, id)
	return err
}
