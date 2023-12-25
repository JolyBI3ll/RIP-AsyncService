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

const token = "Golang"

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

		if data.Token != token {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "нет доступа"})
			return
		}

		requestID := data.Request_id

		// Запуск горутины для отправки статуса
		go sendStatus(requestID, token, fmt.Sprintf("http://localhost:8000/request/%d/result/", requestID))

		c.JSON(http.StatusOK, gin.H{"message": "Result update initiated"})
	})
	router.Run(":8080")

	log.Println("Server down")
}

func genRandomStatus(token string) Result {
	time.Sleep(8 * time.Second)
	result := "W"
	if rand.Intn(100) < 50 {
		result = "F"
	}
	return Result{result, token}
}

// Функция для отправки статуса в отдельной горутине
func sendStatus(requestID int, token string, url string) {
	// Выполнение расчётов с randomStatus
	result := genRandomStatus(token)

	// Отправка PUT-запроса к основному серверу
	_, err := performPUTRequest(url, result)
	if err != nil {
		fmt.Println("Error sending result:", err)
		return
	}

	fmt.Println("Result sent successfully for RequestID:", requestID)
}

type Result struct {
	Result string `json:"result"`
	Token  string `json:"token"`
}

type EventData struct {
	Request_id int    `json:"Request_id"`
	Token      string `json:"token"`
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
