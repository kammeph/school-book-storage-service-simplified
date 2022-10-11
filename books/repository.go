package books

import (
	"context"
	"database/sql"
	"errors"
)

type BooksRepository interface {
	GetAll(ctx context.Context, schoolId string) ([]Book, error)
	GetById(ctx context.Context, id string) (Book, error)
	Create(ctx context.Context, book Book, schoolId, createdBy string) error
	Update(ctx context.Context, book Book, updatedBy string) error
	Delete(ctx context.Context, id string) error
}

type SqlBooksRepository struct {
	db *sql.DB
}

func NewSqlBooksRepository(db *sql.DB) BooksRepository {
	return &SqlBooksRepository{db}
}

func (r *SqlBooksRepository) GetAll(ctx context.Context, schoolId string) ([]Book, error) {
	const getAllBooksQuery = "SELECT b.id, b.isbn, b.name, b.description, b.subject, b.price, g.grade FROM books b INNER JOIN book_grades g ON b.id = g.book_id WHERE b.school_id = ?"
	stmt, err := r.db.PrepareContext(ctx, getAllBooksQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, schoolId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanBooks(rows)
}

func (r *SqlBooksRepository) GetById(ctx context.Context, id string) (Book, error) {
	const getBookByIdQuery = "SELECT b.id, b.isbn, b.name, b.description, b.subject, b.price, g.grade FROM books b INNER JOIN book_grades g ON b.id = g.book_id WHERE b.id = ?"
	stmt, err := r.db.PrepareContext(ctx, getBookByIdQuery)
	if err != nil {
		return Book{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return Book{}, err
	}
	defer rows.Close()

	books, err := scanBooks(rows)
	if err != nil {
		return Book{}, nil
	}

	if len(books) == 0 {
		return Book{}, nil
	}

	if len(books) > 1 {
		return Book{}, errors.New("more than one book found")
	}

	return books[0], nil
}

func (r *SqlBooksRepository) Create(ctx context.Context, book Book, schoolId, createdBy string) error {
	const insertBookQuery = "INSERT INTO books (id, school_id, isbn, name, description, subject, price, created_at, created_by) VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP(), ?)"
	stmt, err := r.db.PrepareContext(ctx, insertBookQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, book.ID, schoolId, book.ISBN, book.Name, book.Description, book.Subject, book.Price, createdBy)
	if err != nil {
		return err
	}
	return r.insertGrades(ctx, book.ID, book.Grades)
}

func (r *SqlBooksRepository) Update(ctx context.Context, book Book, updatedBy string) error {
	const updateBookQuery = "UPDATE books SET isbn = ?, name = ?, description = ?, subject = ?, price = ?, updated_at = CURRENT_TIMESTAMP(), updated_by = ? WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, updateBookQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, book.ISBN, book.Name, book.Description, book.Subject, book.Price, updatedBy, book.ID)
	if err != nil {
		return err
	}

	return r.updateGrades(ctx, book.ID, book.Grades)
}

func (r *SqlBooksRepository) Delete(ctx context.Context, bookId string) error {
	err := r.deleteGrades(ctx, bookId)
	if err != nil {
		return err
	}
	const deleteBookQuery = "DELETE FROM books WHERE id = ?"
	stmt, err := r.db.PrepareContext(ctx, deleteBookQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, bookId)
	return err
}

func (r *SqlBooksRepository) insertGrades(ctx context.Context, bookId string, grades []int) error {
	const insertGradeQuery = "INSERT INTO book_grades (book_id, grade) VALUES (?, ?)"
	stmt, err := r.db.PrepareContext(ctx, insertGradeQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, grade := range grades {
		_, err = stmt.ExecContext(ctx, bookId, grade)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SqlBooksRepository) deleteGrades(ctx context.Context, bookId string) error {
	const deleteGradeQuery = "DELETE FROM book_grades WHERE book_id = ?"
	stmt, err := r.db.PrepareContext(ctx, deleteGradeQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, bookId)
	return err
}

func (r *SqlBooksRepository) updateGrades(ctx context.Context, bookId string, grades []int) error {
	err := r.deleteGrades(ctx, bookId)
	if err != nil {
		return err
	}
	return r.insertGrades(ctx, bookId, grades)
}

func scanBooks(rows *sql.Rows) ([]Book, error) {
	books := []Book{}
	previousID := ""
	for rows.Next() {
		var book Book
		var grade int
		err := rows.Scan(&book.ID, &book.ISBN, &book.Name, &book.Description, &book.Subject, &book.Price, &grade)
		if err != nil {
			return nil, err
		}
		if book.ID != previousID {
			book.Grades = append(book.Grades, grade)
			books = append(books, book)
			previousID = book.ID
		} else {
			books[len(books)-1].Grades = append(books[len(books)-1].Grades, grade)
		}
	}
	return books, nil
}
