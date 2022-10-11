package books

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service-simplified/common"
)

type BooksResponseModel struct {
	Books []Book `json:"books"`
}

func BooksResponse(w http.ResponseWriter, books []Book) {
	common.JsonResponse(w, BooksResponseModel{books})
}

type BookResponseModel struct {
	Book Book `json:"book"`
}

func BookResponse(w http.ResponseWriter, book Book) {
	common.JsonResponse(w, BookResponseModel{book})
}

type BooksController struct {
	repository BooksRepository
}

func NewBooksController(db *sql.DB) BooksController {
	return BooksController{NewSqlBooksRepository(db)}
}

func AddBooksController(db *sql.DB) {
	controller := NewBooksController(db)
	common.Get("/api/books/get-all", common.IsAllowedWithClaims(controller.GetAllBooks, []common.Role{common.User, common.Superuser, common.Admin, common.SysAdmin}))
	common.Get("/api/books/get-by-id", common.IsAllowed(controller.GetBookById, []common.Role{common.User, common.Superuser, common.Admin, common.SysAdmin}))
	common.Post("/api/books/add", common.IsAllowedWithClaims(controller.AddBook, []common.Role{common.Superuser, common.Admin, common.SysAdmin}))
	common.Post("/api/books/update", common.IsAllowedWithClaims(controller.UpdateBook, []common.Role{common.Superuser, common.Admin, common.SysAdmin}))
	common.Post("/api/books/delete", common.IsAllowed(controller.DeleteBook, []common.Role{common.Admin, common.SysAdmin}))
}

func (c BooksController) GetAllBooks(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	if claims.SchoolId == nil {
		common.HttpErrorResponse(w, "user is not assigned to a school")
		return
	}
	books, err := c.repository.GetAll(r.Context(), *claims.SchoolId)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	BooksResponse(w, books)
}

func (c BooksController) GetBookById(w http.ResponseWriter, r *http.Request) {
	bookId := r.URL.Query().Get("id")
	if bookId == "" {
		common.HttpErrorResponse(w, "no book id specified")
		return
	}
	book, err := c.repository.GetById(r.Context(), bookId)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.JsonResponse(w, book)
}

func (c BooksController) AddBook(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	if claims.SchoolId == nil {
		common.HttpErrorResponse(w, "user is not assigned to a school")
		return
	}
	var book Book
	book.ID = uuid.NewString()
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	if err := c.repository.Create(r.Context(), book, *claims.SchoolId, claims.UserId); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}

func (c BooksController) UpdateBook(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	if claims.SchoolId == nil {
		common.HttpErrorResponse(w, "user is not assigned to a school")
		return
	}
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	if err := c.repository.Update(r.Context(), book, claims.UserId); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}

func (c BooksController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookId := r.URL.Query().Get("id")
	if bookId == "" {
		common.HttpErrorResponse(w, "no school id specified")
		return
	}
	if err := c.repository.Delete(r.Context(), bookId); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}
