package logx

import "fmt"

func Info(msg string, args ...any) {
	L().Info(msg, args...)
}

func Error(msg string, args ...any) {
	L().Error(msg, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...any) {
	l.Error(fmt.Sprintf(format, args...))
}
