package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	r := gin.Default()

	routes := map[string]string{
		"/products": "http://product-service:8081",
		"/cart":     "http://cart-service:8082",
		"/orders":   "http://order-service:8083",
		"/auth":     "http://auth-service:8084",
		"/payments": "http://payment-service:8086",
	}

	for prefix, target := range routes {
		// Parse the target URL
		targetURL, err := url.Parse(target)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to parse target URL")
		}
		// Create the proxy handler
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		// Custom Director to ensure the path is forwarded correctly
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			// Modify the request to point to the target host
			req.Host = targetURL.Host
			// e.g., /products/1 -> product-service:8081/products/1
		}
		// Register the route ("/*any" captures all sub-paths)
		r.Any(prefix+"/*any", func(c *gin.Context) {
			proxy.ServeHTTP(c.Writer, c.Request)
		})
		// Also register the exact prefix match (without trailing slash)
		r.Any(prefix, func(c *gin.Context) {
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	}
	// 3. Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "API Gateway Running"})
	})
	log.Info().Msg("API Gateway starting on :8080")
	r.Run(":8080")
}
