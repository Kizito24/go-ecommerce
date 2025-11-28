package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 1. CORS Setup (Keep this from before)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 2. Define the Auth Service URL
	// In Docker, this will be "http://auth-service:5001"
	// Locally, we use "http://localhost:5001"
	authUrl := os.Getenv("AUTH_SVC_URL")
	if authUrl == "" {
		authUrl = "http://localhost:5001"
	}

	// 3. Register the Proxy Route
	// "Any" captures GET, POST, PUT, DELETE, etc.
	// "/*proxyPath" is a wildcard that captures everything after /auth/
	r.Any("/auth/*proxyPath", proxyToService(authUrl))

	log.Println("Gateway running on port 8080")
	r.Run(":8080")
}

// 4. The Reverse Proxy Logic
// This function returns a Gin Handler that forwards the request
func proxyToService(targetHost string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the target URL
		remote, err := url.Parse(targetHost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid proxy target"})
			return
		}

		// Create a Reverse Proxy
		proxy := httputil.NewSingleHostReverseProxy(remote)

		// Define how to rewrite the request (optional, but good practice)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			// Note: We keep the path exactly the same.
			// Gateway: /auth/login -> Auth Service: /auth/login
			req.URL.Path = c.Request.URL.Path
		}

		// Serve the request
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
