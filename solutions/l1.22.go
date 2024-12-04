package solutions

import (
	"fmt"
	"math/big"
)

// DONE
// Разработать программу, которая
// перемножает, делит,
// складывает, вычитает
// две числовых переменных a и b, значение которых > 2^20.

func Solve22() {
	add := func(a, b *big.Int) *big.Int {
		return new(big.Int).Add(a, b)
	}

	subs := func(a, b *big.Int) *big.Int {
		return new(big.Int).Sub(a, b)
	}

	multi := func(a, b *big.Int) *big.Int {
		return new(big.Int).Mul(a, b)
	}

	div := func(a, b *big.Int) *big.Int {
		return new(big.Int).Div(a, b)
	}

	a := big.NewInt(1 << 21)
	b := big.NewInt(1 << 23)

	fmt.Println(a, b)

	fmt.Println("add", add(a, b))
	fmt.Println("sub", subs(a, b))
	fmt.Println("mul", multi(a, b))
	fmt.Println("div", div(b, a))
}
