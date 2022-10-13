package schools

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/kammeph/school-book-storage-service-simplified/common"
)

type SchoolsResponseModel struct {
	Schools []SchoolDto `json:"schools"`
}

func SchoolsResponse(w http.ResponseWriter, schools []SchoolDto) {
	common.JsonResponse(w, SchoolsResponseModel{schools})
}

type SchoolResponseModel struct {
	School SchoolDto `json:"school"`
}

func SchoolResponse(w http.ResponseWriter, school SchoolDto) {
	common.JsonResponse(w, SchoolResponseModel{school})
}

type SchoolsController struct {
	repository SchoolsRepository
}

func NewSchoolsController(db *sql.DB) SchoolsController {
	return SchoolsController{NewSqlSchoolsRepository(db)}
}

func AddSchoolsController(db *sql.DB) {
	controller := NewSchoolsController(db)
	common.Get("/api/schools", common.IsAllowed(controller.GetAll, []common.Role{common.SysAdmin}))
	common.Get("/api/schools/by-id", common.IsAllowedWithClaims(controller.GetById, []common.Role{common.User, common.Superuser, common.Admin, common.SysAdmin}))
	common.Post("/api/schools/add", common.IsAllowedWithClaims(controller.AddSchool, []common.Role{common.SysAdmin}))
	common.Post("/api/schools/update", common.IsAllowedWithClaims(controller.UpdateSchool, []common.Role{common.SysAdmin, common.Admin}))
	common.Post("/api/schools/delete", common.IsAllowedWithClaims(controller.DeleteSchool, []common.Role{common.SysAdmin}))
}

func (c SchoolsController) GetAll(w http.ResponseWriter, r *http.Request) {
	schools, err := c.repository.GetAll(r.Context())
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	SchoolsResponse(w, schools)
}

func (c SchoolsController) GetById(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	schoolId := r.URL.Query().Get("id")
	if schoolId == "" {
		common.HttpErrorResponse(w, "no school id specified")
		return
	}
	if (claims.SchoolId == nil || *claims.SchoolId != schoolId) && !common.IsSysAdmin(claims.Roles) {
		common.HttpErrorResponseWithStatusCode(w, "user missing permissions", http.StatusForbidden)
		return
	}
	school, err := c.repository.GetById(r.Context(), schoolId)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	SchoolResponse(w, school)
}

func (c SchoolsController) AddSchool(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	var addSchool AddSchoolDto
	if err := json.NewDecoder(r.Body).Decode(&addSchool); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	school := NewSchool(addSchool.Name, claims.UserId)
	if err := c.repository.Insert(r.Context(), school); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}

func (c SchoolsController) UpdateSchool(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	var updateSchool UpdateSchoolDto
	if err := json.NewDecoder(r.Body).Decode(&updateSchool); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	school := NewSchoolWithId(updateSchool.ID, updateSchool.Name, claims.UserId)
	if err := c.repository.Update(r.Context(), school); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}

func (c SchoolsController) DeleteSchool(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	schoolId := r.URL.Query().Get("id")
	if schoolId == "" {
		common.HttpErrorResponse(w, "no school id specified")
		return
	}
	if (claims.SchoolId == nil || *claims.SchoolId != schoolId) && !common.IsSysAdmin(claims.Roles) {
		common.HttpErrorResponseWithStatusCode(w, "user missing permissions", http.StatusForbidden)
		return
	}
	if err := c.repository.Delete(r.Context(), schoolId, claims.UserId); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}
