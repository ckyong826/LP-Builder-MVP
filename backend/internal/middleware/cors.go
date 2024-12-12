package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
    // Get allowed origins from environment variable
    allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
    
    // Default origins if none specified
    if len(allowedOrigins) == 0 || (len(allowedOrigins) == 1 && allowedOrigins[0] == "") {
        allowedOrigins = []string{
            "http://localhost:3000",
            "http://localhost:5173",
            "https://*.pinggy.io",
            "https://*.pinggy-ao.net",
						"https://*.a.free.pinggy.link",
            "http://*.pinggy.io",
            "http://*.pinggy-ao.net",
						"http://*.a.free.pinggy.link",
						"https://lp-builder-mvp.vercel.app",
						"https://lp-builder-mvp.vercel.app/",
        }
    }

    return cors.New(cors.Config{
        AllowOrigins:     allowedOrigins,
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept", "Authorization", "X-Requested-With"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        AllowWildcard:    true,  // Important for wildcard domains
        MaxAge:           12 * time.Hour,
    })
}