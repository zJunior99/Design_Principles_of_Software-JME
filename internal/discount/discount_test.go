package discount_test

import (
	"testing"

	"github.com/zJunior99/Design_Principles_of_Software-JME/internal/discount"
)

func TestSeasonalDiscount(t *testing.T) {
	d := discount.SeasonalDiscount{Percentage: 10}

	got := d.Apply(1000)
	want := 900.0

	if got != want {
		t.Errorf("SeasonalDiscount.Apply(1000) = %.2f; quería %.2f", got, want)
	}
}

func TestVIPDiscount(t *testing.T) {
	d := discount.VIPDiscount{}

	got := d.Apply(500)
	want := 400.0

	if got != want {
		t.Errorf("VIPDiscount.Apply(500) = %.2f; quería %.2f", got, want)
	}
}

func TestCouponDiscount(t *testing.T) {
	d := discount.CouponDiscount{Code: "VERANO25", Percentage: 25}

	got := d.Apply(200)
	want := 150.0

	if got != want {
		t.Errorf("CouponDiscount.Apply(200) = %.2f; quería %.2f", got, want)
	}
}

func TestMinimumOrderDiscount_Applied(t *testing.T) {
	d := discount.MinimumOrderDiscount{MinimumTotal: 100, Percentage: 5}

	got := d.Apply(200)
	want := 190.0

	if got != want {
		t.Errorf("MinimumOrderDiscount.Apply(200) = %.2f; quería %.2f", got, want)
	}
}

func TestMinimumOrderDiscount_NotApplied(t *testing.T) {
	d := discount.MinimumOrderDiscount{MinimumTotal: 500, Percentage: 5}

	got := d.Apply(200)
	want := 200.0

	if got != want {
		t.Errorf("MinimumOrderDiscount.Apply(200) con mínimo 500 = %.2f; quería %.2f (sin descuento)", got, want)
	}
}

func TestCalculator_MultipleStrategies(t *testing.T) {
	calc := discount.NewCalculator(
		discount.SeasonalDiscount{Percentage: 10},
		discount.CouponDiscount{Code: "EXTRA5", Percentage: 5},
	)

	// 1000 -> -10% -> 900 -> -5% -> 855
	got := calc.Calculate(1000)
	want := 855.0

	if got != want {
		t.Errorf("Calculator.Calculate(1000) = %.2f; quería %.2f", got, want)
	}
}

func TestCalculator_NoStrategies(t *testing.T) {
	calc := discount.NewCalculator()

	got := calc.Calculate(500)
	want := 500.0

	if got != want {
		t.Errorf("Calculator sin estrategias debería devolver total original; got %.2f, want %.2f", got, want)
	}
}
