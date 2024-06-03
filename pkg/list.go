package handler

import (
	"net/http"
	"strconv"
	todo "todolist"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c) // получаем id пользователя
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input todo.TodoList

	if err := c.BindJSON(&input); err != nil { // парсим input, где находятся созданный список заданий
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.ToDoList.Create(userId, input) // создаём новую задачу
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]any{ // присылаем ответ в виде id
		"id": id,
	})
}

type getAllListResponse struct {
	Data []todo.TodoList `json:"data"`
}

func (h *Handler) getAllList(c *gin.Context) {
	userId, err := getUserId(c) // получаем id пользователя
	if err != nil {
		return
	}

	allList, err := h.service.ToDoList.GetAllLists(userId)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, getAllListResponse{
		Data: allList,
	})
}

func (h *Handler) getListById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	list, err := h.service.ToDoList.GetById(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handler) updateList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	idList, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
	}

	type list todo.UpdateListInput

	if err = c.BindJSON(&list); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedList, err := h.service.ToDoList.UpdateList(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	c.JSON(http.StatusOK)
}

func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.service.ToDoList.DeleteList(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
