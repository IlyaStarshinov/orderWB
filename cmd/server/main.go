package main

import (
	"log"
	"net/http"
	"os"

	"github.com/IlyaStarshinov/orderWB/internal/cache"
	"github.com/IlyaStarshinov/orderWB/internal/db"
	"github.com/IlyaStarshinov/orderWB/internal/handler"
	"github.com/IlyaStarshinov/orderWB/internal/kafka"

	"github.com/gorilla/mux"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Запуск приложения...")

	// Инициализация БД
	log.Println("Инициализация подключения к базе данных...")
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}
	log.Println("База данных успешно инициализирована")

	// Инициализация кэша
	log.Println("Инициализация кэша...")
	cache := cache.NewCache()
	if err := cache.RestoreFromDB(dbConn); err != nil {
		log.Fatal("Ошибка восстановления кэша из БД:", err)
	}
	log.Println("Кэш успешно инициализирован")

	// Запуск Kafka Consumer
	log.Println("Запуск Kafka Consumer...")
	go kafka.StartConsumer(
		[]string{"localhost:9092"},
		"orders",
		dbConn,
		cache,
	)
	log.Println("Kafka Consumer запущен в фоновом режиме")

	// HTTP сервер
	r := mux.NewRouter()
	r.HandleFunc("/order/{order_uid}", handler.GetOrderHandler(cache)).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))

	log.Println("Сервер запущен на порту :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
