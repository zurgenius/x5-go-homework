package main

import "fmt"

type Notifier interface {
	Send(message string) error
}

type BaseNotifier struct {
	Channel string
}

func (b BaseNotifier) Label(message string) string {
	return "[" + b.Channel + "] " + message
}

type EmailNotifier struct {
	BaseNotifier
	Address  string
	Messages []string
}

type PushNotifier struct {
	BaseNotifier
	DeviceID string
	Messages []string
}

func (n *EmailNotifier) Send(message string) error {
	n.Messages = append(n.Messages, n.Label(message))
	return nil
}

func (n *PushNotifier) Send(message string) error {
	n.Messages = append(n.Messages, n.Label(message))
	return nil
}

func Dispatch(all []Notifier, message string) error {
	for _, n := range all {
		if err := n.Send(message); err != nil {
			return err
		}
	}
	return nil
}

func EmailAddress(n Notifier) (string, bool) {
	email, ok := n.(*EmailNotifier)
	if !ok || email == nil {
		return "", false
	}
	return email.Address, true
}

func Describe(n Notifier) string {
	switch v := n.(type) {
	case *EmailNotifier:
		if v == nil {
			return "unknown"
		}
		return "email: " + v.Address
	case *PushNotifier:
		if v == nil {
			return "unknown"
		}
		return "push: " + v.DeviceID
	default:
		return "unknown"
	}
}

var _ Notifier = (*EmailNotifier)(nil)

var _ Notifier = (*PushNotifier)(nil)

func main() {
	email := &EmailNotifier{
		BaseNotifier: BaseNotifier{Channel: "email"},
		Address:      "student@example.com",
	}
	push := &PushNotifier{
		BaseNotifier: BaseNotifier{Channel: "push"},
		DeviceID:     "phone-1",
	}
	err := Dispatch([]Notifier{email, push}, "Заказ готов")
	fmt.Println("ошибка:", err)
	fmt.Println("email:", email.Messages)
	fmt.Println("push:", push.Messages)
}
