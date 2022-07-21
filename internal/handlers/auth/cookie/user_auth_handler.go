package cookie

import (
	"fmt"
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/handlers/auth/interfaces"
	"github.com/antonioo83/shot-url-service/internal/models"
	repositoryInterfaces "github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/go-chi/jwtauth"
	"net/http"
	"time"
)

type userAuthHandler struct {
	tokenAuth      *jwtauth.JWTAuth
	userRepository repositoryInterfaces.UserRepository
	config         config.Config
}

// NewUserAuthHandler create userAuthHandler instance.
func NewUserAuthHandler(tokenAuth *jwtauth.JWTAuth,
	userRepository repositoryInterfaces.UserRepository, config config.Config) interfaces.UserAuthHandler {
	return &userAuthHandler{tokenAuth, userRepository, config}
}

// GetAuthUser get authorized user.
func (a userAuthHandler) GetAuthUser(r *http.Request, w http.ResponseWriter) (models.User, error) {
	var user models.User
	var lastModel *models.User
	var err error
	token, _ := a.GetToken(r)
	isValidate, _ := a.ValidateToken(token)
	if token == "" || !isValidate {
		lastModel, err = a.userRepository.GetLastModel()
		if err != nil {
			return user, fmt.Errorf("i can't get auth user: %w", err)
		}
		if lastModel.Code == 0 {
			user.Code = 1
			user.UID, err = a.GenerateToken(user.Code)
		} else {
			user = *lastModel
			user.Code = user.Code + 1
			user.UID, err = a.GenerateToken(user.Code)
		}
		a.SetToken(w, user.UID)
	} else {
		user.UID = token
		user.Code, err = a.GetUserCode(token)
	}

	if err != nil {
		return user, fmt.Errorf("i can't get auth user: %w", err)
	}

	return user, nil
}

// SetToken set token to a cookie.
func (a userAuthHandler) SetToken(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:    a.config.Auth.TokenName,
		Value:   token,
		Expires: time.Now().Add(a.config.Auth.CookieTTL),
	}
	http.SetCookie(w, &cookie)
}

// GetToken get token from a cookie.
func (a userAuthHandler) GetToken(r *http.Request) (string, error) {
	token, err := r.Cookie(a.config.Auth.TokenName)
	if err != nil {
		return "", err
	}

	return token.Value, err
}

// GetUserCode get user code.
func (a userAuthHandler) GetUserCode(token string) (int, error) {
	data, err := a.getExtractedData(token)
	userCode, ok := data["user_code"]
	if !ok {
		return 0, fmt.Errorf("can't get user code: %v", err)
	}

	var val, ok2 = userCode.(float64)
	if !ok2 {
		return 0, fmt.Errorf("can't reduce user code: %s", userCode)
	}

	return int(val), nil
}

// GenerateToken generate token.
func (a userAuthHandler) GenerateToken(code int) (string, error) {
	tokenAuth := jwtauth.New(
		a.config.Auth.Alg,
		a.config.Auth.SignKey,
		nil,
	)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_code": code})
	if err != nil {
		return "", fmt.Errorf("i can't generate token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validate token.
func (a userAuthHandler) ValidateToken(token string) (bool, error) {
	user, err := a.getExtractedData(token)
	if err != nil {
		return false, fmt.Errorf("can't decode token from a cookie: %v", err)
	}
	_, ok := user["user_code"]

	return ok, nil
}

func (a userAuthHandler) getExtractedData(token string) (map[string]interface{}, error) {
	tokenAuth := jwtauth.New(
		a.config.Auth.Alg,
		a.config.Auth.SignKey,
		nil,
	)
	jwtToken, err := tokenAuth.Decode(token)
	if err != nil {
		return nil, fmt.Errorf("can't decode token from a cookie: %v", err)
	}

	return jwtToken.PrivateClaims(), nil
}
