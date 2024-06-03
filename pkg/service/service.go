package service

import (
	todo "todolist"
	"todolist/pkg/repository"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
	UpdateList()
}

type ToDoList interface {
	Create(userId int, list todo.TodoList) (int, error)
	GetAllLists(userId int) ([]todo.TodoList, error)
	GetById(userId, listId int) (todo.TodoList, error)
	UpdateList(userId, listId, list todo.UpdateListInput) error
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
		ToDoList:      NewToDoListService(repos.ToDoList),
	}
}
