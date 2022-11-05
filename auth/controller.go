package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/kammeph/school-book-storage-service-simplified/common"
	"github.com/kammeph/school-book-storage-service-simplified/users"
	"golang.org/x/crypto/bcrypt"
)

type AccessTokenResponseModel struct {
	AccessToken string `json:"accessToken"`
}

func AccessTokenResponse(w http.ResponseWriter, accessToken string) {
	common.JsonResponse(w, AccessTokenResponseModel{accessToken})
}

type AuthController struct {
	usersRepo users.UsersRepository
}

func NewAuthController(db *sql.DB) AuthController {
	return AuthController{users.NewSqlUserRepository(db)}
}

func AddAuthController(db *sql.DB) {
	controller := NewAuthController(db)
	common.Post("/api/auth/login", controller.Login)
	common.Post("/api/auth/logout", controller.Logout)
	common.Post("/api/auth/register", controller.Register)
	common.Get("/api/auth/refresh", controller.Refresh)
}

func (c AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	ctx := r.Context()
	passwordHash, err := c.usersRepo.GetCredentialsByName(ctx, credentials.Username)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(credentials.Password)); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	user, err := c.usersRepo.GetByName(ctx, credentials.Username)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	accessToken, err := common.CreateAccessToken(user.ID, user.Username, user.SchoolId, user.Roles, user.Locale)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	refreshToken, err := common.CreateRefreshToken(user.ID)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)
	AccessTokenResponse(w, accessToken)
}

func (c AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)
	common.HttpSuccessResponse(w)
}

func (c AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		common.HttpErrorResponse(w, err.Error())
		return
	}
	user := users.NewUserWithDefaultRole(credentials.Username, string(passwordHash))
	if err := c.usersRepo.Insert(r.Context(), user); err != nil {
		common.HttpErrorResponse(w, err.Error())
	}
	common.HttpSuccessResponse(w)
}

func (c AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	tokenString, err := common.GetRefreshToken(r)
	if err != nil {
		common.HttpErrorResponseWithStatusCode(w, err.Error(), http.StatusUnauthorized)
		return
	}
	claims := &common.RefreshClaims{}
	if err := common.GetClaimsFromToken(r, tokenString, claims); err != nil {
		common.HttpErrorResponseWithStatusCode(w, err.Error(), http.StatusUnauthorized)
		return
	}
	user, err := c.usersRepo.GetById(r.Context(), claims.UserID)
	if err != nil {
		common.HttpErrorResponseWithStatusCode(w, err.Error(), http.StatusUnauthorized)
		return
	}
	accessToken, err := common.CreateAccessToken(user.ID, user.Username, user.SchoolId, user.Roles, user.Locale)
	if err != nil {
		common.HttpErrorResponseWithStatusCode(w, err.Error(), http.StatusUnauthorized)
		return
	}
	AccessTokenResponse(w, accessToken)
}
