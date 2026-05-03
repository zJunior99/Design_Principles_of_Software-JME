package service

import (
	"fmt"
	"time"

	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/discount"
	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/domain"
	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/logger"
	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/notifier"
	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/repository"
)

// OrderService orquesta el procesamiento de pedidos.
// Principio DIP: depende de interfaces, no de implementaciones concretas.
// Principio SRP: su única responsabilidad es coordinar el flujo de un pedido.
type OrderService struct {
	repo       repository.OrderRepository
	calculator *discount.Calculator
	notifier   notifier.Notifier
	logger     logger.Logger
}

// NewOrderService construye el servicio con sus dependencias inyectadas.
func NewOrderService(
	repo repository.OrderRepository,
	calc *discount.Calculator,
	notif notifier.Notifier,
	log logger.Logger,
) *OrderService {
	return &OrderService{
		repo:       repo,
		calculator: calc,
		notifier:   notif,
		logger:     log,
	}
}

// ProcessOrder valida, calcula el precio final, persiste y notifica sobre un pedido.
func (s *OrderService) ProcessOrder(order domain.Order) error {
	s.logger.Info("Iniciando procesamiento de pedido", "id", order.ID, "cliente", order.CustomerID)

	// Validación básica
	if err := s.validate(order); err != nil {
		s.logger.Warn("Pedido inválido rechazado", "id", order.ID, "error", err)
		return fmt.Errorf("pedido inválido: %w", err)
	}

	// Calcular precio final con descuentos
	order.FinalTotal = s.calculator.Calculate(order.Total)
	order.Status = domain.StatusProcessed
	order.CreatedAt = time.Now()

	s.logger.Info("Descuentos aplicados",
		"total_original", order.Total,
		"total_final", order.FinalTotal,
		"ahorro", order.Total-order.FinalTotal,
	)

	// Persistir el pedido
	if err := s.repo.Save(order); err != nil {
		s.logger.Error("Error al guardar pedido", "id", order.ID, "error", err)
		return fmt.Errorf("error al guardar pedido %s: %w", order.ID, err)
	}

	// Notificar al cliente (fallo no bloquea el flujo principal)
	if err := s.notifier.Notify(order); err != nil {
		s.logger.Warn("No se pudo notificar al cliente", "id", order.ID, "error", err)
		// Decisión de negocio: el pedido ya fue guardado, la notificación es best-effort
	}

	s.logger.Info("Pedido procesado exitosamente", "id", order.ID, "total_final", order.FinalTotal)
	return nil
}

// GetOrder recupera un pedido por su ID.
func (s *OrderService) GetOrder(id string) (domain.Order, error) {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return domain.Order{}, fmt.Errorf("pedido %s no encontrado: %w", id, err)
	}
	return order, nil
}

// ListOrders devuelve todos los pedidos registrados.
func (s *OrderService) ListOrders() ([]domain.Order, error) {
	orders, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("error al listar pedidos: %w", err)
	}
	return orders, nil
}

// validate verifica que el pedido tenga los datos mínimos requeridos.
// Principio SRP: esta validación es parte de la responsabilidad del servicio.
func (s *OrderService) validate(order domain.Order) error {
	if order.ID == "" {
		return fmt.Errorf("el ID del pedido no puede estar vacío")
	}
	if order.CustomerID == "" {
		return fmt.Errorf("el pedido debe tener un cliente asociado")
	}
	if order.Total <= 0 {
		return fmt.Errorf("el total del pedido debe ser mayor a cero, se recibió: %.2f", order.Total)
	}
	if len(order.Items) == 0 {
		return fmt.Errorf("el pedido debe contener al menos un ítem")
	}
	return nil
}
