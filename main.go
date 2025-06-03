package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	logDir        = "logs"
	logFileFormat = "log_%s.log"
	timeFormat    = "2006-01-02 15:04:05"
	dateFormat    = "2006-01-02"
)

// ensureLogDir ensures that the log directory exists
func ensureLogDir() error {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		return os.MkdirAll(logDir, 0755)
	}
	return nil
}

// getLogFilePath returns the log file path for a given date
func getLogFilePath(date string) string {
	return filepath.Join(logDir, fmt.Sprintf(logFileFormat, date))
}

// filterSensitiveData removes sensitive information from the body
func filterSensitiveData(data map[string]interface{}) map[string]interface{} {
	sensitiveKeys := []string{"password", "pwd", "secret", "token", "api_key", "apikey", "key"}
	filtered := make(map[string]interface{})

	for k, v := range data {
		isSensitive := false
		for _, sk := range sensitiveKeys {
			if strings.ToLower(k) == sk || strings.Contains(strings.ToLower(k), sk) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			filtered[k] = "[REDACTED]"
		} else if nestedMap, ok := v.(map[string]interface{}); ok {
			filtered[k] = filterSensitiveData(nestedMap)
		} else {
			filtered[k] = v
		}
	}

	return filtered
}

// logRequest logs the request to a file
func logRequest(c *gin.Context) {
	// Ensure log directory exists
	if err := ensureLogDir(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create log directory"})
		return
	}

	// Read request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}
	// Restore the request body for further processing
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	// Get current time
	now := time.Now()
	timeStr := now.Format(timeFormat)
	dateStr := now.Format(dateFormat)

	// Format body for logging
	var bodyFormatted string
	if len(bodyBytes) > 0 {
		var bodyData map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &bodyData); err == nil {
			// Filter sensitive data
			filteredData := filterSensitiveData(bodyData)

			// Convert filtered data to JSON string
			filteredJSON, err := json.MarshalIndent(filteredData, "", "  ")
			if err == nil {
				bodyFormatted = string(filteredJSON)
			} else {
				bodyFormatted = string(bodyBytes)
			}
		} else {
			// If not valid JSON, use as string
			bodyFormatted = string(bodyBytes)
		}
	}

	// Convert headers from map[string][]string to map[string]string
	flattenedHeaders := make(map[string]string)
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			flattenedHeaders[key] = values[0]
		}
	}

	// Format headers for logging
	headerJSON, err := json.MarshalIndent(flattenedHeaders, "", "  ")
	var headerFormatted string
	if err == nil {
		headerFormatted = string(headerJSON)
	} else {
		headerFormatted = fmt.Sprintf("%v", flattenedHeaders)
	}

	// Create log entry in the required format
	logEntry := fmt.Sprintf("#---\ntime: %s\n%s %s\n%s\n\n%s\n---\n",
		timeStr,
		c.Request.Method,
		c.Request.URL.String(),
		headerFormatted,
		bodyFormatted,
	)

	// Open log file in append mode
	logFilePath := getLogFilePath(dateStr)
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open log file"})
		return
	}
	defer logFile.Close()

	// Write log entry to file
	if _, err := logFile.WriteString(logEntry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write log entry"})
		return
	}

	// Continue with the request
	c.Next()
}

// getLogHandler handles requests to get logs for a specific date
func getLogHandler(c *gin.Context) {
	// Get date parameter
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date parameter is required"})
		return
	}

	// Validate date format
	if _, err := time.Parse(dateFormat, date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	// Get log file path
	logFilePath := getLogFilePath(date)

	// Check if log file exists
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log file not found for the specified date"})
		return
	}

	// Read log file
	logContent, err := os.ReadFile(logFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read log file"})
		return
	}

	// Return log content
	c.String(http.StatusOK, string(logContent))
}

func main() {
	// Create a new Gin router
	r := gin.Default()

	// Apply log middleware to all routes except the log endpoint
	r.Use(func(c *gin.Context) {
		// Skip logging for the log endpoint
		if c.Request.URL.Path == "/api/v1/log" {
			c.Next()
			return
		}

		// Log the request
		logRequest(c)

		// Return 200 OK
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// Define the log endpoint
	r.GET("/api/v1/log", getLogHandler)

	// Start the server
	fmt.Println("Server is running on :8081")
	r.Run(":8081")
}
