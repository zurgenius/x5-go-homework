package main

import (
	"errors"
	"fmt"
)

var ErrProductNotFound = errors.New("товар не найден")

type ProductError struct {
	Product string
	Reason  string
}

func (e *ProductError) Error() string {
	return fmt.Sprintf("товар %q: %s", e.Product, e.Reason)
}

func stock(product string) (int, error) {
	if product == "" {
		return 0, &ProductError{Product: product, Reason: "не задано название"}
	}
	products := map[string]int{
		"чай":  5,
		"кофе": 0,
	}
	quantity, ok := products[product]
	if !ok {
		return 0, ErrProductNotFound
	}
	return quantity, nil
}

func checkStock(product string) (int, error) {
	quantity, err := stock(product)
	if err != nil {
		return 0, fmt.Errorf("проверить остаток %q: %w", product, err)
	}
	return quantity, nil
}

func main() {
	for _, product := range []string{"чай", "кофе", "сок", ""} {
		quantity, err := checkStock(product)
		if err == nil {
			fmt.Printf("%q: остаток %d\n", product, quantity)
			continue
		}

		var productErr *ProductError
		switch {
		case errors.Is(err, ErrProductNotFound):
			fmt.Printf("%q: проверьте название товара\n", product)
		case errors.As(err, &productErr):
			fmt.Println("ошибка ввода:", productErr.Reason)
		default:
			fmt.Println("ошибка:", err)
		}
	}
}
