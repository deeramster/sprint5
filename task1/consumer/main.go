package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/joho/godotenv"
	"github.com/riferrei/srclient"
)

type Message struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
}

func main() {
	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	// Чтение параметров из переменных окружения
	broker := os.Getenv("KAFKA_BROKERS")
	schemaURL := os.Getenv("SCHEMA_REGISTRY_URL")
	topic := os.Getenv("KAFKA_TOPIC")
	schemaSubject := os.Getenv("SCHEMA_SUBJECT")

	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	caLocation := os.Getenv("SSL_CA_LOCATION")
	schemaAPIKey := os.Getenv("SCHEMA_REGISTRY_API_KEY")
	consumerGroup := os.Getenv("KAFKA_CONSUMER_GROUP")

	// Подключение к Schema Registry
	schemaRegistryClient := srclient.CreateSchemaRegistryClient(schemaURL)
	schemaRegistryClient.SetCredentials(schemaAPIKey, "")

	_, err := schemaRegistryClient.GetLatestSchema(schemaSubject)
	if err != nil {
		log.Fatalf("Ошибка получения схемы: %v", err)
	}

	// Создание консьюмера Kafka
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          consumerGroup,
		"auto.offset.reset": "earliest",
		"security.protocol": "SASL_SSL",
		"sasl.mechanism":    "SCRAM-SHA-512",
		"sasl.username":     username,
		"sasl.password":     password,
		"ssl.ca.location":   caLocation,
	})
	if err != nil {
		log.Fatalf("Ошибка создания консьюмера: %v", err)
	}
	defer c.Close()

	// Подписка на топик
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("Ошибка подписки на топик: %v", err)
	}

	// Чтение сообщений
	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			log.Printf("Ошибка чтения сообщения: %v", err)
			continue
		}

		var message Message
		err = json.Unmarshal(msg.Value, &message)
		if err != nil {
			log.Printf("Ошибка декодирования JSON: %v", err)
			continue
		}

		fmt.Printf("Получено сообщение: %+v\n", message)
		fmt.Printf("Метаданные сообщения: Топик=%s, Партиция=%d, Оффсет=%d\n",
			*msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
	}
}
