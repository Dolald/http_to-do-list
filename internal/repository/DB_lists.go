package repository

import (
	"fmt"
	"strings"
	todo "todolist"

	"github.com/jmoiron/sqlx"
)

type ToDoListPostgres struct {
	db *sqlx.DB
}

func NewToDoListPostgres(db *sqlx.DB) *ToDoListPostgres {
	return &ToDoListPostgres{db: db}
}

func (t *ToDoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	transaction, err := t.db.Begin() // создаём транзакцию
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := transaction.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		transaction.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = transaction.Exec(createUsersListQuery, userId, id)
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	return id, transaction.Commit()
}

func (t *ToDoListPostgres) GetAllLists(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList // создаём экземпляр структуры из наших тасков

	getAllListsQuery := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl JOIN %s ul ON (tl.id = ul.list_id) WHERE ul.user_id = $1", todoListsTable, usersListsTable)

	err := t.db.Select(&lists, getAllListsQuery, userId) // выполняем операцию и засовываем в list

	return lists, err
}

func (t *ToDoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList

	getOneList := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl JOIN %s ul ON (tl.id = ul.list_id) WHERE ul.user_id = $1 AND ul.list_id = $2", todoListsTable, usersListsTable)

	err := t.db.Get(&list, getOneList, userId, listId)

	return list, err
}

func (t *ToDoListPostgres) DeleteList(userId, listId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND $1 = ul.user_id AND ul.list_id = $2", todoListsTable, usersListsTable)

	_, err := t.db.Exec(query, userId, listId) // Exec - выполняет запрос, не возвращая никаних значений

	return err
}

func (t *ToDoListPostgres) UpdateList(userId, listId int, list todo.UpdateListInput) error {
	values := make([]string, 0)
	args := make([]any, 0)
	argsId := 1

	if list.Title != nil { // проверяем что в запросе пользователя изменено
		values = append(values, fmt.Sprintf("title=$%d", argsId)) //"title=$1"
		args = append(args, *list.Title)
		argsId++
	}

	if list.Description != nil {
		values = append(values, fmt.Sprintf("description=$%d", argsId)) //"description=$2"
		args = append(args, *list.Description)
		argsId++
	}

	setQuery := strings.Join(values, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id = $%d AND ul.user_id = $%d",
		todoListsTable, setQuery, usersListsTable, argsId, argsId+1)

	args = append(args, listId, userId)

	_, err := t.db.Exec(query, args...)

	return err
}
