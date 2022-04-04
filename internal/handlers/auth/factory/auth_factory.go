package factory

import (
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/handlers/auth/cookie"
	"github.com/antonioo83/shot-url-service/internal/handlers/auth/interfaces"
	repositoryInterfaces "github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/go-chi/jwtauth"
)

func NewAuthHandler(tokenAuth *jwtauth.JWTAuth, userRepository repositoryInterfaces.UserRepository, config config.Config) interfaces.UserAuthHandler {
	return cookie.NewUserAuthHandler(tokenAuth, userRepository, config)
}
