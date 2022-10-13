package storages

import (
	"context"
	"database/sql"
)

type StoragesRepository interface {
	GetAll(ctx context.Context, schoolId string) ([]Storage, error)
	GetById(ctx context.Context, id string) (Storage, error)
	Insert(ctx context.Context, storage Storage, schoolId, createdBy string) error
	Update(ctx context.Context, storage Storage, updatedBy string) error
	Delete(ctx context.Context, id string) error
}

type SqlStoragesRepository struct {
	db *sql.DB
}

func NewSqlSchoolsRepository(db *sql.DB) StoragesRepository {
	return &SqlStoragesRepository{db}
}

func (r *SqlStoragesRepository) GetAll(ctx context.Context, schoolId string) ([]Storage, error) {
	const getAllQuery = "SELECT id, name, location FROM storages WHERE school_id = ?"
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

	storagees := []Storage{}
	for rows.Next() {
		var storage Storage
		if err := rows.Scan(&storage.ID, &storage.Name, &storage.Location); err != nil {
			return nil, err
		}
		storagees = append(storagees, storage)
	}
	return storagees, nil
}

func (r *SqlStoragesRepository) GetById(ctx context.Context, id string) (Storage, error) {
	const getByIdQuery = "SELECT id, name, location FROM storages WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, getByIdQuery)
	if err != nil {
		return Storage{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)
	if row.Err() != nil {
		return Storage{}, err
	}

	var storage Storage
	err = row.Scan(&storage.ID, &storage.Name, &storage.Location)
	return storage, err
}

func (r *SqlStoragesRepository) Insert(ctx context.Context, storage Storage, schoolId, createdBy string) error {
	const createQuery = "INSERT INTO storages (id, school_id, name, location, created_at, created_by) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP(), ?)"
	stmt, err := r.db.PrepareContext(ctx, createQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, storage.ID, schoolId, storage.Name, storage.Location, createdBy)
	return err
}

func (r *SqlStoragesRepository) Update(ctx context.Context, storage Storage, updatedBy string) error {
	const updateQuery = "UPDATE storages SET name = ?, location = ?, updated_at = CURRENT_TIMESTAMP(), updated_by = ? WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, updateQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, storage.Name, storage.Location, updatedBy, storage.ID)
	return err
}

func (r *SqlStoragesRepository) Delete(ctx context.Context, id string) error {
	const deleteQuery = "DELETE FROM storages WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, deleteQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	return err
}
