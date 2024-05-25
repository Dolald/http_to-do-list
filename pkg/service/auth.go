package service

import (
	"crypto/sha1"
	"fmt"
	todo "todolist"
	"todolist/pkg/repository"
)

const salt = "asdkoawk2jdkji"

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

func (a *AuthService) GenerateToken(username, password string) (string, error) {

	return "", nil
}

func gereratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
