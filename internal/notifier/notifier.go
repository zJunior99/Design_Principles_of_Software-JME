package notifier

import (
	"fmt"

	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/domain"
)

// Notifier define el contrato para enviar notificaciones.
// Principio LSP: cualquier implementación debe comportarse de forma predecible.
type Notifier interface {
	Notify(order domain.Order) error
}

// EmailNotifier envía notificaciones por correo electrónico.
type EmailNotifier struct {
	From       string
	SMTPServer string
}

// Notify envía un email de confirmación al cliente.
// Cumple el contrato: retorna error si no puede notificar.
func (e *EmailNotifier) Notify(order domain.Order) error {
	if order.CustomerEmail == "" {
		return fmt.Errorf("no se puede enviar email: el pedido %s no tiene email de cliente", order.ID)
	}
	// En producción, aquí iría la llamada real al servidor SMTP.
	fmt.Printf("[EmailNotifier] Enviando email a %s para pedido %s (total: $%.2f)\n",
		order.CustomerEmail, order.ID, order.FinalTotal)
	return nil
}

// SMSNotifier envía notificaciones por mensaje de texto.
type SMSNotifier struct {
	APIKey string
}

// Notify envía un SMS al cliente.
// Cumple el contrato: retorna error si no puede notificar.
func (s *SMSNotifier) Notify(order domain.Order) error {
	if order.CustomerPhone == "" {
		return fmt.Errorf("no se puede enviar SMS: el pedido %s no tiene teléfono de cliente", order.ID)
	}
	fmt.Printf("[SMSNotifier] Enviando SMS a %s para pedido %s\n",
		order.CustomerPhone, order.ID)
	return nil
}

// MultiNotifier envía notificaciones por múltiples canales.
// Implementa el patrón Composite.
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier crea un notificador que delega en múltiples canales.
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

// Notify intenta notificar por todos los canales.
// Continúa aunque alguno falle, reportando todos los errores.
func (m *MultiNotifier) Notify(order domain.Order) error {
	var errs []string
	for _, n := range m.notifiers {
		if err := n.Notify(order); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("errores de notificación: %v", errs)
	}
	return nil
}

// NoOpNotifier es un notificador que no hace nada (Null Object Pattern).
// Útil para tests y entornos donde las notificaciones no aplican.
type NoOpNotifier struct{}

func (n *NoOpNotifier) Notify(_ domain.Order) error { return nil }
