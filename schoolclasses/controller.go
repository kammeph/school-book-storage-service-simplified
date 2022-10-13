package schoolclasses

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service-simplified/common"
)

type SchoolClassesResponseModel struct {
	SchoolClasses []SchoolClass `json:"schoolClasses"`
}

func BooksResponse(w http.ResponseWriter, schoolClasses []SchoolClass) {
	common.JsonResponse(w, SchoolClassesResponseModel{schoolClasses})
}

type SchoolClassResponseModel struct {
	SchoolClass SchoolClass `json:"schoolClass"`
}

func SchoolClassResponse(w http.ResponseWriter, schoolClass SchoolClass) {
	common.JsonResponse(w, SchoolClassResponseModel{schoolClass})
}

type SchoolClassController struct {
	repository SchoolClassesRepository
}

func NewBooksController(db *sql.DB) SchoolClassController {
	return SchoolClassController{NewSqlSchoolsRepository(db)}
}

func AddBooksController(db *sql.DB) {
	controller := NewBooksController(db)
	common.Get("/api/schoolclasses/get-all", common.IsAllowedWithClaims(controller.GetAllSchoolClasses, []common.Role{common.User, common.Superuser, common.Admin, common.SysAdmin}))
	common.Get("/api/schoolclasses/get-by-id", common.IsAllowed(controller.GetBookById, []common.Role{common.User, common.Superuser, common.Admin, common.SysAdmin}))
	common.Post("/api/schoolclasses/add", common.IsAllowedWithClaims(controller.AddBook, []common.Role{common.Superuser, common.Admin, common.SysAdmin}))
	common.Post("/api/schoolclasses/update", common.IsAllowedWithClaims(controller.UpdateBook, []common.Role{common.Superuser, common.Admin, common.SysAdmin}))
	common.Post("/api/schoolclasses/delete", common.IsAllowed(controller.DeleteBook, []common.Role{common.Admin, common.SysAdmin}))
}

func (c SchoolClassController) GetAllSchoolClasses(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	if claims.SchoolId == nil {
		common.HttpErrorResponse(w, "user is not assigned to a school")
		return
	}
	schoolClasses, err := c.repository.GetAll(r.Context(), *claims.SchoolId)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	BooksResponse(w, schoolClasses)
}

func (c SchoolClassController) GetBookById(w http.ResponseWriter, r *http.Request) {
	schoolClassId := r.URL.Query().Get("id")
	if schoolClassId == "" {
		common.HttpErrorResponse(w, "no school class id specified")
		return
	}
	schoolClass, err := c.repository.GetById(r.Context(), schoolClassId)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.JsonResponse(w, schoolClass)
}

func (c SchoolClassController) AddBook(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	if claims.SchoolId == nil {
		common.HttpErrorResponse(w, "user is not assigned to a school")
		return
	}
	var schoolClass SchoolClass
	schoolClass.ID = uuid.NewString()
	if err := json.NewDecoder(r.Body).Decode(&schoolClass); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	if err := c.repository.Insert(r.Context(), schoolClass, *claims.SchoolId, claims.UserId); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}

func (c SchoolClassController) UpdateBook(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	if claims.SchoolId == nil {
		common.HttpErrorResponse(w, "user is not assigned to a school")
		return
	}
	var schoolClass SchoolClass
	if err := json.NewDecoder(r.Body).Decode(&schoolClass); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	if err := c.repository.Update(r.Context(), schoolClass, claims.UserId); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}

func (c SchoolClassController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookId := r.URL.Query().Get("id")
	if bookId == "" {
		common.HttpErrorResponse(w, "no school class id specified")
		return
	}
	if err := c.repository.Delete(r.Context(), bookId); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}
