package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// User представляет структуру данных
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

// generateRandomData генерирует случайные данные для пользователей
func generateRandomData() User {
	names := []string{"Екатерина", "Никита", "Майя", "Алексей", "Илья", "Виктория"}
	emails := []string{"kate@yandex.ru", "nikita@yandex.ru", "maya@yandex.ru", "alex@yandex.ru", "ilya@yandex.ru", "vika@yandex.ru"}

	return User{
		ID:    rand.Intn(1000), // случайный ID
		Name:  names[rand.Intn(len(names))],
		Age:   rand.Intn(30) + 18, // возраст от 18 до 47
		Email: emails[rand.Intn(len(emails))],
	}
}

// sendToNiFi отправляет данные в NiFi
func sendToNiFi(user User) error {
	url := "http://nifi:8085" // Адрес ListenHTTP в NiFi

	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("ошибка при маршалинге JSON: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("ошибка при отправке в NiFi: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("NiFi вернул ошибку: %s", resp.Status)
	}

	log.Println("Данные успешно отправлены в NiFi:", user)
	return nil
}

func main() {
	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Потоковая генерация данных каждую секунду
	for {
		// Генерируем случайного пользователя
		user := generateRandomData()

		// Отправляем данные в NiFi
		if err := sendToNiFi(user); err != nil {
			log.Println("Ошибка:", err)
		}

		// Ожидание 1 секунду перед следующей отправкой
		time.Sleep(1 * time.Second)
	}
}
