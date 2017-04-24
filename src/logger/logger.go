package logger

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry/gosteno"
	"github.com/gin-gonic/gin"
	"time"
)

type Logger interface {
	Info(subject string, messages ...map[string]string)
	Debug(subject string, messages ...map[string]string)
	Error(subject string, err error, messages ...map[string]string)
}

type RealLogger struct {
	steno *gosteno.Logger
}

func NewRealLogger(steno *gosteno.Logger) *RealLogger {
	return &RealLogger{
		steno: steno,
	}
}

func (logger *RealLogger) Debug(subject string, messages ...map[string]string) {
	logger.steno.Debug(subject + logger.parseMessages(messages))
}

func (logger *RealLogger) Info(subject string, messages ...map[string]string) {
	logger.steno.Info(subject + logger.parseMessages(messages))
}

func (logger *RealLogger) Error(subject string, err error, messages ...map[string]string) {
	logger.steno.Error(subject + " - Error:" + err.Error() + logger.parseMessages(messages))
}

func (logger *RealLogger) parseMessages(messages []map[string]string) string {
	messageString := ""
	for _, message := range messages {
		messageBytes, _ := json.Marshal(message)
		messageString += " - " + string(messageBytes)
	}

	return messageString
}

func FormatLogger(start time.Time, statusCode int, c *gin.Context, errorString string) string {
	end := time.Now()
	latency := end.Sub(start)
	clientIP := c.Request.RemoteAddr
	method := c.Request.Method
	logString := fmt.Sprintf(" %3d | %v | %12v | %s | %s | %s | ERROR: %s",
		statusCode,
		end.Format("2006/01/02 - 15:04:05"),
		latency,
		clientIP,
		method,
		c.Request.URL.Path,
		errorString,
	)
	return logString
}
