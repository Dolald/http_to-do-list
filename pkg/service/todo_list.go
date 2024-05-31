package service

import (
	todo "todolist"
	"todolist/pkg/repository"
)

type ToDoListService struct {
	repo repository.ToDoList
}

func NewToDoListService(repo repository.ToDoList) *ToDoListService {
	return &ToDoListService{repo: repo}
}

func (t *ToDoListService) Create(userId int, list todo.TodoList) (int, error) {
	return t.repo.Create(userId, list)
}

func (t *ToDoListService) GetAllLists(userId int) ([]todo.TodoList, error) {
	return t.repo.GetAllLists(userId)
}