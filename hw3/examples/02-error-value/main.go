package main

import (
	"errors"
	"fmt"
)

func doubleQuantity(quantity int) (int, error) {
	if quantity < 0 {
		return 0, errors.New("количество не может быть отрицательным")
	}
	return quantity * 2, nil
}

func show(quantity int) {
	result, err := doubleQuantity(quantity)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	fmt.Println("результат:", result)
}

func main() {
	show(3)
	show(-1)
	fmt.Println("main продолжает работу")
}
