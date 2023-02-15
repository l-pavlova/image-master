package logging

import (
	"log"
	"os"
)

type Level string

const (
	Info    Level = "info"
	Warning Level = "warning"
	Error   Level = "error"
)

const (
	Ldate     = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                     // the time in the local time zone: 01:23:23
	LstdFlags = Ldate | Ltime // initial values for the standard logger
)

type ImageMasterLogger struct {
	warningLogger *log.Logger
	infoLogger    *log.Logger
	errorLogger   *log.Logger
}

// initialize a new logger with log levels |info, warning, error|
func NewImageMasterLogger() *ImageMasterLogger {
	return &ImageMasterLogger{
		infoLogger:    log.New(os.Stdout, "INFO: ", LstdFlags),
		warningLogger: log.New(os.Stdout, "WARNING: ", LstdFlags),
		errorLogger:   log.New(os.Stdout, "ERROR: ", LstdFlags),
	}
}

// Log function accepts as first string level = ("info", "warning", "error" and the corresponding implementation logs the message and the additional params from it to stdout)
func (i *ImageMasterLogger) Log(level Level, message string, a ...any) {
	switch level {
	case Info:
		i.infoLogger.Println(message, a)
		return
	case Warning:
		i.warningLogger.Println(message, a)
		return
	case Error:
		i.errorLogger.Println(message, a)
		return
	}
}
