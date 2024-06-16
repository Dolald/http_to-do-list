package repository

import (
	"errors"
	"testing"
	todo "todolist"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func Test_TodoItemPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// создаем обязательное закрытие нашей моковой db
	defer db.Close()

	// создаем структуру todoPostgres с доступом к db
	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		item   todo.TodoItem
	}

	type mockBehavior func(args args, id int)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		want    int
		wantErr bool
	}{
		{
			name: "ok",
			input: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "test title",
					Description: "test description",
				},
			},
			want: 2,
			mock: func(args args, id int) {
				mock.ExpectBegin()
				// пишем какой ответ хотим получить
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)

				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").WithArgs(args.listId, id).WillReturnResult(sqlmock.NewResult(1, 1))
				// 1 - id вставленной строчки, 1 - количество затронутых строчек командой
				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			input: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "",
					Description: "description",
				},
			},
			mock: func(args args, id int) {
				mock.ExpectBegin()
				// хотим получить ошибку
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(0, errors.New("insert error"))
				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)
				// ожидаем откат транзакции
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Failed 2nd Insert",
			input: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "title",
					Description: "description",
				},
			},
			mock: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").WithArgs(args.listId, id).WillReturnError(errors.New("insert error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input, tt.want)

			got, err := r.Create(tt.input.listId, tt.input.item)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoItemPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		userId int
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []todo.TodoItem
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(1, "title1", "description1", true).
					AddRow(2, "title2", "description2", false).
					AddRow(3, "title3", "description3", false)

				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li on (.+) INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1, 1).WillReturnRows(rows)
			},
			input: args{
				listId: 1,
				userId: 1,
			},
			want: []todo.TodoItem{
				{1, "title1", "description1", true},
				{2, "title2", "description2", false},
				{3, "title3", "description3", false},
			},
		},
	}
}
