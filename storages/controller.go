package storages

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service-simplified/common"
)

type StoragesResponseModel struct {
	Storages []Storage `json:"storages"`
}

func StoragesResponse(w http.ResponseWriter, storagees []Storage) {
	common.JsonResponse(w, StoragesResponseModel{storagees})
}

type StorageResponseModel struct {
	Storage Storage `json:"storage"`
}

func StorageResponse(w http.ResponseWriter, storage Storage) {
	common.JsonResponse(w, StorageResponseModel{storage})
}

type StoragesController struct {
	repository StoragesRepository
}

func NewBooksController(db *sql.DB) StoragesController {
	return StoragesController{NewSqlSchoolsRepository(db)}
}

func AddStoragesController(db *sql.DB) {
	controller := NewBooksController(db)
	common.Get("/api/storages/get-all", common.IsAllowedWithClaims(controller.GetAllStoragees, []common.Role{common.User, common.Superuser, common.Admin, common.SysAdmin}))
	common.Get("/api/storages/get-by-id", common.IsAllowed(controller.GetStorageById, []common.Role{common.User, common.Superuser, common.Admin, common.SysAdmin}))
	common.Post("/api/storages/add", common.IsAllowedWithClaims(controller.AddStorage, []common.Role{common.Superuser, common.Admin, common.SysAdmin}))
	common.Post("/api/storages/update", common.IsAllowedWithClaims(controller.UpdateStorage, []common.Role{common.Superuser, common.Admin, common.SysAdmin}))
	common.Post("/api/storages/delete", common.IsAllowed(controller.DeleteStorage, []common.Role{common.Admin, common.SysAdmin}))
}

func (c StoragesController) GetAllStoragees(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	if claims.SchoolId == nil {
		common.HttpErrorResponse(w, "user is not assigned to a school")
		return
	}
	storagees, err := c.repository.GetAll(r.Context(), *claims.SchoolId)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	StoragesResponse(w, storagees)
}

func (c StoragesController) GetStorageById(w http.ResponseWriter, r *http.Request) {
	storageId := r.URL.Query().Get("id")
	if storageId == "" {
		common.HttpErrorResponse(w, "no storage id specified")
		return
	}
	storage, err := c.repository.GetById(r.Context(), storageId)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	StorageResponse(w, storage)
}

func (c StoragesController) AddStorage(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	if claims.SchoolId == nil {
		common.HttpErrorResponse(w, "user is not assigned to a school")
		return
	}
	var storage Storage
	storage.ID = uuid.NewString()
	if err := json.NewDecoder(r.Body).Decode(&storage); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	if err := c.repository.Insert(r.Context(), storage, *claims.SchoolId, claims.UserId); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}

func (c StoragesController) UpdateStorage(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	if claims.SchoolId == nil {
		common.HttpErrorResponse(w, "user is not assigned to a school")
		return
	}
	var storage Storage
	if err := json.NewDecoder(r.Body).Decode(&storage); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	if err := c.repository.Update(r.Context(), storage, claims.UserId); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}

func (c StoragesController) DeleteStorage(w http.ResponseWriter, r *http.Request) {
	storageId := r.URL.Query().Get("id")
	if storageId == "" {
		common.HttpErrorResponse(w, "no storage id specified")
		return
	}
	if err := c.repository.Delete(r.Context(), storageId); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}
