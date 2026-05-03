package service_test

import (
	"fmt"
	"testing"

	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/discount"
	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/domain"
	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/logger"
	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/notifier"
	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/repository"
	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/service"
)

// --- Helpers para construir pedidos de prueba ---

func newValidOrder(id string) domain.Order {
	return domain.Order{
		ID:            id,
		CustomerID:    "CUST-001",
		CustomerEmail: "cliente@example.com",
		CustomerPhone: "+51999000111",
		Total:         1000.00,
		Items: []domain.OrderItem{
			{ProductID: "PROD-1", Name: "Laptop", Quantity: 1, UnitPrice: 1000.00},
		},
	}
}

func newService(strategies ...discount.Strategy) *service.OrderService {
	return service.NewOrderService(
		repository.NewInMemoryRepository(),
		discount.NewCalculator(strategies...),
		&notifier.NoOpNotifier{},
		logger.NewNoOpLogger(),
	)
}

// --- Tests de procesamiento ---

func TestProcessOrder_ValidOrder_Succeeds(t *testing.T) {
	svc := newService()
	order := newValidOrder("ORD-001")

	err := svc.ProcessOrder(order)

	if err != nil {
		t.Fatalf("ProcessOrder() error inesperado: %v", err)
	}
}

func TestProcessOrder_AppliesSeasonalDiscount(t *testing.T) {
	svc := newService(discount.SeasonalDiscount{Percentage: 10})
	order := newValidOrder("ORD-002")
	order.Total = 1500.00

	if err := svc.ProcessOrder(order); err != nil {
		t.Fatalf("ProcessOrder() error inesperado: %v", err)
	}

	saved, err := svc.GetOrder("ORD-002")
	if err != nil {
		t.Fatalf("GetOrder() error inesperado: %v", err)
	}

	want := 1350.00 // 1500 - 10%
	if saved.FinalTotal != want {
		t.Errorf("FinalTotal = %.2f; quería %.2f", saved.FinalTotal, want)
	}
}

func TestProcessOrder_AppliesMultipleDiscounts(t *testing.T) {
	svc := newService(
		discount.SeasonalDiscount{Percentage: 10},
		discount.CouponDiscount{Code: "EXTRA5", Percentage: 5},
	)
	order := newValidOrder("ORD-003")
	order.Total = 1000.00

	if err := svc.ProcessOrder(order); err != nil {
		t.Fatalf("ProcessOrder() error inesperado: %v", err)
	}

	saved, _ := svc.GetOrder("ORD-003")

	// 1000 -> -10% -> 900 -> -5% -> 855
	want := 855.0
	if saved.FinalTotal != want {
		t.Errorf("FinalTotal = %.2f; quería %.2f", saved.FinalTotal, want)
	}
}

func TestProcessOrder_SetsStatusToProcessed(t *testing.T) {
	svc := newService()
	order := newValidOrder("ORD-004")

	_ = svc.ProcessOrder(order)

	saved, _ := svc.GetOrder("ORD-004")
	if saved.Status != domain.StatusProcessed {
		t.Errorf("Status = %q; quería %q", saved.Status, domain.StatusProcessed)
	}
}

// --- Tests de validación ---

func TestProcessOrder_EmptyID_ReturnsError(t *testing.T) {
	svc := newService()
	order := newValidOrder("")

	err := svc.ProcessOrder(order)

	if err == nil {
		t.Error("ProcessOrder() con ID vacío debería retornar error")
	}
}

func TestProcessOrder_EmptyCustomerID_ReturnsError(t *testing.T) {
	svc := newService()
	order := newValidOrder("ORD-005")
	order.CustomerID = ""

	err := svc.ProcessOrder(order)

	if err == nil {
		t.Error("ProcessOrder() sin CustomerID debería retornar error")
	}
}

func TestProcessOrder_ZeroTotal_ReturnsError(t *testing.T) {
	svc := newService()
	order := newValidOrder("ORD-006")
	order.Total = 0

	err := svc.ProcessOrder(order)

	if err == nil {
		t.Error("ProcessOrder() con total 0 debería retornar error")
	}
}

func TestProcessOrder_NoItems_ReturnsError(t *testing.T) {
	svc := newService()
	order := newValidOrder("ORD-007")
	order.Items = nil

	err := svc.ProcessOrder(order)

	if err == nil {
		t.Error("ProcessOrder() sin ítems debería retornar error")
	}
}

// --- Tests de notificación tolerante a fallos ---

func TestProcessOrder_NotificationFailure_DoesNotBlockOrder(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	calc := discount.NewCalculator()
	failingNotifier := &AlwaysFailNotifier{}
	log := logger.NewNoOpLogger()

	svc := service.NewOrderService(repo, calc, failingNotifier, log)
	order := newValidOrder("ORD-008")

	// El pedido debe procesarse aunque la notificación falle
	err := svc.ProcessOrder(order)
	if err != nil {
		t.Errorf("ProcessOrder() no debería fallar por error de notificación: %v", err)
	}

	// El pedido debe estar guardado igualmente
	_, err = svc.GetOrder("ORD-008")
	if err != nil {
		t.Errorf("El pedido debería estar guardado aunque la notificación fallara: %v", err)
	}
}

// --- Tests de consulta ---

func TestGetOrder_NotFound_ReturnsError(t *testing.T) {
	svc := newService()

	_, err := svc.GetOrder("NO-EXISTE")

	if err == nil {
		t.Error("GetOrder() con ID inexistente debería retornar error")
	}
}

func TestListOrders_ReturnsAllProcessedOrders(t *testing.T) {
	svc := newService()

	for i := 1; i <= 3; i++ {
		_ = svc.ProcessOrder(newValidOrder(fmt.Sprintf("ORD-%03d", i)))
	}

	orders, err := svc.ListOrders()
	if err != nil {
		t.Fatalf("ListOrders() error inesperado: %v", err)
	}
	if len(orders) != 3 {
		t.Errorf("ListOrders() devolvió %d pedidos; quería 3", len(orders))
	}
}

// --- Mocks auxiliares ---

// AlwaysFailNotifier simula un notificador que siempre falla.
type AlwaysFailNotifier struct{}

func (f *AlwaysFailNotifier) Notify(_ domain.Order) error {
	return fmt.Errorf("servicio de notificación no disponible")
}
