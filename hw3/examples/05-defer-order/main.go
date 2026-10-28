package main

import "fmt"

func work() {
	value := 1
	defer fmt.Println("первый defer, значение:", value)
	defer fmt.Println("второй defer")

	value = 2
	fmt.Println("тело функции, значение:", value)
}

func main() {
	work()
	fmt.Println("вернулись в main")
}
