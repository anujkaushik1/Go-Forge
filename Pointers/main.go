package main

import "fmt"

func main() {

	x := 10

	ptr := &x

	fmt.Println(ptr)  //address
	fmt.Println(*ptr) //value (10)
	println("\n")
	y := 100

	ptrofy := &y

	*ptrofy = 89

	println(y) //89 and not 100

	var ptr3 *int
	ptr3 = ptrofy

	println(*ptr3) //89

	// new

	val1 := 25

	pointer1 := &val1

	val2 := 65

	pointer1 = &val2

	println(*pointer1) // 65
	println(val1)      // 25

	*pointer1 = 189

	println(val1, " ---- ", val2) // 25, 189

	pointer1 = &val1

	println(val1, " ---- ", val2) // 25, 189

}
