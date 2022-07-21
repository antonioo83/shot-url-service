package interfaces

import (
	"github.com/antonioo83/shot-url-service/internal/models"
	"net/http"
)

type UserAuthHandler interface {
	GetAuthUser(r *http.Request, w http.ResponseWriter) (models.User, error)
	SetToken(w http.ResponseWriter, token string)
	GetToken(r *http.Request) (string, error)
	GetUserCode(token string) (int, error)
	GenerateToken(code int) (string, error)
	ValidateToken(token string) (bool, error)
}
