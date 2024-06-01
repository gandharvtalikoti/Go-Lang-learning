package main

import (
	"fmt"
	"math"
	"math/rand"
)

// package level var statement 
var c, python, java bool


func add(x, y int) int {
	return x + y

}

func swap(a, b string) (string, string) {
	return b, a
}

func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return
	// this is naked return
	//Naked return statements should be used only in short functions, as with the example shown here.
	// They can harm readability in longer functions.

}
func main() {
	// function level var statement
	var i,j = 1,2

	yo := 9.0

	var name, age, salary = "Nikon", 21, 20_000.0
	fmt.Println(i,j, c, python, java, yo)
	fmt.Println(name, age, salary)

	fmt.Println("hello world", rand.Intn(10))
	fmt.Println("hello world", math.Sqrt(4))
	fmt.Println("hello world", math.Pi)
	fmt.Println("a + b = ", add(2, 3))
	a, b := swap("hello", "world")
	fmt.Println(a, b)
	fmt.Println(split(17))
}
