package main

import (
	"log"
	"net/http"

	"github.com/khavazh7/matchmaker/handler"

	"github.com/khavazh7/matchmaker/internal/matchmaker"

	"github.com/khavazh7/matchmaker/config"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Выбираем хранилище
	var storage matchmaker.Storage
	switch cfg.StorageType {
	case "memory":
		storage = matchmaker.NewInMemoryStorage()
	case "postgres":
		var err error
		storage, err = matchmaker.NewPostgresStorage(cfg.DBConfig)
		if err != nil {
			log.Fatalf("Failed to connect to the database: %v", err)
		}
	default:
		log.Fatalf("Unknown storage type: %s", cfg.StorageType)
	}

	// Создаем матчмейкер
	matcher := matchmaker.NewMatcher(storage, cfg.GroupSize)

	// Запускаем HTTP сервер
	http.HandleFunc("/users", handler.CreateUserHandler(matcher))
	log.Printf("Starting server on %s", cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, nil))
}
