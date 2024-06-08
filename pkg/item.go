package handler

import (
	"net/http"
	"strconv"
	todo "todolist"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c)

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	var input todo.TodoItem

	if err != c.BindJSON(&input) {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.ToDoItem.Create(userId, listId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"id": id,
	})
}

func (h *Handler) getAllItem(c *gin.Context) {

}

func (h *Handler) getItemById(c *gin.Context) {

}

func (h *Handler) updateItem(c *gin.Context) {

}

func (h *Handler) deleteItem(c *gin.Context) {

}
