package repository

import (
	"fmt"
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
	var lists []todo.TodoList

	getAllListsQuery := fmt.Sprintf("SELECT * FROM %s WHERE $1", todoListsTable)

	row, err := t.db.Exec(getAllListsQuery, userId)
	if err != nil {
		return lists, err
	}

	if err := row.Scan(&lists); err != nil {
		return lists, err
	}

	return lists, err
}
