package common

import "github.com/golang-jwt/jwt"

type Role string

const (
	SysAdmin  Role = "SYS_ADMIN"
	Admin     Role = "ADMIN"
	Superuser Role = "SUPERUSER"
	User      Role = "USER"
)

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
