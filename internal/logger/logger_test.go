package logger_test

import (
	"testing"

	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/logger"
)

// ConsoleLogger escribe en stdout — verificamos que no entre en pánico.
func TestConsoleLogger_DoesNotPanic(t *testing.T) {
	log := logger.NewConsoleLogger()

	// Ninguno de estos debe entrar en pánico
	log.Info("mensaje de info", "key", "value")
	log.Warn("mensaje de advertencia", "key", 42)
	log.Error("mensaje de error", "key", true)
}

func TestConsoleLogger_WithNoArgs_DoesNotPanic(t *testing.T) {
	log := logger.NewConsoleLogger()

	log.Info("sin argumentos extra")
	log.Warn("sin argumentos extra")
	log.Error("sin argumentos extra")
}

func TestNoOpLogger_DoesNotPanic(t *testing.T) {
	log := logger.NewNoOpLogger()

	log.Info("ignorado", "key", "value")
	log.Warn("ignorado", "key", "value")
	log.Error("ignorado", "key", "value")
}

// Verifica que NoOpLogger implementa la interfaz Logger.
func TestNoOpLogger_ImplementsInterface(t *testing.T) {
	var _ logger.Logger = logger.NewNoOpLogger()
}

// Verifica que ConsoleLogger implementa la interfaz Logger.
func TestConsoleLogger_ImplementsInterface(t *testing.T) {
	var _ logger.Logger = logger.NewConsoleLogger()
}
