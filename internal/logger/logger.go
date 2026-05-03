package logger

import (
	"fmt"
	"time"
)

// Logger define el contrato para registro de eventos.
// Principio DIP: el servicio depende de esta abstracción, no de log.Printf.
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// ConsoleLogger escribe logs formateados en la salida estándar.
type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) Info(msg string, args ...any) {
	l.log("INFO", msg, args...)
}

func (l *ConsoleLogger) Warn(msg string, args ...any) {
	l.log("WARN", msg, args...)
}

func (l *ConsoleLogger) Error(msg string, args ...any) {
	l.log("ERROR", msg, args...)
}

func (l *ConsoleLogger) log(level, msg string, args ...any) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] %s — %s", timestamp, level, msg)
	for i := 0; i+1 < len(args); i += 2 {
		fmt.Printf(" | %v=%v", args[i], args[i+1])
	}
	fmt.Println()
}

// NoOpLogger descarta todos los mensajes. Ideal para pruebas unitarias.
type NoOpLogger struct{}

func NewNoOpLogger() *NoOpLogger               { return &NoOpLogger{} }
func (n *NoOpLogger) Info(_ string, _ ...any)  {}
func (n *NoOpLogger) Warn(_ string, _ ...any)  {}
func (n *NoOpLogger) Error(_ string, _ ...any) {}
