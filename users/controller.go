package users

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/kammeph/school-book-storage-service-simplified/common"
	"github.com/kammeph/school-book-storage-service-simplified/fp"
)

type UsersController struct {
	usersRepository UsersRepository
}

func NewUsersController(db *sql.DB) UsersController {
	return UsersController{NewSqlUserRepository(db)}
}

func AddUsersController(db *sql.DB) {
	controller := NewUsersController(db)
	common.Get("/api/users/me",
		common.IsAllowedWithClaims(controller.GetMe, []common.Role{common.User, common.Superuser, common.Admin, common.SysAdmin}))
	common.Get("/api/users",
		common.IsAllowed(controller.GetUsers, []common.Role{common.SysAdmin}))
	common.Get("/api/users/by-id",
		common.IsAllowedWithClaims(controller.GetUserById, []common.Role{common.User, common.Superuser, common.Admin, common.SysAdmin}))
}

func (c UsersController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.usersRepository.GetUsers(context.Background())
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpResponse(w, users)
}

func (c UsersController) GetUserById(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	userId := r.URL.Query().Get("userId")
	if userId != claims.UserId && !fp.Some(claims.Roles, func(r common.Role) bool { return r == common.SysAdmin }) {
		common.HttpErrorResponseWithStatusCode(w, "user missing permissions", http.StatusForbidden)
		return
	}
	user, err := c.usersRepository.GetUserById(context.Background(), userId)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpResponse(w, user)
}

func (c UsersController) GetMe(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	user, err := c.usersRepository.GetUserById(context.Background(), claims.UserId)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpResponse(w, user)
}
