package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"strings"
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

type Vertex struct {
	X int
	Y int
}
type Person struct {
	Name string
	Age  int
}

func main() {
	// function level var statement
	var i, j = 1, 2

	yo := 9.0

	var name, age, salary = "Nikon", 21, 20_000.0
	fmt.Println(i, j, c, python, java, yo)
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

	for i := 0; i < 3; i++ {
		fmt.Print(i, " ")
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
	fmt.Println("address of p", p)

	// type declaration
	type fahrenheit int
	type celsius int

	var f fahrenheit = 32
	var c celsius = 0
	fmt.Println(f, c)

	// strings

	fmt.Println("Stirngs")
	var str string = "hello"
	fmt.Println("%c", str[0])
	fmt.Println(str[:4])
	fmt.Println()
	var sb strings.Builder
	sb.WriteString("hello")
	fmt.Println(sb.String())
	fmt.Println(sb.Cap())
	fmt.Println(sb.Len())
	//structure
	v := Vertex{1, 2}
	fmt.Println(v)

	arr := [5]int{1, 2, 3, 4, 5}
	s2 := arr[:4]   // Creates a slice from arr[1] to arr[3]
	fmt.Println(s2) // Output: [2 3 4]

	// Slices can be appended to
	s2 = append(s2, 6, 7)
	fmt.Println(s2)  // Output: [1 2 3 4 5]
	fmt.Println(arr) // Output: [1 2 3 4 5 4 5]

	// Person
	var pt = Person{"Nikon", 21}
	fmt.Println(pt.Name, pt.Age)

	// go routine
	Concepts()
}
