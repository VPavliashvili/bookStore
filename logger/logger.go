package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

var lgr struct {
	*slog.Logger
	uuid string
}

func init() {
	f, err := os.OpenFile("./log.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	wr := io.MultiWriter(os.Stdout, f)

	lgr = struct {
		*slog.Logger
		uuid string
	}{
		Logger: slog.New(slog.NewJSONHandler(wr, nil)),
	}
}

func SetNewUUID(uuid string) {
	lgr.uuid = uuid
}

func wrapUUID(msg string) string {
	if lgr.uuid != "" {
		msg = fmt.Sprintf("UUID: %s, msg: %s", lgr.uuid, msg)
		return msg
	}
	return msg
}

func Info(msg string) {
	msg = wrapUUID(msg)
	lgr.Info(msg)
}

func Error(msg string) {
	msg = wrapUUID(msg)
	lgr.Error(msg)
}

func Warn(msg string) {
	msg = wrapUUID(msg)
	lgr.Warn(msg)
}

func Debug(msg string) {
	msg = wrapUUID(msg)
	lgr.Info(msg)
}
