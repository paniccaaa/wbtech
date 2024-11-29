package solutions

import "fmt"

// DONE
// Дана структура Human (с произвольным набором полей и методов).
// Реализовать встраивание методов в структуре Action от родительской структуры Human (аналог наследования).

type Human struct {
	age  int
	name string
}

func (h *Human) Greeter() {
	fmt.Printf("Hello! My name is %s\n", h.name)
}

func (h *Human) printAge() int {
	return h.age
}

type Action struct {
	Human
}

func (a *Action) Walk() {
	fmt.Printf("%s is walking, and he is %d", a.name, a.printAge())
}

func Solve1() {
	h := Human{age: 10, name: "Petya"}
	a := Action{h}

	a.Greeter()
	a.Walk()
}
