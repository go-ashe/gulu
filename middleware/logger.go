package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"uuid"

	"github.com/gin-gonic/gin"
	"github.com/go-ashe/gulu/logx"
)

func Logger(c *gin.Context) {
	// 1. 生成请求 ID
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	// 2. 读取 Body（限制 10MB，防止内存溢出）
	body, err := readRequestBody(c.Request, 10*1024*1024)
	if err != nil {
		logx.FromGin(c).Error("读取 Body 失败",
			"error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体过大"})
		c.Abort()
		return
	}

	// 3. 将 Body 重新写回 Request（供后续 Handler 使用）
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// 4. 从全局 logger 派生子 logger（零分配）
	reqLogger := logx.FromGin(c).With(
		"request_id", requestID,
	)
	// 存入 Context，供后续 Gin handler 使用
	c.Set("logger", reqLogger)

	bodyStr := "[二进制/不记录]"
	// 5. 打印 Body（仅 JSON/Form，二进制不打印）
	if shouldLogBody(c.Request.Header.Get("Content-Type")) {
		// 截断超长 Body（最多打印 2000 字符）
		bodyStr = truncateString(string(body), 2000)
	}

	// 定义不需要认证的路径
	skipPaths := []string{
		"map/keepAlive",
		"alarm/sum",
	}

	// 6. 记录请求开始
	needLog := true
	for _, skipPath := range skipPaths {
		if strings.Contains(c.Request.URL.Path, skipPath) {
			needLog = false
			break
		}
	}

	if needLog {
		reqLogger.Info("请求开始",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"remote_addr", c.ClientIP(),
			"body", bodyStr)
	}

	// 7. 继续处理请求
	start := time.Now()
	c.Next() // 执行后续 handler 链
	duration := time.Since(start).String()

	// 8. 获取响应状态码（Gin 自动记录）
	status := c.Writer.Status()
	if status == 0 {
		status = 200
	}

	// 8. 记录请求结束（自动包含所有字段）
	if needLog {
		reqLogger.Info("请求完成",
			slog.String("elapsed", duration),
			slog.Int("status", status),
			slog.Int("body_size", c.Writer.Size()),
		)
	}
}

// readRequestBody 安全读取 Body（带大小限制）
func readRequestBody(r *http.Request, maxSize int64) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	// 限制读取大小
	limitedReader := io.LimitReader(r.Body, maxSize+1) // +1 用于检测溢出
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	// 检测是否超过限制
	if int64(len(body)) > maxSize {
		return nil, io.ErrShortBuffer // 返回特定错误
	}

	return body, nil
}

// shouldLogBody 判断是否需要打印 Body
func shouldLogBody(contentType string) bool {
	// 只打印 JSON 和 Form
	return strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "application/x-www-form-urlencoded") ||
		strings.Contains(contentType, "multipart/form-data")
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + " [TRUNCATED]"
}
