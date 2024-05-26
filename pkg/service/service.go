package service

import (
	todo "todolist"
	"todolist/pkg/repository"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username, password string) (string, error)
}

type ToDoList interface {
}

type ToDoItem interface {
}

type Service struct {
	Authorization
	ToDoList
	ToDoItem
}

func NewService(repos *repository.Repository) *Service { //
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
