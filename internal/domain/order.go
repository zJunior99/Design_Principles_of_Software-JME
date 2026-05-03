package domain

import "time"

// Order representa un pedido en el sistema.
type Order struct {
	ID            string
	CustomerID    string
	CustomerEmail string
	CustomerPhone string
	Items         []OrderItem
	Total         float64
	FinalTotal    float64
	Status        OrderStatus
	CreatedAt     time.Time
}

// OrderItem representa un producto dentro de un pedido.
type OrderItem struct {
	ProductID string
	Name      string
	Quantity  int
	UnitPrice float64
}

// OrderStatus representa el estado de un pedido.
type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusProcessed  OrderStatus = "processed"
	StatusCancelled  OrderStatus = "cancelled"
)
