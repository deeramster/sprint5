package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

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
	schemaAPIKey := os.Getenv("SCHEMA_REGISTRY_API_KEY")
	caLocation := os.Getenv("SSL_CA_LOCATION")

	// Подключение к Schema Registry
	schemaRegistryClient := srclient.CreateSchemaRegistryClient(schemaURL)
	schemaRegistryClient.SetCredentials(schemaAPIKey, "")

	// Проверка наличия схемы (просто для использования переменной)
	_, err := schemaRegistryClient.GetLatestSchema(schemaSubject)
	if err != nil {
		log.Printf("Предупреждение при получении схемы: %v", err)
		// Не fatal, так как можно продолжить работу
	}

	// Создание продюсера Kafka
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"security.protocol": "SASL_SSL",
		"sasl.mechanism":    "SCRAM-SHA-512",
		"sasl.username":     username,
		"sasl.password":     password,
		"ssl.ca.location":   caLocation,
	})
	if err != nil {
		log.Fatalf("Ошибка создания продюсера: %v", err)
	}
	defer p.Close()

	// Создание JSON-сообщения
	message := Message{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Text:      "Привет, Kafka с JSON!",
		Timestamp: time.Now().Unix(),
	}

	// Кодирование в JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Ошибка кодирования JSON: %v", err)
	}

	// Канал для обработки результатов
	deliveryChan := make(chan kafka.Event, 1)

	// Отправка сообщения в Kafka
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          jsonData,
	}, deliveryChan)
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %v", err)
	}

	// Ожидание подтверждения отправки
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Ошибка доставки сообщения: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Сообщение отправлено в топик %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	p.Flush(15 * 1000)
}
