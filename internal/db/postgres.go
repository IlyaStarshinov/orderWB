package db

import (
	"log"

	"github.com/IlyaStarshinov/orderWB/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := "host=localhost user=order_user password=1234 dbname=orders_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Ошибка подключения к базе данных:", err)
		return nil, err
	}

	db.Exec("SET search_path TO public")

	// Проверка подключения
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Ошибка получения экземпляра БД:", err)
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		log.Println("Ошибка ping базы данных:", err)
		return nil, err
	}
	log.Println("Успешное подключение к PostgreSQL!")

	// Автомиграция
	log.Println("Запуск AutoMigrate...")
	err = db.AutoMigrate(
		&model.Order{},
		&model.Delivery{},
		&model.Payment{},
		&model.Item{},
	)
	if err != nil {
		log.Println("Ошибка AutoMigrate:", err)
		return nil, err
	}
	log.Println("AutoMigrate успешно завершен")

	return db, nil
}

func SaveOrder(db *gorm.DB, order *model.Order) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}
