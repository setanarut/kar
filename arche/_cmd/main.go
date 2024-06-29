package main

import "fmt"

func a(callback func(int, an any)) {
	callback(3, an)
}

func c(x int, y interface{}) {
	fmt.Println(x, y)
}

func main() {
	a(c(1, 3))

}
