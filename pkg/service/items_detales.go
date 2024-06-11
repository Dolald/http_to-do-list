package service

import (
	todo "todolist"
	"todolist/pkg/repository"
)

type TodoItemService struct { // здесь находятся на перепутье 2 интерфейса
	repo     repository.ToDoItem // interface Item
	listRepo repository.ToDoList // interface List
}

func NewTodoItemService(repo repository.ToDoItem, listRepo repository.ToDoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemService) Create(userId, listId int, item todo.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId) // проверяем существует ли list с таким id или нет
	if err != nil {
		return 0, err
	}
	return s.repo.Create(listId, item) // переходим на нижний уровень, уровень репозитория
}

func (s *TodoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *TodoItemService) GetById(userId, itemId int) (todo.TodoItem, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *TodoItemService) DeleteItem(userId, itemId int) error {
	return s.repo.DeleteItem(userId, itemId)
}
