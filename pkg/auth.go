package handler

import (
	"net/http"
	todo "todolist"

	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) { // c *gin.Context - вариант w http.ResponseWriter, r *http.Request в стандартном пакете http
	var input todo.User // создаём экземпляр переменной с типом нашей структуры из другого пакета

	if err := c.BindJSON(&input); err != nil { // засовываем ответ в input
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

}
func (h *Handler) signIn(c *gin.Context) {

}
