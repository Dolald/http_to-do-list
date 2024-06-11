package repository

import (
	"fmt"
	todo "todolist"

	"github.com/jmoiron/sqlx"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, item todo.TodoItem) (int, error) {
	transaction, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int

	createItemQuery := fmt.Sprintf(`INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id`, todoItemsTable) // засовываем в таблицу данные
	row := transaction.QueryRow(createItemQuery, item.Title, item.Description)                                         // QueryRow - возвращает строку

	err = row.Scan(&itemId) // принимаем id item'а
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	createListsItemsQuery := fmt.Sprintf(`INSERT INTO %s (list_id, item_id) VALUES ($1, $2)`, listsItemsTable)
	_, err = transaction.Exec(createListsItemsQuery, listId, itemId) // exec - просто выполняем запрос без возврата чего-либо
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	return itemId, transaction.Commit()
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem

	query := fmt.Sprintf(`
	SELECT ti.id, ti.title, ti.description, ti.done 
	FROM "%s" ti
	JOIN "%s" li ON ti.id = li.item_id
	JOIN "%s" ul ON ul.list_id = li.list_id
	WHERE ul.user_id = $1 AND li.list_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Select(&items, query, userId, listId); err != nil { // Select засовывает в переменную все нужные данные
		return nil, err
	}

	return items, nil
}

func (r *TodoItemPostgres) GetById(userId, itemId int) (todo.TodoItem, error) {
	var todoItem todo.TodoItem

	query := fmt.Sprintf(`
	SELECT ti.id, ti.title, ti.description, ti.done
	FROM %s ti
	JOIN %s li ON li.item_id = ti.id
	JOIN %s ul ON ul.list_id = li.list_id
	WHERE ul.user_id = $1 AND ti.id = $2`, todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Get(&todoItem, query, userId, itemId); err != nil { // get засовывает в переменную все нужные данные и возвращает только 1 значение
		return todoItem, err
	}

	return todoItem, nil

}

func (r *TodoItemPostgres) DeleteItem(userId, itemId int) error {
	query := fmt.Sprintf(`
	DELETE FROM %s ti
	USING %s li, %s ul
	WHERE ti.id = li.item_id AND ul.list_id = li.list_id AND ul.user_id = $1 AND ul.list_id = $2`, todoItemsTable, listsItemsTable, usersListsTable)

	_, err := r.db.Exec(query, userId, itemId) // Exec - выполняет запрос, не возвращая никаних значений

	return err
}
