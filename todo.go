package todo

import "errors"

type TodoList struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type UsersList struct { // связная таблица с 2 таблицами User и TodoList
	Id     int
	UserId int
	ListId int
}

type TodoItem struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done" db:"done"`
}

type ListsItems struct { // связная таблица с 2 таблицами TodoList и TodoItem
	Id     int
	ListId int // TodoList.id
	ItemId int // TodoItem.id
}

type UpdateListInput struct { // создаём экземпляр структуры из полученного запроса
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (i UpdateListInput) Validate() error {
	if i.Description == nil && i.Title == nil {
		return errors.New("update structures has no values")
	}
	return nil
}

type UpdateItemInput struct { // создаём экземпляр структуры из полученного запроса
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

func (i UpdateItemInput) Validate() error {
	if i.Description == nil && i.Title == nil && i.Done == nil {
		return errors.New("update structures has no values")
	}
	return nil
}
