package logger

import (
	"io"
	"log/slog"
)

type LogStream struct {
	InfoStream  *slog.Logger
	DebugStream *slog.Logger
}

var Log LogStream

func InitInfoTextlog(w io.Writer) {
	Log.InfoStream = slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func Info(info string) {
	Log.InfoStream.Info(info)
}

func InitErrorJSONlog(w io.Writer) {
	Log.DebugStream = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func Error(err string) {
	Log.DebugStream.Error(err)
}
