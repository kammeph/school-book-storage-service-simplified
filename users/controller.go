package users

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/kammeph/school-book-storage-service-simplified/common"
	"github.com/kammeph/school-book-storage-service-simplified/fp"
	"golang.org/x/crypto/bcrypt"
)

type UsersResponseModel struct {
	Users []UserDto `json:"users"`
}

func UsersResponse(w http.ResponseWriter, users []UserDto) {
	response := UsersResponseModel{users}
	common.JsonResponse(w, response)
}

type UserResponseModel struct {
	User UserDto `json:"user"`
}

func UserResponse(w http.ResponseWriter, user UserDto) {
	response := UserResponseModel{user}
	common.JsonResponse(w, response)
}

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
	common.Get("/api/users/get-all",
		common.IsAllowed(controller.GetUsers, []common.Role{common.SysAdmin}))
	common.Get("/api/users/get-by-id",
		common.IsAllowedWithClaims(controller.GetUserById, []common.Role{common.User, common.Superuser, common.Admin, common.SysAdmin}))
	common.Post("/api/users/update",
		common.IsAllowedWithClaims(controller.UpdateUser, []common.Role{common.User, common.Superuser, common.Admin, common.SysAdmin}))
	common.Post("/api/users/delete",
		common.IsAllowed(controller.DeleteUser, []common.Role{common.SysAdmin}))
	common.Post("/api/users/change-password",
		common.IsAllowedWithClaims(controller.ChangePassword, []common.Role{common.User, common.Superuser, common.Admin, common.SysAdmin}))
}

func (c UsersController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.usersRepository.GetAll(r.Context())
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	UsersResponse(w, users)
}

func (c UsersController) GetUserById(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	userId := r.URL.Query().Get("userId")
	if userId != claims.UserId && !fp.Some(claims.Roles, func(r common.Role) bool { return r == common.SysAdmin }) {
		common.HttpErrorResponseWithStatusCode(w, "user missing permissions", http.StatusForbidden)
		return
	}
	user, err := c.usersRepository.GetById(r.Context(), userId)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	UserResponse(w, user)
}

func (c UsersController) GetMe(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	user, err := c.usersRepository.GetById(r.Context(), claims.UserId)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	UserResponse(w, user)
}

func (c UsersController) UpdateUser(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	var user UserDto
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	if user.ID != claims.UserId && !fp.Some(claims.Roles, func(r common.Role) bool { return r == common.SysAdmin }) {
		common.HttpErrorResponseWithStatusCode(w, "user missing permissions", http.StatusForbidden)
		return
	}
	if err := c.usersRepository.Update(r.Context(), user); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}

func (c UsersController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")
	if userId == "" {
		common.HttpErrorResponse(w, "no user id specified")
	}
	if err := c.usersRepository.Delete(r.Context(), userId); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}

func (c UsersController) ChangePassword(w http.ResponseWriter, r *http.Request, claims common.AccessClaims) {
	var passwordUpdate PasswordUpdate
	if err := json.NewDecoder(r.Body).Decode(&passwordUpdate); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	if passwordUpdate.OldPassword == passwordUpdate.NewPassword {
		common.HttpErrorResponse(w, "old and new password must be unequal")
		return
	}
	oldPasswordHash, err := c.usersRepository.GetCredentialsById(r.Context(), claims.UserId)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(oldPasswordHash), []byte(passwordUpdate.OldPassword)); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(passwordUpdate.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	if err := c.usersRepository.UpdatePassword(r.Context(), claims.UserId, string(newPasswordHash)); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	common.HttpSuccessResponse(w)
}
