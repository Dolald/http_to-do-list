package repository

import (
	"testing"
	todo "todolist/internal/domain"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	// создаем мок для sqlx
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// создаем обязательное закрытие нашей моковой db
	defer db.Close()

	// создаем структуру authPostgres с доступом к db
	r := NewAuthPostgres(db)

	// создаем массив структур из тестов
	tests := []struct {
		name    string
		mock    func()    // определение мока db
		input   todo.User // входные данные
		want    int       // какой должен быть ответ
		wantErr bool
	}{
		{
			name: "Ok",
			// определяем поведение мока
			mock: func() {
				// создаем мок ответ от db - "id": 1
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				// делаем запрос "INSERT INTO users" с параметрами "test", ожидая возврат данных - "id": 1
				mock.ExpectQuery("INSERT INTO users").WithArgs("test", "test", "test").WillReturnRows(rows)
			},
			input: todo.User{
				Name:     "test",
				Username: "test",
				Password: "test",
			},
			want: 1,
		},
		{
			name: "Empty fields",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO users").WithArgs("test", "test", "").WillReturnRows(rows)
			},
			input: todo.User{
				Name:     "test",
				Username: "test",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// проводим мок операции с db
			tt.mock()

			got, err := r.CreateUser(tt.input)
			// если есть ошибка, сопоставляем её с ожидающей ошибкой
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				// если нет ошибки, то сравниваем с нашим ответом
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func testAuthPostgres_getUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// создаем обязательное закрытие нашей моковой db
	defer db.Close()
	// создаем структуру authPostgres с доступом к db
	r := NewAuthPostgres(db)

	// делаем структуру для авторизации
	type args struct {
		username string
		password string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    todo.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password"}).AddRow(1, "test", "test", "test")
				mock.ExpectQuery("SELECT (.+) FROM users").WithArgs("test", "password").WillReturnRows(rows)
			},
			input: args{"test", "test"},
			want:  todo.User{1, "test", "test", "test"},
		},
		{
			name: "Now found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "password"})
				mock.ExpectQuery("SELECT (.+) FROM users").WithArgs("not", "found").WillReturnRows(rows)
			},
			input:   args{"not", "found"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := r.GetUser(tt.input.username, tt.input.password)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			// проверяем, что все ожидания ожиданы
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
