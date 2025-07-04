package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/IlyaStarshinov/orderWB/internal/cache"
	dbpkg "github.com/IlyaStarshinov/orderWB/internal/db"
	"github.com/IlyaStarshinov/orderWB/internal/model"
	"gorm.io/gorm"
)

func StartConsumer(brokers []string, topic string, db *gorm.DB, cache *cache.Cache) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	config.Version = sarama.V2_0_0_0

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatal("Ошибка создания потребителя Kafka:", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Ошибка создания потребителя партиции:", err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			processMessage(msg.Value, db, cache)
		case err := <-partitionConsumer.Errors():
			log.Println("Ошибка Kafka:", err)
		}
	}
}

func processMessage(data []byte, db *gorm.DB, cache *cache.Cache) {
	log.Printf("Получено сообщение: %s", string(data))

	var order model.Order
	if err := json.Unmarshal(data, &order); err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		return
	}

	log.Printf("Обработка заказа: %s", order.OrderUID)

	if err := dbpkg.SaveOrder(db, &order); err != nil {
		log.Printf("Ошибка сохранения в БД: %v", err)
		return
	}

	cache.Set(&order)
	log.Printf("Заказ успешно обработан: %s", order.OrderUID)
}
