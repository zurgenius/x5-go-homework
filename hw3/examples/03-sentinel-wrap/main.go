package main

import (
	"errors"
	"fmt"
)

var ErrEmptyMessage = errors.New("пустое сообщение")

func main() {
	channelErr := fmt.Errorf("email: %w", ErrEmptyMessage)
	err := fmt.Errorf("заказ 42: %w", channelErr)
	textOnly := fmt.Errorf("email: %v", ErrEmptyMessage)
	another := errors.New("пустое сообщение")

	fmt.Println("ошибка:", err)
	fmt.Println("равенство:", err == ErrEmptyMessage)
	fmt.Println("Is через два уровня:", errors.Is(err, ErrEmptyMessage))
	fmt.Println("Unwrap один уровень:", errors.Unwrap(err))
	fmt.Println("одинаковый текст:", another.Error() == ErrEmptyMessage.Error())
	fmt.Println("одинаковая причина:", errors.Is(another, ErrEmptyMessage))
	fmt.Println("Is после текстовой обёртки:", errors.Is(textOnly, ErrEmptyMessage))
}
