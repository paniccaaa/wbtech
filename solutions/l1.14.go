package solutions

import "fmt"

// DONE
// Разработать программу, которая в рантайме способна определить тип переменной:
// int, string, bool, channel из переменной типа interface{}.

func Solve14() {
	t := func(a interface{}) {
		switch v := a.(type) {
		case int:
			fmt.Println("int", v)
		case string:
			fmt.Println("string", v)
		case bool:
			fmt.Println("bool", v)
		case chan int:
			fmt.Println("chan int", v)
		default:
			fmt.Println("undefined", v)
		}
	}
	var a interface{} = 1
	t(a)

	a = "hello world"
	t(a)

	a = false
	t(a)

	a = make(chan int)
	t(a)

	a = []rune(`hello rune`)
	t(a)
}
