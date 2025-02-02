package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Email struct to represent the email sending request
type Email struct {
	To      string `json:"to" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

// Stats struct to hold email statistics
type Stats struct {
	EmailsSent  int
	EmailLimit  int
	EmailWarmUp bool
}

var stats = Stats{EmailLimit: 10, EmailWarmUp: true}

// Create Gin router
func createRouter() *gin.Engine {
	router := gin.Default()

	// Enable CORS
	router.Use(cors.Default())

	// Root route for testing
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the mock SES API!"})
	})

	// Routes
	router.POST("/sendEmail", sendEmail)
	router.GET("/stats", getStats)
	router.GET("/healthcheck", healthCheck)

	return router
}

// sendEmail handles the sending of email (mock behavior)
func sendEmail(c *gin.Context) {
	var email Email
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Simulate email warming (only a few emails allowed during warm-up period)
	if stats.EmailWarmUp && stats.EmailsSent >= stats.EmailLimit {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Email warm-up period active. Please try again later."})
		return
	}

	// Simulate email sending
	stats.EmailsSent++
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully", "to": email.To})
}

// getStats returns the statistics of the mock SES API
func getStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"emails_sent": stats.EmailsSent,
		"email_limit": stats.EmailLimit,
	})
}

// healthCheck checks if the API is alive
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

// Vercel expects a handler to be exported for serverless
func Handler(w http.ResponseWriter, r *http.Request) {
	router := createRouter()
	router.ServeHTTP(w, r)
}

func main() {
	// Local testing: Run the server
	// This should be disabled in production serverless environments like Vercel
	if err := http.ListenAndServe(":8080", http.HandlerFunc(Handler)); err != nil {
		log.Fatal("Failed to start the server: ", err)
	}
}
