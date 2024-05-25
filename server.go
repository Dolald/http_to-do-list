package todo

import (
	"context"
	"net/http"
	"time"

	_ "github.com/lib/pq" // add this
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error { // метод по запуску сервера
	s.httpServer = &http.Server{
		Addr:           ":" + port,       // адрес порта
		Handler:        handler,          // используем наш созданный обработчик
		MaxHeaderBytes: 1 >> 20,          // определяет максимальное количество байт, которое сервер будет считывать при анализе ключей и значений заголовка запроса
		ReadTimeout:    10 * time.Second, // макс продолжительность обработки запроса
		WriteTimeout:   10 * time.Second, // макс время до истечения обработки запроса
	}
	return s.httpServer.ListenAndServe() // слушаем созданное ТСР соединение
}

func (s *Server) Shutdown(ctx context.Context) error { // не знаю что это :(
	return s.httpServer.Shutdown(ctx)
}
