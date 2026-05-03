package repository

import (
	"fmt"
	"sync"

	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/domain"
)

// InMemoryRepository es una implementación en memoria del OrderRepository.
// Útil para pruebas y desarrollo local.
type InMemoryRepository struct {
	mu     sync.RWMutex
	orders map[string]domain.Order
}

// NewInMemoryRepository crea un nuevo repositorio en memoria.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		orders: make(map[string]domain.Order),
	}
}

// Save persiste un pedido en memoria.
func (r *InMemoryRepository) Save(order domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return nil
}

// FindByID recupera un pedido por su ID.
func (r *InMemoryRepository) FindByID(id string) (domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, ok := r.orders[id]
	if !ok {
		return domain.Order{}, fmt.Errorf("pedido con ID %q no encontrado", id)
	}
	return order, nil
}

// FindAll devuelve todos los pedidos almacenados.
func (r *InMemoryRepository) FindAll() ([]domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	orders := make([]domain.Order, 0, len(r.orders))
	for _, o := range r.orders {
		orders = append(orders, o)
	}
	return orders, nil
}
