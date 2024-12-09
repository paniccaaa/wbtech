package solutions

import (
	"fmt"
)

// DONE
// К каким негативным последствиям может привести данный фрагмент кода,
// и как это исправить? Приведите корректный пример реализации.

// var justString string

// func someFunc() {
//   v := createHugeString(1 << 10)
//   justString = v[:100]
// }

// func main() {
//   someFunc()
// }

// Потенциальные проблемы:
// Удержание большого объема памяти из-за среза v[:100] (остальной массив остается в памяти)

// Если createHugeString вернет строку длиной менее 100 символов, попытка взять срез v[:100] вызовет панику из-за выхода за границы.

func createHugeString(size int) string {
	return string(make([]byte, size))
}

func someFunc() string {
	v := createHugeString(1 << 10) // 2^10

	if len(v) < 100 {
		return v
	}

	// copy first 100 symb
	return string(v[:100])
}

func Solve15() {
	fmt.Println(someFunc())
}
