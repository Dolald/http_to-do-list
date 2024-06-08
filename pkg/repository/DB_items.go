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

	createItemQuery := fmt.Sprintf(`INSERT %s (title, description) VALUES $1, $2 RETURNING id`, todoItemsTable) // засовываем в таблицу данные
	row := transaction.QueryRow(createItemQuery, item.Title, item.Description)                                  // QueryRow - возвращает строку

	err = row.Scan(&itemId) // принимаем id таблицы с кайфом
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	createListsItemsQuery := fmt.Sprintf(`INSERT %s (list_id, item_id) VALUES $1, $2`, listsItemsTable)
	_, err = transaction.Exec(createListsItemsQuery, listId, itemId) // exec - просто выполняем запрос без возврата чего-либо
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	// JOIN %s ON todo_item.id = lists_item.id
	// JOIN %s ON lists_item.id = todo_lists.id
	// WHERE user_`, todoItemsTable, listsItemsTable, todoListsTable)

	return 0, nil
}
