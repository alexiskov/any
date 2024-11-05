package logger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
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

func InitErrorTemplog(w io.Writer) {
	Log.DebugStream = slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelError}))
}

func Error(err string) {
	Log.DebugStream.Error(err)
}

func CreateTXTlog() (w *os.File) {
	if _, err := os.Stat("logs"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatal(fmt.Errorf("log directory checking error: %w", err))
		}
		if err = os.Mkdir("logs", os.FileMode(int(0777))); err != nil {
			log.Fatal(fmt.Errorf("log dir create error: %w", err))
		}
	}

	w, err := os.OpenFile("logs/temp.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(fmt.Errorf("log-file opening error: %w", err))
	}

	return
}
