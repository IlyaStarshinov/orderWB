# Order Processing Microservice

Микросервис на языке Go для приёма, хранения и отображения заказов с использованием Apache Kafka, PostgreSQL и кэша в оперативной памяти.

## Функциональные возможности
- Получение сообщений в формате JSON из Kafka
- Сохранение заказов и связанных сущностей в базе данных PostgreSQL
- Кэширование заказов в памяти для ускоренного доступа
- Восстановление кэша из базы данных при старте приложения
- HTTP API и веб-интерфейс для просмотра информации о заказах

## Используемые технологии
- Язык программирования Go 1.21+
- PostgreSQL 14
- Apache Kafka 3.9
- Zookeeper
- GORM (ORM для Go)
- Sarama (клиент Kafka для Go)
- Gorilla/mux (роутер для HTTP)
- HTML и JavaScript (веб-интерфейс)

## Инструкция по установке и запуску

### Установка PostgreSQL

Скачать установщик с официального сайта: https://www.postgresql.org/download/windows/

### Создание базы данных и пользователя

psql -U postgres

В интерактивной консоли PostgreSQL выполнить:

CREATE DATABASE orders_db;
CREATE USER order_user WITH PASSWORD 'order_password';
GRANT ALL PRIVILEGES ON DATABASE orders_db TO order_user;

### Применение SQL-модуля с миграциями

psql -U order_user -d orders_db -f init.sql

## Запуск Kafka и Zookeeper (без Docker)

### Запуск (в двух отдельных окнах терминала)

### Zookeeper:

D:\kafka_2.13-3.9.1\bin\windows\zookeeper-server-start.bat D:\kafka_2.13-3.9.1\config\zookeeper.properties

### Kafka:

D:\kafka_2.13-3.9.1\bin\windows\kafka-server-start.bat D:\kafka_2.13-3.9.1\config\server.properties

### Создание топика Kafka

D:\kafka_2.13-3.9.1\bin\windows\kafka-topics.bat --create --topic orders --bootstrap-server localhost:9092

## Сборка и запуск Go-приложения

### Установка зависимостей

go mod tidy

### Запуск сервера

go run cmd/server/main.go

## Отправка тестового сообщения в Kafka

1. Создать файл model.json с валидным JSON-сообщением заказа.
2. Отправить его через консоль:

$json = Get-Content -Raw .\model.json
$json | D:\kafka_2.13-3.9.1\bin\windows\kafka-console-producer.bat --broker-list localhost:9092 --topic orders
