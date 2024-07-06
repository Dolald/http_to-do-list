package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"
	todo "todolist/internal/domain"
	"todolist/internal/repository"

	"github.com/dgrijalva/jwt-go"
)

const (
	salt        = "asdkoawk2jdkji"
	tokenTTL    = 99 * time.Hour
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
	hash := sha1.New()           // создаём новый хеш функцию
	hash.Write([]byte(password)) // прокидываем  password через хеш функцию

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (a *AuthService) GenerateToken(username, password string) (string, error) { // создаём токен для пользователя
	userId, err := a.repo.GetUser(username, gereratePasswordHash(password)) // получили ID пользователя
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{ // jwt.NewWithClaims - создаём новый токен, jwt.SigningMethodHS256 - алгоритм подписи, пое*ень-трава, которую мы раньше создали
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(), // через 99 часов токен превращается в тыкву
			IssuedAt:  time.Now().Unix(),               // дата происзводства сейчас
		},
		userId.Id, // засовывает полученный ранее ID пользователя
	})

	return token.SignedString([]byte(signgingKey)) // подписывает токен с использованием секретного ключа и превращаем его в строку
}

func (a *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (any, error) { // Функция jwt.ParseWithClaims разбирает токен, разделяя его на header, payload и signature, заполняем tokenClaims signgingKey
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // token.Method - метод подписи, если токен с методом SigningMethodHMAC есть, то всё норм
			return nil, errors.New("invalid signing method")
		}
		return []byte(signgingKey), nil // возвращаем секретный ключ
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims) // устанавливаем тип tokenClaims в JWT токен
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}
