package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	
	// host/userororganizationname/project/(dir)/package
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

var (
	ToBe   bool       = false
	MaxInt uint64     = 1<<64 - 1
	z      complex128 = cmplx.Sqrt(-5 + 12i)
)

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


	fmt.Printf("Type: %T Value: %v\n", ToBe, ToBe)
	fmt.Printf("Type: %T Value: %v\n", MaxInt, MaxInt)
	fmt.Printf("Type: %T Value: %v\n", z, z)


	for i:=0; i<3; i++{
		fmt.Print(i," ")
	}
	// POINTERS
	fmt.Println("\nPOINTERS")
	var x int = 1
	var p *int = &x

	// x := 1
	// p := &x
	// p : holds the memory address of x
	// *p : holds the value of x
	fmt.Println("value of p : ", *p)
	fmt.Println("address of p",p)
	
	// t
}
