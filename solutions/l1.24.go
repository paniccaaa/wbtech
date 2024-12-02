package solutions

import (
	"fmt"
	"math"
)

// DONE
// Разработать программу нахождения расстояния между двумя точками,
// которые представлены в виде структуры Point
// с инкапсулированными параметрами x,y и конструктором.

type Point struct {
	x, y float64
}

func (p *Point) Distance(otherPoint *Point) float64 {
	x := otherPoint.x - p.x
	y := otherPoint.y - p.y

	return math.Sqrt(x*x + y*y)
}

func NewPoint(x, y float64) *Point {
	return &Point{
		x: x,
		y: y,
	}
}

func Solve24() {
	point1 := NewPoint(4, 7)
	point2 := NewPoint(2, 6)

	distance := point1.Distance(point2)

	fmt.Println(distance)
}
