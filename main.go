package handler

import (
	"net/http"
	"sync"

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
	sync.Mutex
	EmailsSent  int
	EmailLimit  int
	EmailWarmUp bool
}

var stats = Stats{EmailLimit: 10, EmailWarmUp: true}

// Handler function for the Go app, used in Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	router := gin.Default()

	// **Enable CORS with specific frontend origin**
	config := cors.Config{
		AllowOrigins:     []string{"https://aws-ses-api-fgq1.vercel.app"}, // Replace with your frontend URL
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}
	router.Use(cors.New(config))

	// Root route for testing
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the mock SES API!"})
	})

	// Routes
	router.POST("/sendEmail", sendEmail)
	router.GET("/stats", getStats)
	router.GET("/healthcheck", healthCheck)

	// Forward the request to Gin
	router.ServeHTTP(w, r)
}

// sendEmail handles the sending of email (mock behavior)
func sendEmail(c *gin.Context) {
	var email Email
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	stats.Lock()
	defer stats.Unlock()

	if stats.EmailWarmUp && stats.EmailsSent >= stats.EmailLimit {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Email warm-up period active. Please try again later."})
		return
	}

	stats.EmailsSent++
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully", "to": email.To})
}

// getStats returns the statistics of the mock SES API
func getStats(c *gin.Context) {
	stats.Lock()
	defer stats.Unlock()
	c.JSON(http.StatusOK, gin.H{
		"emails_sent": stats.EmailsSent,
		"email_limit": stats.EmailLimit,
	})
}

// healthCheck checks if the API is alive
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}
