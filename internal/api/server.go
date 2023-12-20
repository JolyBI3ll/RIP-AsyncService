package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const password = "Golang"

func StartServer() {
	log.Println("Server start up")

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.POST("/result", func(c *gin.Context) {
		var data EventData

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		requestID := data.Request_id

		// Запуск горутины для отправки статуса
		go sendStatus(requestID, password, fmt.Sprintf("http://localhost:8000/request/%d/status/", requestID))

		c.JSON(http.StatusOK, gin.H{"message": "Status update initiated"})
	})
	router.Run(":8080")

	log.Println("Server down")
}

func genRandomStatus(password string) Result {
	time.Sleep(8 * time.Second)
	status := "W"
	if rand.Intn(100) < 50 {
		status = "F"
	}
	return Result{status, password}
}

// Функция для отправки статуса в отдельной горутине
func sendStatus(requestID int, password string, url string) {
	// Выполнение расчётов с randomStatus
	result := genRandomStatus(password)

	// Отправка PUT-запроса к основному серверу
	_, err := performPUTRequest(url, result)
	if err != nil {
		fmt.Println("Error sending status:", err)
		return
	}

	fmt.Println("Status sent successfully for RequestID:", requestID)
}

type Result struct {
	Status   string `json:"status"`
	Password string `json:"password"`
}

type EventData struct {
	Request_id int `json:"Request_id"`
}

func performPUTRequest(url string, data Result) (*http.Response, error) {
	// Сериализация структуры в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Создание PUT-запроса
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}
