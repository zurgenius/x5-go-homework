package main

import "fmt"

type Task struct {
	Name   string
	Broken bool
}

func run(task Task) {
	if task.Broken {
		panic(task.Name)
	}
	fmt.Println("выполнено:", task.Name)
}

func runWithOuterRecover(tasks []Task) {
	defer func() {
		if value := recover(); value != nil {
			fmt.Println("перехват вокруг цикла:", value)
		}
	}()

	for _, task := range tasks {
		run(task)
	}
}

func safeRun(task Task) (err error) {
	defer func() {
		if value := recover(); value != nil {
			err = fmt.Errorf("сбой задания: %v", value)
		}
	}()
	run(task)
	return nil
}

func main() {
	tasks := []Task{
		{Name: "первое", Broken: true},
		{Name: "второе"},
	}

	fmt.Println("одна граница вокруг цикла")
	runWithOuterRecover(tasks)

	fmt.Println("отдельная граница на каждый вызов")
	for _, task := range tasks {
		if err := safeRun(task); err != nil {
			fmt.Println(err)
		}
	}
}
