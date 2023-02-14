package logging

import "log"

func LogWarning(l *log.Logger, message string) {
	currPrefix := l.Prefix()
	l.SetPrefix("WARNING: ")
	l.Println(message)
	l.SetPrefix(currPrefix)
}
