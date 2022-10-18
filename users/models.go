package users

import (
	"github.com/google/uuid"
	"github.com/kammeph/school-book-storage-service-simplified/common"
)

type UserModel struct {
	ID           string
	SchoolId     *string
	Username     string
	PasswordHash string
	Active       bool
	Roles        []common.Role
	Locale       common.Locale
}

type UserDto struct {
	ID       string        `json:"id"`
	SchoolId *string       `json:"schoolId"`
	Username string        `json:"username"`
	Roles    []common.Role `json:"roles"`
	Locale   common.Locale `json:"locale"`
}

type PasswordUpdate struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UserRole struct {
	UserID string
	Role   common.Role
}

func NewUser(schoolId *string, username, passwordHash string, roles []common.Role, locale common.Locale) UserModel {
	return UserModel{
		ID:           uuid.NewString(),
		SchoolId:     schoolId,
		Username:     username,
		PasswordHash: passwordHash,
		Active:       true,
		Roles:        roles,
		Locale:       locale,
	}
}

func NewUserWithDefaultRole(username, passwordHash string) UserModel {
	return UserModel{
		ID:           uuid.NewString(),
		SchoolId:     nil,
		Username:     username,
		PasswordHash: passwordHash,
		Active:       true,
		Roles:        []common.Role{common.User},
		Locale:       common.DE,
	}
}
