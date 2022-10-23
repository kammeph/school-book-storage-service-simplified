package storagebooks

import (
	"context"
	"database/sql"

	"github.com/kammeph/school-book-storage-service-simplified/fp"
)

type StorageBooksRepository interface {
	GetBooksByStorage(ctx context.Context, storageId string) ([]Book, error)
	StoreBooks(ctx context.Context, books []Book, storageId, requestor string) error
	AddBookToStorage(ctx context.Context, book Book, storageId, createdBy string) error
	ChangeBooksCountInStorage(ctx context.Context, books []Book, storageId, updatedBy string) error
	RemoveBookFromStorage(ctx context.Context, storageId, bookId string) error
}

type SqlStorageBooksRepository struct {
	db *sql.DB
}

func NewSqlStorageBooksRepository(db *sql.DB) StorageBooksRepository {
	return &SqlStorageBooksRepository{db}
}

func (r *SqlStorageBooksRepository) GetBooksByStorage(ctx context.Context, storageId string) ([]Book, error) {
	const getBooksByStorageQuery = "SELECT b.id, b.name, sb.count FROM storage_book sb INNER JOIN books b ON sb.book_id = b.id WHERE sb.storage_id = ?"
	stmt, err := r.db.PrepareContext(ctx, getBooksByStorageQuery)
	books := []Book{}
	if err != nil {
		return books, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, storageId)
	if err != nil {
		return books, err
	}
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Name, &book.Count); err != nil {
			return books, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *SqlStorageBooksRepository) StoreBooks(ctx context.Context, books []Book, storageId, requestor string) error {
	bookIds, err := r.getBookIdsByStorage(ctx, storageId)
	if err != nil {
		return err
	}
	const insertQuery = "INSERT INTO storage_book (storage_id, book_id, count, created_at, created_by) VALUES (?, ?, ?, CURRENT_TIMESTAMP(), ?)"
	insertStmt, err := r.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		return err
	}
	defer insertStmt.Close()
	const updateQuery = "UPDATE storage_book SET count = ?, updated_at = CURRENT_TIMESTAMP(), updated_by = ? WHERE storage_id = ? AND book_id = ?"
	updateStmt, err := r.db.PrepareContext(ctx, updateQuery)
	if err != nil {
		return err
	}
	defer updateStmt.Close()

	for _, book := range books {
		if len(bookIds) <= 0 || !fp.Some(bookIds, func(bookId string) bool { return bookId == book.ID }) {
			_, err = insertStmt.ExecContext(ctx, storageId, book.ID, book.Count, requestor)
			if err != nil {
				return err
			}
		} else {
			_, err = updateStmt.ExecContext(ctx, book.Count, requestor, storageId, book.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *SqlStorageBooksRepository) AddBookToStorage(ctx context.Context, book Book, storageId, createdBy string) error {
	const insertQuery = "INSERT INTO storage_book (storage_id, book_id, count, created_at, created_by) VALUES (?, ?, ?, CURRENT_TIMESTAMP(), ?)"
	stmt, err := r.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, storageId, book.ID, book.Count, createdBy)
	return err
}

func (r *SqlStorageBooksRepository) ChangeBooksCountInStorage(ctx context.Context, books []Book, storageId, updatedBy string) error {
	const updateQuery = "UPDATE storage_book SET count = ?, updated_at = CURRENT_TIMESTAMP(), updated_by = ? WHERE storage_id = ? AND book_id = ?"
	stmt, err := r.db.PrepareContext(ctx, updateQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, book := range books {
		_, err = stmt.ExecContext(ctx, book.Count, updatedBy, storageId, book.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SqlStorageBooksRepository) RemoveBookFromStorage(ctx context.Context, storageId, bookId string) error {
	const updateBooksInStorageQuery = "DELETE FROM storage_book WHERE storage_id = ? AND book_id = ?"
	stmt, err := r.db.PrepareContext(ctx, updateBooksInStorageQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, storageId, bookId)
	return err
}

func (r *SqlStorageBooksRepository) getBookIdsByStorage(ctx context.Context, storageId string) ([]string, error) {
	bookIds := []string{}
	const getBooksByStorageQuery = "SELECT book_id FROM storage_book WHERE storage_id = ?"
	stmt, err := r.db.PrepareContext(ctx, getBooksByStorageQuery)
	if err != nil {
		return bookIds, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, storageId)
	if err != nil {
		return bookIds, err
	}
	for rows.Next() {
		var bookId string
		if err := rows.Scan(&bookId); err != nil {
			return bookIds, err
		}
		bookIds = append(bookIds, bookId)
	}
	return bookIds, nil
}
