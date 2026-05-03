package main

import (
	"fmt"

	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/discount"
	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/domain"
	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/logger"
	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/notifier"
	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/repository"
	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/service"
)

func main() {
	fmt.Println("=== Sistema de Procesamiento de Pedidos ===")
	fmt.Println("Principios SOLID aplicados en Go")

	// --- Composición de dependencias (Principio DIP) ---
	// Este es el único lugar donde se conocen las implementaciones concretas.

	repo := repository.NewInMemoryRepository()

	calc := discount.NewCalculator(
		discount.SeasonalDiscount{Percentage: 10},
		discount.MinimumOrderDiscount{MinimumTotal: 1000, Percentage: 5},
	)

	notif := notifier.NewMultiNotifier(
		&notifier.EmailNotifier{From: "noreply@tienda.com", SMTPServer: "smtp.tienda.com"},
		&notifier.SMSNotifier{APIKey: "demo-key"},
	)

	log := logger.NewConsoleLogger()

	svc := service.NewOrderService(repo, calc, notif, log)

	// --- Procesar pedidos de ejemplo ---

	orders := []domain.Order{
		{
			ID:            "ORD-001",
			CustomerID:    "CUST-123",
			CustomerEmail: "ana@example.com",
			CustomerPhone: "+51999000111",
			Total:         1500.00,
			Items: []domain.OrderItem{
				{ProductID: "P001", Name: "Laptop Gamer", Quantity: 1, UnitPrice: 1500.00},
			},
		},
		{
			ID:            "ORD-002",
			CustomerID:    "CUST-456",
			CustomerEmail: "carlos@example.com",
			CustomerPhone: "+51999000222",
			Total:         800.00,
			Items: []domain.OrderItem{
				{ProductID: "P002", Name: "Teclado Mecánico", Quantity: 2, UnitPrice: 400.00},
			},
		},
	}

	for _, order := range orders {
		fmt.Printf("\n--- Procesando pedido %s ---\n", order.ID)
		if err := svc.ProcessOrder(order); err != nil {
			log.Error("Error procesando pedido", "id", order.ID, "error", err)
			continue
		}

		saved, _ := svc.GetOrder(order.ID)
		fmt.Printf("✅ Pedido %s guardado | Original: $%.2f → Final: $%.2f (ahorro: $%.2f)\n\n",
			saved.ID, saved.Total, saved.FinalTotal, saved.Total-saved.FinalTotal)
	}

	// --- Listar todos los pedidos ---
	allOrders, _ := svc.ListOrders()
	fmt.Printf("\n=== Resumen: %d pedido(s) procesado(s) ===\n", len(allOrders))
	for _, o := range allOrders {
		fmt.Printf("  • %s | Cliente: %s | Total: $%.2f\n", o.ID, o.CustomerID, o.FinalTotal)
	}
}
