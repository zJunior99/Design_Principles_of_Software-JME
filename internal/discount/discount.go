package discount

import "fmt"

// Strategy define el contrato para cualquier tipo de descuento.
// Principio OCP: abierto para extensión, cerrado para modificación.
type Strategy interface {
	Apply(total float64) float64
	Name() string
}

// Calculator aplica múltiples estrategias de descuento en secuencia.
type Calculator struct {
	strategies []Strategy
}

// NewCalculator crea un calculador con las estrategias proporcionadas.
func NewCalculator(strategies ...Strategy) *Calculator {
	return &Calculator{strategies: strategies}
}

// Calculate aplica todos los descuentos sobre el total dado.
func (c *Calculator) Calculate(total float64) float64 {
	result := total
	for _, s := range c.strategies {
		result = s.Apply(result)
	}
	return result
}

// AppliedDiscounts retorna el detalle de cada descuento aplicado.
func (c *Calculator) AppliedDiscounts(total float64) []string {
	details := make([]string, 0, len(c.strategies))
	current := total
	for _, s := range c.strategies {
		after := s.Apply(current)
		saved := current - after
		details = append(details, fmt.Sprintf("%s — ahorro: $%.2f", s.Name(), saved))
		current = after
	}
	return details
}

// NoDiscount es una estrategia que no aplica ningún descuento.
// Útil como valor nulo (Null Object Pattern).
type NoDiscount struct{}

func (n NoDiscount) Apply(total float64) float64 { return total }
func (n NoDiscount) Name() string                { return "Sin descuento" }

// SeasonalDiscount aplica un descuento porcentual de temporada.
type SeasonalDiscount struct {
	Percentage float64
}

func (s SeasonalDiscount) Apply(total float64) float64 {
	return total * (1 - s.Percentage/100)
}

func (s SeasonalDiscount) Name() string {
	return fmt.Sprintf("Descuento de temporada (%.0f%%)", s.Percentage)
}

// VIPDiscount aplica un 20% de descuento a clientes VIP.
type VIPDiscount struct{}

func (v VIPDiscount) Apply(total float64) float64 {
	return total * 0.80
}

func (v VIPDiscount) Name() string {
	return "Descuento VIP (20%)"
}

// CouponDiscount aplica un descuento mediante un código de cupón.
type CouponDiscount struct {
	Code       string
	Percentage float64
}

func (c CouponDiscount) Apply(total float64) float64 {
	return total * (1 - c.Percentage/100)
}

func (c CouponDiscount) Name() string {
	return fmt.Sprintf("Cupón %s (%.0f%%)", c.Code, c.Percentage)
}

// MinimumOrderDiscount aplica un descuento solo si el total supera un mínimo.
type MinimumOrderDiscount struct {
	MinimumTotal float64
	Percentage   float64
}

func (m MinimumOrderDiscount) Apply(total float64) float64 {
	if total >= m.MinimumTotal {
		return total * (1 - m.Percentage/100)
	}
	return total
}

func (m MinimumOrderDiscount) Name() string {
	return fmt.Sprintf("Descuento por compra mínima $%.0f (%.0f%%)", m.MinimumTotal, m.Percentage)
}
