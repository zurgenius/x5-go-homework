package main

import (
	"errors"
	"fmt"
)

var ErrCalculationPanic = errors.New("паника вычисления")

func calculate(broken bool) int {
	if broken {
		panic("нарушено внутреннее условие")
	}
	return 42
}

func safeCalculate(broken bool) (result int, err error) {
	defer func() {
		if value := recover(); value != nil {
			err = fmt.Errorf("%w: %v", ErrCalculationPanic, value)
		}
	}()

	result = calculate(broken)
	fmt.Println("вычисление закончено")
	return result, nil
}

func main() {
	result, err := safeCalculate(false)
	fmt.Println("обычный вызов:", result, err)

	result, err = safeCalculate(true)
	fmt.Println("сломанный вызов:", result, err)
	fmt.Println("распознали панику:", errors.Is(err, ErrCalculationPanic))
	fmt.Println("main продолжает работу")
}
