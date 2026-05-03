package repository_test

import (
	"testing"

	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/domain"
	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/repository"
)

func newOrder(id string) domain.Order {
	return domain.Order{
		ID:         id,
		CustomerID: "CUST-001",
		Total:      500.0,
		Items:      []domain.OrderItem{{ProductID: "P1", Name: "Item", Quantity: 1, UnitPrice: 500}},
	}
}

func TestInMemoryRepository_SaveAndFindByID(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	order := newOrder("ORD-001")

	if err := repo.Save(order); err != nil {
		t.Fatalf("Save() error inesperado: %v", err)
	}

	got, err := repo.FindByID("ORD-001")
	if err != nil {
		t.Fatalf("FindByID() error inesperado: %v", err)
	}
	if got.ID != order.ID {
		t.Errorf("FindByID() = %q; quería %q", got.ID, order.ID)
	}
}

func TestInMemoryRepository_FindByID_NotFound(t *testing.T) {
	repo := repository.NewInMemoryRepository()

	_, err := repo.FindByID("NO-EXISTE")
	if err == nil {
		t.Error("FindByID() con ID inexistente debería retornar error")
	}
}

func TestInMemoryRepository_FindAll(t *testing.T) {
	repo := repository.NewInMemoryRepository()

	_ = repo.Save(newOrder("ORD-001"))
	_ = repo.Save(newOrder("ORD-002"))
	_ = repo.Save(newOrder("ORD-003"))

	orders, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll() error inesperado: %v", err)
	}
	if len(orders) != 3 {
		t.Errorf("FindAll() devolvió %d pedidos; quería 3", len(orders))
	}
}

func TestInMemoryRepository_Save_Overwrites(t *testing.T) {
	repo := repository.NewInMemoryRepository()

	order := newOrder("ORD-001")
	_ = repo.Save(order)

	order.Total = 999.0
	_ = repo.Save(order)

	got, _ := repo.FindByID("ORD-001")
	if got.Total != 999.0 {
		t.Errorf("Save() debería sobreescribir el pedido existente; Total = %.2f; quería 999.00", got.Total)
	}
}

func TestInMemoryRepository_FindAll_Empty(t *testing.T) {
	repo := repository.NewInMemoryRepository()

	orders, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll() en repositorio vacío error inesperado: %v", err)
	}
	if len(orders) != 0 {
		t.Errorf("FindAll() en repositorio vacío devolvió %d pedidos; quería 0", len(orders))
	}
}
