package cookie

import (
	"fmt"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/go-chi/jwtauth"
	"net/http"
	"time"
)

var tokenAuth *jwtauth.JWTAuth

func SetToken(w http.ResponseWriter, user models.User) {
	cookie := http.Cookie{Name: "token", Value: user.UID, Expires: time.Now().Add(60 * 30 * time.Second)}
	http.SetCookie(w, &cookie)
}

func GetToken(r *http.Request) string {
	token, err := r.Cookie("token")
	if err != nil {
		return ""
	}

	return token.Value
}

func GetUserCode(token string) int {
	data := getExtractedData(token)
	userCode, ok := data["user_code"]
	if !ok {
		return 0
	}

	var val = userCode.(float64)
	return int(val)
}

func GenerateToken(code int) (string, error) {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_code": code})
	if err != nil {
		return "", fmt.Errorf("i can't generate token: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(token string) bool {
	user := getExtractedData(token)
	_, ok := user["user_code"]

	return ok
}

func getExtractedData(token string) map[string]interface{} {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	jwtToken, err := tokenAuth.Decode(token)
	if err != nil {
		return nil
	}

	return jwtToken.PrivateClaims()
}
