package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/IlyaStarshinov/orderWB/internal/cache"

	"github.com/gorilla/mux"
)

func GetOrderHandler(cache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderUID := vars["order_uid"]
		log.Printf("Получен запрос для заказа UID: %s", orderUID)

		// Пробуем получить из кэша
		order, exists := cache.Get(orderUID)
		if !exists {
			log.Printf("Заказ %s не найден в КЭШЕ", orderUID)
			http.Error(w, "Заказ не найден", http.StatusNotFound)
			return
		}

		log.Printf("Заказ %s предоставлен из КЭША", orderUID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
	}
}
