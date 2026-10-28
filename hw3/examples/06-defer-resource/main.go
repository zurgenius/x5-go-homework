package main

import (
	"errors"
	"fmt"
)

// Учебная модель ресурса: открытие и закрытие только печатают сообщения.
type Resource struct {
	Name string
}

func openResource(name string) (*Resource, error) {
	if name == "" {
		return nil, errors.New("не задано имя ресурса")
	}
	fmt.Println("открыт:", name)
	return &Resource{Name: name}, nil
}

func (r *Resource) Close() {
	fmt.Println("закрыт:", r.Name)
}

func readReport(name string, broken bool) error {
	resource, err := openResource(name)
	if err != nil {
		return err
	}
	defer resource.Close()

	if broken {
		return errors.New("не удалось прочитать отчёт")
	}
	fmt.Println("отчёт прочитан")
	return nil
}

func main() {
	fmt.Println("успешное чтение")
	fmt.Println("результат:", readReport("report", false))
	fmt.Println("ранний возврат")
	fmt.Println("результат:", readReport("report", true))
	fmt.Println("ошибка открытия")
	fmt.Println("результат:", readReport("", false))
}
