package common

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kammeph/school-book-storage-service-simplified/fp"
)

var (
	jwtSecretKey             = os.Getenv("JWT_SECRET_KEY")
	jwtAccessTokenExpiry, _  = strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_EXPIRY_SEC"))
	jwtRefreshTokenExpiry, _ = strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_EXPIRY_SEC"))
)

func IsAllowed(handler func(w http.ResponseWriter, r *http.Request), roles []Role) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := getAccessToken(r)
		if err != nil {
			HttpErrorResponseWithStatusCode(w, err.Error(), http.StatusUnauthorized)
			return
		}
		claims := &AccessClaims{}
		if err := GetClaimsFromToken(r, tokenString, claims); err != nil {
			HttpErrorResponseWithStatusCode(w, err.Error(), http.StatusUnauthorized)
			return
		}
		for _, role := range roles {
			if fp.Some(claims.Roles, func(r Role) bool { return r == role }) {
				handler(w, r)
				return
			}
		}
		HttpErrorResponseWithStatusCode(w, "user missing permissions", http.StatusUnauthorized)
	}
}

func IsAllowedWithClaims(
	handler func(w http.ResponseWriter, r *http.Request, claims AccessClaims),
	roles []Role,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := getAccessToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		claims := &AccessClaims{}
		if err := GetClaimsFromToken(r, tokenString, claims); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if HasRoles(claims.Roles, roles) {
			handler(w, r, *claims)
			return
		}
		http.Error(w, "user missing permissions", http.StatusForbidden)
	}
}

func getAccessToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", errors.New("access token is not set")
	}
	if !strings.ContainsAny(auth, "Bearer") {
		return "", errors.New("no bearer token found")
	}
	token := strings.Split(auth, " ")[1]
	return token, nil
}

func GetRefreshToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func GetClaimsFromToken(r *http.Request, tokenString string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("access token is invalid")
	}
	return nil
}

func CreateAccessToken(userId, username string, schoolId *string, roles []Role, locale Locale) (string, error) {
	expirationTime := time.Now().Add(time.Duration(jwtAccessTokenExpiry) * time.Second)
	claims := AccessClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    username,
		},
		UserId: userId, SchoolId: schoolId, Username: username, Roles: roles, Locale: locale,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecretKey))
}

func CreateRefreshToken(userID string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(jwtRefreshTokenExpiry) * time.Second)
	claims := RefreshClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		UserID: userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecretKey))
}
