package notifier_test

import (
	"testing"

	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/domain"
	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/notifier"
)

func orderWith(email, phone string) domain.Order {
	return domain.Order{
		ID:            "ORD-001",
		CustomerEmail: email,
		CustomerPhone: phone,
		FinalTotal:    100.0,
	}
}

func TestEmailNotifier_ValidEmail_Succeeds(t *testing.T) {
	n := &notifier.EmailNotifier{From: "no-reply@tienda.com", SMTPServer: "smtp.tienda.com"}

	err := n.Notify(orderWith("cliente@example.com", ""))
	if err != nil {
		t.Errorf("Notify() error inesperado con email válido: %v", err)
	}
}

func TestEmailNotifier_EmptyEmail_ReturnsError(t *testing.T) {
	n := &notifier.EmailNotifier{From: "no-reply@tienda.com", SMTPServer: "smtp.tienda.com"}

	err := n.Notify(orderWith("", ""))
	if err == nil {
		t.Error("Notify() debería retornar error cuando el email está vacío")
	}
}

func TestSMSNotifier_ValidPhone_Succeeds(t *testing.T) {
	n := &notifier.SMSNotifier{APIKey: "demo-key"}

	err := n.Notify(orderWith("", "+51999000111"))
	if err != nil {
		t.Errorf("Notify() error inesperado con teléfono válido: %v", err)
	}
}

func TestSMSNotifier_EmptyPhone_ReturnsError(t *testing.T) {
	n := &notifier.SMSNotifier{APIKey: "demo-key"}

	err := n.Notify(orderWith("", ""))
	if err == nil {
		t.Error("Notify() debería retornar error cuando el teléfono está vacío")
	}
}

func TestNoOpNotifier_AlwaysSucceeds(t *testing.T) {
	n := &notifier.NoOpNotifier{}

	err := n.Notify(orderWith("", ""))
	if err != nil {
		t.Errorf("NoOpNotifier.Notify() nunca debería retornar error: %v", err)
	}
}

func TestMultiNotifier_AllSucceed(t *testing.T) {
	n := notifier.NewMultiNotifier(
		&notifier.EmailNotifier{From: "a@b.com", SMTPServer: "smtp"},
		&notifier.SMSNotifier{APIKey: "key"},
	)

	err := n.Notify(orderWith("cliente@example.com", "+51999000111"))
	if err != nil {
		t.Errorf("MultiNotifier.Notify() error inesperado: %v", err)
	}
}

func TestMultiNotifier_PartialFailure_ReturnsError(t *testing.T) {
	n := notifier.NewMultiNotifier(
		&notifier.EmailNotifier{From: "a@b.com", SMTPServer: "smtp"},
		&notifier.SMSNotifier{APIKey: "key"},
	)

	// Email ok, SMS falla (sin teléfono)
	err := n.Notify(orderWith("cliente@example.com", ""))
	if err == nil {
		t.Error("MultiNotifier.Notify() debería retornar error cuando algún canal falla")
	}
}

func TestMultiNotifier_Empty_Succeeds(t *testing.T) {
	n := notifier.NewMultiNotifier()

	err := n.Notify(orderWith("", ""))
	if err != nil {
		t.Errorf("MultiNotifier sin canales no debería retornar error: %v", err)
	}
}
