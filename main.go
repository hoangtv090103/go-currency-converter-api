package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2" // Thực hiện các yêu cầu HTTP tới API bên thứ ba.
	"github.com/joho/godotenv"
)

type ConvesionRequest struct {
	From   string  `json:"from" binding:"required"`
	To     string  `json: "to" binding: "required"`
	Amount float64 `json: "amount" binding: "required"`
}

type ConversionResponse struct {
	ConvertedAmount float64 `json:"converted_amount"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := gin.Default()

	r.GET("/convert", func(c *gin.Context) {
		var request ConvesionRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		apiKey := os.Getenv("EXCHANGE_RATE_API_KEY")

		url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/pair/%s/%s/%f", apiKey, request.From, request.To, request.Amount)
		var result map[string]interface{}

		client := resty.New()
		resp, err := client.
			R().
			SetHeader("Accept", "application/json").
			SetResult(&result).
			Get(url)

		if err != nil || resp.StatusCode() != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get exchange rate"})
			return
		}
		fmt.Println(result["conversion_result"])
		res, ok := result["conversion_result"].(float64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid conversion rate"})
			return
		}

		c.JSON(http.StatusOK, ConversionResponse{ConvertedAmount: res})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
