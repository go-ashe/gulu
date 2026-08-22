package middleware

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-ashe/gulu/logx"
)

const (
	maxRequestBodySize = 10 * 1024 * 1024
	maxLoggedBodySize  = 2000
	maxLoggedDataSize  = 200
	requestIDKey       = "X-Request-Id"
)

// Logger 记录请求和响应信息。合法 JSON 请求体会作为结构化字段输出，
// 避免将 JSON 当作字符串再次编码后产生大量转义字符。
func Logger(c *gin.Context) {
	body, err := readRequestBody(c.Request, maxRequestBodySize)
	if err != nil {
		logx.FromGin(c).Error("读取 Body 失败", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体过大"})
		c.Abort()
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	requestID := getRequestID(c)
	reqLogger := logx.FromGin(c).WithAttrs("request_id", requestID)
	c.Set("logger", reqLogger)

	reqLogger.Info("请求开始",
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"remote_addr", c.ClientIP(),
		"body", requestBodyLogValue(body, c.GetHeader("Content-Type")),
	)

	responseWriter := &responseLogWriter{ResponseWriter: c.Writer}
	c.Writer = responseWriter

	start := time.Now()
	c.Next()

	status := c.Writer.Status()
	if status == 0 {
		status = http.StatusOK
	}

	reqLogger.Info("请求完成",
		slog.String("elapsed", time.Since(start).String()),
		slog.Int("status", status),
		slog.Int("body_size", c.Writer.Size()),
		slog.Any("body", responseWriter.body),
	)
}

// readRequestBody 安全读取 Body（带大小限制）
func readRequestBody(r *http.Request, maxSize int64) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, maxSize+1))
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > maxSize {
		return nil, io.ErrShortBuffer
	}
	return body, nil
}

func requestBodyLogValue(body []byte, contentType string) any {
	if len(body) == 0 {
		return nil
	}

	switch {
	case strings.Contains(contentType, "application/json"):
		if len(body) <= maxLoggedBodySize && json.Valid(body) {
			return json.RawMessage(body)
		}
		return truncateString(string(body), maxLoggedBodySize)
	case strings.Contains(contentType, "application/x-www-form-urlencoded"):
		return truncateString(string(body), maxLoggedBodySize)
	case strings.Contains(contentType, "multipart/form-data"):
		return "[multipart/不记录]"
	default:
		return "[二进制/不记录]"
	}
}

func getRequestID(c *gin.Context) string {
	if value, exists := c.Get(requestIDKey); exists {
		if requestID, ok := value.(string); ok && requestID != "" {
			return requestID
		}
	}
	if requestID := c.GetHeader(requestIDKey); requestID != "" {
		return requestID
	}

	var value [16]byte
	if _, err := rand.Read(value[:]); err == nil {
		return hex.EncodeToString(value[:])
	}
	return "unknown"
}

func truncateString(value string, maxLen int) string {
	if len(value) <= maxLen {
		return value
	}
	return value[:maxLen] + " [TRUNCATED]"
}

type responseLogWriter struct {
	gin.ResponseWriter
	body     any
	captured bool
}

func (w *responseLogWriter) Write(data []byte) (int, error) {
	w.capture(data)
	return w.ResponseWriter.Write(data)
}

func (w *responseLogWriter) WriteString(data string) (int, error) {
	w.capture([]byte(data))
	return w.ResponseWriter.WriteString(data)
}

func (w *responseLogWriter) capture(data []byte) {
	if w.captured {
		return
	}
	w.captured = true
	w.body = responseBodyLogValue(data)
}

func responseBodyLogValue(body []byte) any {
	if len(body) == 0 {
		return nil
	}

	var response map[string]json.RawMessage
	if err := json.Unmarshal(body, &response); err != nil {
		return truncateString(string(body), maxLoggedBodySize)
	}

	data, exists := response["data"]
	if !exists || len(bytes.TrimSpace(data)) < maxLoggedDataSize {
		return json.RawMessage(body)
	}

	delete(response, "data")
	filteredBody, err := json.Marshal(response)
	if err != nil {
		return "[响应日志序列化失败]"
	}
	return json.RawMessage(filteredBody)
}
