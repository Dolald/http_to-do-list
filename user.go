package todo

type User struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"` // binding:"required" - значит, что в ответе обязательно должны быть данные поля, иначе будет ошибка
}
