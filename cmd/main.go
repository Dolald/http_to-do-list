package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	todo "todolist"
	"todolist/pkg/handler"
	"todolist/pkg/repository"
	"todolist/pkg/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// @title Todo App API
// @version 1.0
// @description API Server for TodoList Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter)) // ошибка выводится в формате JSON

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil { // читаем файл окружения
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{ // подключаемся к нашей БД
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	server := new(todo.Server) // создаём новый сервер

	go func() { // запускаем сервер в анонимной горутине чтоб не блокировал основной поток ?
		if err := server.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil { // запускаем сервер на 8000 порту через заданный обработчик handler.InitRoutes()
			logrus.Fatalf("error occured while runnung http server: %s", err.Error())
		}
	}()

	logrus.Print("TodoApp started")

	quit := make(chan os.Signal, 1)                      // делаем канал типа сигнал размером 1
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT) // подписываемся на 2 системных сигнала
	<-quit                                               // ожидаем сигнал из канала, до этого программа не останавливается

	logrus.Print("ToDoApp shuttong down")

	if err := server.Shutdown(context.Background()); err != nil { // аккуратно завершаем работу сервера
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs") // говорим искать папку configs
	viper.SetConfigName("config")  // говорим искать файл config
	return viper.ReadInConfig()    // читаем
}
