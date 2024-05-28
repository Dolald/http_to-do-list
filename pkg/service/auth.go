package service

import (
	"crypto/sha1"
	"fmt"
	"time"
	todo "todolist"
	"todolist/pkg/repository"

	"github.com/dgrijalva/jwt-go"
)

const (
	salt        = "asdkoawk2jdkji"
	tokenTTL    = 12 * time.Hour
	signgingKey = "sa2000dddwli29d2kpld"
)

type tokenClaims struct { // создаём структуру токена
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (a *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = gereratePasswordHash(user.Password)

	return a.repo.CreateUser(user)
}

func gereratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (a *AuthService) GenerateToken(username, password string) (string, error) { // получаем токен пользователя из базы
	userId, err := a.repo.GetUser(username, gereratePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId.Id,
	})

	return token.SignedString([]byte(signgingKey))
}

func (a *AuthService) ParseToken(token string) (int, error) {
	jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {

	})
	return 0, nil
}
