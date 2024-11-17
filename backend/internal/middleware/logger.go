package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// ANSI color codes
const (
    reset  = "\033[0m"
    red    = "\033[31m"
    green  = "\033[32m"
    yellow = "\033[33m"
    blue   = "\033[34m"
    purple = "\033[35m"
    cyan   = "\033[36m"
    gray   = "\033[37m"
    white  = "\033[97m"
)

// getStatusColor returns the appropriate color for the HTTP status code
func getStatusColor(code int) string {
    switch {
    case code >= 200 && code < 300:
        return green
    case code >= 300 && code < 400:
        return yellow
    case code >= 400 && code < 500:
        return red
    default:
        return purple
    }
}

func RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery
        
        // Get request headers
        headers := c.Request.Header

        c.Next()

        // end := time.Now()
        // latency := end.Sub(start)
        statusCode := c.Writer.Status()
        // clientIP := c.ClientIP()
        method := c.Request.Method
        // size := c.Writer.Size()
        
        if raw != "" {
            path = path + "?" + raw
        }

        statusColor := getStatusColor(statusCode)
        
        fmt.Printf("\n%s┌─────────────────────────────────────────────────%s\n", blue, reset)
        fmt.Printf("%s│%s 🌐 Request Details\n", blue, reset)
        fmt.Printf("%s├─────────────────────────────────────────────────%s\n", blue, reset)
        fmt.Printf("%s│%s %s%s %s%s%s\n", blue, reset, yellow, method, cyan, path, reset)
        fmt.Printf("%s│%s ⚡ Status: %s%d%s\n", blue, reset, statusColor, statusCode, reset)
        // fmt.Printf("%s│%s ⏱️  Latency: %s%v%s\n", blue, reset, purple, latency, reset)
        // fmt.Printf("%s│%s 🌍 IP: %s%s%s\n", blue, reset, green, clientIP, reset)
        // fmt.Printf("%s│%s 📦 Size: %s%d bytes%s\n", blue, reset, cyan, size, reset)
        
        // Print headers if they exist
        if len(headers) > 0 {
            fmt.Printf("%s│%s 📋 Headers:\n", blue, reset)
            for key, values := range headers {
                fmt.Printf("%s│%s   %s%s: %s%v%s\n", blue, reset, yellow, key, white, values, reset)
            }
        }
        
        // If there are any errors, print them
        if len(c.Errors) > 0 {
            fmt.Printf("%s│%s ❌ Errors: %s%v%s\n", blue, reset, red, c.Errors.String(), reset)
        }
        
        fmt.Printf("%s└─────────────────────────────────────────────────%s\n", blue, reset)
    }
}