package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type SortStep struct {
	Array []int `json:"array"`
	Pivot int   `json:"pivot,omitempty"`
	Left  int   `json:"left,omitempty"`
	Right int   `json:"right,omitempty"`
}

var steps []SortStep

func quickSort(arr []int, low, high int) {
	if low < high {
		pi := partition(arr, low, high)
		quickSort(arr, low, pi-1)
		quickSort(arr, pi+1, high)
	}
}

func partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
			logStep(arr, pivot, i, j)
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	logStep(arr, pivot, i+1, high)
	return i + 1
}

func logStep(arr []int, pivot, left, right int) {
	// Create a deep copy of the array to avoid mutation issues
	arrCopy := make([]int, len(arr))
	copy(arrCopy, arr)
	steps = append(steps, SortStep{Array: arrCopy, Pivot: pivot, Left: left, Right: right})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	// Apply the CORS middleware to every request
	r.Use(CORSMiddleware())

	r.POST("/quicksort", func(c *gin.Context) {
		var input struct {
			Array []int `json:"array"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		steps = nil // reset steps
		quickSort(input.Array, 0, len(input.Array)-1)
		c.JSON(http.StatusOK, gin.H{"steps": steps})
	})
	return r
}

func main() {
	r := setupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	r.Run(":" + port)
}

