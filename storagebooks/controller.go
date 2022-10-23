package storagebooks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kammeph/school-book-storage-service-simplified/common"
	"github.com/kammeph/school-book-storage-service-simplified/fp"
)

type StoreBooksBody struct {
	StorageId string `json:"storageId"`
	Books     []Book `json:"books"`
}

type AddBookBody struct {
	StorageId string `json:"storageId"`
	Book      Book   `json:"book"`
}

type BooksResponseModel struct {
	Books []Book `json:"books"`
}

func StorageBooksResponse(w http.ResponseWriter, books []Book) {
	common.JsonResponse(w, BooksResponseModel{books})
}

type BookResponseModel struct {
	Book Book `json:"book"`
}

func StorageBookResponse(w http.ResponseWriter, book Book) {
	common.JsonResponse(w, BookResponseModel{book})
}

type StorageBooksController struct {
	repository StorageBooksRepository
}

func NewStorageBooksController(db *sql.DB) StorageBooksController {
	return StorageBooksController{NewSqlStorageBooksRepository(db)}
}

func AddStorageBooksController(db *sql.DB) {
	controller := NewStorageBooksController(db)
	common.Get("/api/storagebooks/get-books-by-storage", common.IsAllowed(controller.GetBooksByStorage, []common.Role{common.User, common.Superuser, common.Admin, common.SysAdmin}))
	common.Post("/api/storagebooks/store-books", common.IsAllowedWithClaims(controller.StoreBooks, []common.Role{common.Superuser, common.Admin, common.SysAdmin}))
	common.Post("/api/storagebooks/add-book", common.IsAllowedWithClaims(controller.AddBooksToStorage, []common.Role{common.Superuser, common.Admin, common.SysAdmin}))
	common.Post("/api/storagebooks/change-books-count", common.IsAllowedWithClaims(controller.ChangeBooksCountInStorage, []common.Role{common.Superuser, common.Admin, common.SysAdmin}))
	common.Post("/api/storagebooks/remove-book-from-storage", common.IsAllowed(controller.RemoveBookFromStorge, []common.Role{common.Admin, common.SysAdmin}))
}

func (c StorageBooksController) GetBooksByStorage(w http.ResponseWriter, r *http.Request) {
	storageId := r.URL.Query().Get("id")
	if storageId == "" {
		common.HttpErrorResponse(w, "no book id specified")
		return
	}
	books, err := c.repository.GetBooksByStorage(r.Context(), storageId)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	StorageBooksResponse(w, books)
}

func (c StorageBooksController) StoreBooks(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	var body StoreBooksBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	if body.StorageId == "" {
		common.HttpErrorResponse(w, "no storage id specified")
		return
	}
	if len(body.Books) <= 0 {
		common.HttpErrorResponse(w, "no books submitted")
		return
	}
	if fp.Some(body.Books, func(book Book) bool { return book.Count < 0 }) {
		common.HttpErrorResponse(w, "not all book counts are valid")
		return
	}
	if err := c.repository.StoreBooks(r.Context(), body.Books, body.StorageId, claims.UserId); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}

func (c StorageBooksController) AddBooksToStorage(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	var body AddBookBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	fmt.Println(body)
	if body.StorageId == "" {
		common.HttpErrorResponse(w, "no storage id specified")
		return
	}
	if body.Book == (Book{}) {
		common.HttpErrorResponse(w, "no book submitted")
		return
	}
	if body.Book.Count < 0 {
		common.HttpErrorResponse(w, "book count is not valid")
		return
	}
	if err := c.repository.AddBookToStorage(r.Context(), body.Book, body.StorageId, claims.UserId); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}

func (c StorageBooksController) ChangeBooksCountInStorage(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	var body StoreBooksBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	if body.StorageId == "" {
		common.HttpErrorResponse(w, "no storage id specified")
		return
	}
	if len(body.Books) <= 0 {
		common.HttpErrorResponse(w, "no books submitted")
		return
	}
	if fp.Some(body.Books, func(book Book) bool { return book.Count < 0 }) {
		common.HttpErrorResponse(w, "not all book counts are valid")
		return
	}
	if err := c.repository.ChangeBooksCountInStorage(r.Context(), body.Books, body.StorageId, claims.UserId); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}

func (c StorageBooksController) RemoveBookFromStorge(w http.ResponseWriter, r *http.Request) {
	storageId := r.URL.Query().Get("storageId")
	if storageId == "" {
		common.HttpErrorResponse(w, "no storage id specified")
		return
	}
	bookId := r.URL.Query().Get("bookId")
	if bookId == "" {
		common.HttpErrorResponse(w, "no book id specified")
		return
	}
	if err := c.repository.RemoveBookFromStorage(r.Context(), storageId, bookId); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}
