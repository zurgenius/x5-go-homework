package main

import (
	"errors"
	"fmt"
)

type SendError struct {
	Channel string
	Reason  string
}

func (e *SendError) Error() string {
	return fmt.Sprintf("канал %s: %s", e.Channel, e.Reason)
}

func sendPush(deviceID string) error {
	if deviceID == "" {
		return &SendError{
			Channel: "push",
			Reason:  "не задан device id",
		}
	}
	return nil
}

func main() {
	original := sendPush("")
	err := fmt.Errorf("заказ 42: %w", original)
	fmt.Println("ошибка:", err)

	_, ok := err.(*SendError)
	fmt.Println("прямое утверждение типа:", ok)

	var sendErr *SendError
	if errors.As(err, &sendErr) {
		fmt.Println("канал:", sendErr.Channel)
		fmt.Println("причина:", sendErr.Reason)
	}
}
