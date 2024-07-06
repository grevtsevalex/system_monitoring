package logger

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
)

const (
	InfoLevel = iota
	DebugLevel
	WarnLevel
	ErrorLevel
	InfoLevelLiteral  = "INFO"
	DebugLevelLiteral = "DEBUG"
	WarnLevelLiteral  = "WARN"
	ErrorLevelLiteral = "ERROR"
)

var (
	ErrWriteToStorage = errors.New("write to storage")
	ErrWriteZeroBytes = errors.New("write zero bytes to storage")
	loggerLevels      = map[string]int{
		InfoLevelLiteral: InfoLevel, DebugLevelLiteral: DebugLevel,
		WarnLevelLiteral: WarnLevel, ErrorLevelLiteral: ErrorLevel,
	}
)

// Logger модель логгера.
type Logger struct {
	level    string
	storages []io.Writer
}

// LogMsg модель сообщения.
type LogMsg struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

// New конструктор логгера.
func New(level string, storages ...io.Writer) *Logger {
	return &Logger{level: strings.ToUpper(level), storages: storages}
}

// Log логирование.
func (l Logger) Log(msg string) {
	switch l.level {
	case ErrorLevelLiteral:
		l.Error(msg)
	case DebugLevelLiteral:
		l.Debug(msg)
	case WarnLevelLiteral:
		l.Warn(msg)
	default:
		l.Info(msg)
	}
}

// Info логирование информации.
func (l Logger) Info(msg string) {
	err := l.write(msg, InfoLevelLiteral)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}

// Warn логирование предупреждений.
func (l Logger) Warn(msg string) {
	err := l.write(msg, WarnLevelLiteral)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}

// Error логирование ошибки.
func (l Logger) Error(msg string) {
	err := l.write(msg, ErrorLevelLiteral)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}

// Debug логирование дебага.
func (l Logger) Debug(msg string) {
	err := l.write(msg, DebugLevelLiteral)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}

// write запись данных в хранилище.
func (l Logger) write(msg string, level string) error {
	if loggerLevels[level] < loggerLevels[l.level] {
		return nil
	}
	lm := LogMsg{
		Level:   level,
		Message: msg,
	}
	s, err := json.Marshal(lm)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	writer := io.MultiWriter(l.storages...)

	n, err := writer.Write(append(s, 10))
	if err != nil {
		return ErrWriteToStorage
	}

	if n == 0 {
		return ErrWriteZeroBytes
	}

	return nil
}
