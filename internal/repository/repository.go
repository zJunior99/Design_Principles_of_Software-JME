package repository

import "github.com/zJunior99/Design_Principles_of_Software-JME/internal/domain"

// OrderWriter define operaciones de escritura sobre pedidos.
// Principio ISP: interfaces pequeñas y específicas.
type OrderWriter interface {
	Save(order domain.Order) error
}

// OrderReader define operaciones de lectura sobre pedidos.
type OrderReader interface {
	FindByID(id string) (domain.Order, error)
	FindAll() ([]domain.Order, error)
}

// OrderRepository combina lectura y escritura para acceso completo.
type OrderRepository interface {
	OrderWriter
	OrderReader
}
