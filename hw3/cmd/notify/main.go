package main

import (
	"errors"
	"fmt"

	"lesson5/student/practice"
)

func main() {
	email := &practice.EmailNotifier{
		BaseNotifier: practice.BaseNotifier{Channel: "email"},
		Address:      "user@example.com",
	}
	push := &practice.PushNotifier{
		BaseNotifier: practice.BaseNotifier{Channel: "push"},
		DeviceID:     "device-1",
	}
	channels := []practice.Notifier{email, push}

	fmt.Println("1. Успешная отправка")
	err := practice.Dispatch(channels, "Заказ готов")
	fmt.Println("Ошибка:", err)
	fmt.Println("История email:", email.Messages)
	fmt.Println("История push:", push.Messages)

	fmt.Println("\n2. Пустой текст")
	err = practice.Dispatch(channels, "")
	if errors.Is(err, practice.ErrEmptyMessage) {
		fmt.Println("Распознана ошибка пустого сообщения:", err)
	} else {
		fmt.Println("Неожиданный результат:", err)
	}
	fmt.Println("История email:", email.Messages)
	fmt.Println("История push:", push.Messages)

	fmt.Println("\n3. Нет настройки")
	unconfiguredPush := &practice.PushNotifier{BaseNotifier: practice.BaseNotifier{Channel: "push"}}
	err = practice.Dispatch([]practice.Notifier{unconfiguredPush}, "Заказ готов")
	var sendErr *practice.SendError
	if errors.As(err, &sendErr) {
		fmt.Printf("Канал: %s, причина: %s\n", sendErr.Channel, sendErr.Reason)
	} else {
		fmt.Println("Неожиданный результат:", err)
	}

	fmt.Println("\n4. Паника канала")
	broken := &practice.BrokenNotifier{BaseNotifier: practice.BaseNotifier{Channel: "broken"}}
	goodPush := &practice.PushNotifier{
		BaseNotifier: practice.BaseNotifier{Channel: "push"},
		DeviceID:     "device-2",
	}
	failures := practice.SafeDispatch([]practice.Notifier{broken, goodPush}, "Заказ готов")
	for _, failure := range failures {
		if errors.Is(failure, practice.ErrNotifierPanic) {
			fmt.Println("Распознана паника:", failure)
		} else {
			fmt.Println("Ошибка:", failure)
		}
	}
	fmt.Println("История push:", goodPush.Messages)
}
