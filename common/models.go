package common

import (
	"github.com/golang-jwt/jwt"
	"github.com/kammeph/school-book-storage-service-simplified/fp"
)

type Role string

const (
	SysAdmin  Role = "SYS_ADMIN"
	Admin     Role = "ADMIN"
	Superuser Role = "SUPERUSER"
	User      Role = "USER"
)

func HasRoles(userRoles, allowedRoles []Role) bool {
	for _, role := range allowedRoles {
		if fp.Some(userRoles, func(r Role) bool { return r == role }) {
			return true
		}
	}
	return false
}

func IsSysAdmin(userRoles []Role) bool {
	for _, role := range userRoles {
		if role == SysAdmin {
			return true
		}
	}
	return false
}

type Locale string

const (
	DE Locale = "DE"
	EN Locale = "EN"
)

type AccessClaims struct {
	jwt.StandardClaims
	UserId   string  `json:"userId"`
	SchoolId *string `json:"schoolId"`
	Username string  `json:"userName"`
	Roles    []Role  `json:"roles"`
	Locale   Locale  `json:"locale"`
}

type RefreshClaims struct {
	jwt.StandardClaims
	UserID string `json:"userId"`
}
