package logx

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"
)

func (l *Logger) log(level slog.Level, msg string, args ...any) {
	ctx := context.Background()
	if !l.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])

	record := slog.NewRecord(time.Now(), level, msg, pcs[0])
	record.Add(args...)
	_ = l.Handler().Handle(ctx, record)
}

func Info(msg string, args ...any) {
	L().log(slog.LevelInfo, msg, args...)
}

func Error(msg string, args ...any) {
	L().log(slog.LevelError, msg, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.log(slog.LevelInfo, fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...any) {
	l.log(slog.LevelError, fmt.Sprintf(format, args...))
}
