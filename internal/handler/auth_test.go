package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"
	todo "todolist"
	"todolist/internal/service"
	service_mocks "todolist/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestSignUp(t *testing.T) {

	type mockBehavior func(r *service_mocks.MockAuthorization, user todo.User)

	tests := []struct {
		name                 string       // итог теста
		inputBody            string       // тело запроса
		inputUser            todo.User    // наш юзер
		mockBehavior         mockBehavior // определяем поведение теста
		expectedStatusCode   int          // ожидаемый статус код
		expectedResponseBody string       // ожидаемый ответ
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "username", "name": "Test Name", "password": "qwerty"}`,
			inputUser: todo.User{
				Username: "username",
				Name:     "Test Name",
				Password: "qwerty",
			},
			mockBehavior: func(r *service_mocks.MockAuthorization, user todo.User) {
				r.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Wrong input",
			inputBody:            `{"username": "username"}`,
			inputUser:            todo.User{},
			mockBehavior:         func(r *service_mocks.MockAuthorization, user todo.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"username": "username", "name": "Test Name", "password": "qwerty"}`,
			inputUser: todo.User{
				Username: "username",
				Name:     "Test Name",
				Password: "qwerty",
			},
			mockBehavior: func(r *service_mocks.MockAuthorization, user todo.User) {
				r.EXPECT().CreateUser(user).Return(0, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		// запускаем каждый отдельный тест
		t.Run(test.name, func(t *testing.T) {
			// создаём новый контроллер моков
			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockAuthorization(c) // создаём новый мок объект
			test.mockBehavior(repo, test.inputUser)       // настраиваем тест

			services := &service.Service{Authorization: repo} // создаём серви с нашим моком
			handler := Handler{services}                      // делаем хендлер с нашим сервисом

			r := gin.New()                     // создаем новый роутер
			r.POST("/sign-up", handler.signUp) // настраиваем маршрут POST

			w := httptest.NewRecorder()                                                           // Создаем запрос
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(test.inputBody)) // делаем запрос

			r.ServeHTTP(w, req) // Make Request

			assert.Equal(t, w.Code, test.expectedStatusCode)            // сверяем кода статуса
			assert.Equal(t, w.Body.String(), test.expectedResponseBody) // сравниваем результат и ожидаемый результат
		})
	}
}
