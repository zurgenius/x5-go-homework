package main

import "fmt"

func broken() {
	defer fmt.Println("defer внутри broken")
	fmt.Println("до panic")
	panic("сломалась учебная реализация")
}

func main() {
	defer fmt.Println("defer внутри main")
	broken()
	fmt.Println("эта строка не выполнится")
}
