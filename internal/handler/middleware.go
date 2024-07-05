package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentify(c *gin.Context) {
	header := c.GetHeader(authorizationHeader) // получаем методанные от хедера HTTP запроса
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ") // в хедере Authorization может находиться такое "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ" или такое "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXV"
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
	}

	userId, err := h.service.Authorization.ParseToken(headerParts[1]) // "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ" - headerParts[1] = "QWxhZGRpbjpvcGVuIHNlc2FtZQ"
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is not found")
		return 0, errors.New("user id is not found")
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return 0, errors.New("user id is not found")
	}

	return idInt, nil
}
