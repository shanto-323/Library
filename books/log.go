package books

import (
	"log/slog"
	"os"
)

var loger = *slog.New(slog.NewJSONHandler(os.Stdout, nil))

type LogLevel = slog.Level

const (
	WARN  LogLevel = slog.LevelWarn
	INFO  LogLevel = slog.LevelInfo
	ERROR LogLevel = slog.LevelError
)

type LogType struct {
	L LogLevel
	M string
	E error
	D any
}

func SlogLogger(l LogType) {
	switch l.L {
	case WARN:
		loger.Warn(l.M, "Issue", l.E, "Data", &l.D)
	case INFO:
		loger.Info(l.M, "Data", l.D)
	case ERROR:
		loger.Error(l.M, "Error", l.E, "More", l.D)
	}
}

func LogError(level LogLevel, path string, err error, more any) {
	SlogLogger(
		LogType{
			L: slog.LevelError,
			M: path,
			E: err,
			D: more,
		},
	)
}

func LogInfo(level LogLevel, path string, data any) {
	SlogLogger(
		LogType{
			L: slog.LevelInfo,
			M: path,
			D: data,
		},
	)
}
