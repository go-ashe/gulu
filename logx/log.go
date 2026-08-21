package logx

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ashe/gulu/env"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

var (
	once         sync.Once
	globalLogger *Logger
)

type Logger struct {
	*slog.Logger
}

func Init(e env.RunMode) {
	once.Do(func() {
		// 根据环境配置
		opts := &slog.HandlerOptions{
			Level:     getLogLevel(e),
			AddSource: true,
			// 🔧 自定义 source 格式
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.SourceKey && len(groups) == 0 {
					if source, ok := a.Value.Any().(*slog.Source); ok {
						shortFile := filepath.Base(filepath.Dir(source.File)) + "/" + filepath.Base(source.File)
						a.Key = "caller"
						a.Value = slog.StringValue(fmt.Sprintf("%s:%d", shortFile, source.Line))
					}
				}
				return a
			},
		}
		// 设置全局 logger
		sl := slog.New(slog.NewJSONHandler(os.Stdout, opts))
		globalLogger = &Logger{sl}
		// 替换 slog 默认 logger（兼容 log.Printf）
		slog.SetDefault(sl)
	})
}

// 动态获取日志级别（支持热更新）
func getLogLevel(e env.RunMode) slog.Level {
	if e != env.RunModeProd {
		return slog.LevelDebug
	}
	// 生产环境默认 Info，可通过 LevelVar 动态调整
	return slog.LevelInfo
}

// L 获取全局 logger（非请求场景使用）
func L() *Logger {
	if globalLogger == nil {
		return &Logger{slog.Default()}
	}
	return globalLogger
}

func (l *Logger) WithAttrs(args ...any) *Logger {
	if len(args) == 0 {
		return l
	}
	return &Logger{l.With(args...)}
}

// 请求场景：通过 context 传递
type ctxKey struct{}

func FromCtx(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(ctxKey{}).(*Logger); ok {
		return logger
	}
	return globalLogger
}

func ToCtx(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

// WithCtx 扩展 context 中的 logger（推荐用法）
func WithCtx(ctx context.Context, fields ...any) context.Context {
	logger := FromCtx(ctx).WithAttrs(fields...)
	return ToCtx(ctx, logger)
}

// FromGin 辅助函数
func FromGin(c *gin.Context) *Logger {
	if lg, exists := c.Get("logger"); exists {
		if logger, ok := lg.(*Logger); ok {
			return logger
		}
	}
	return globalLogger
}
