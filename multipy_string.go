package main

import (
	"fmt"
)

func main() {
	var source string
	var times int
	// гарантируется, что значения корректные
	fmt.Scan(&source, &times)

	// возьмите строку `source` и повторите ее `times` раз
	// запишите результат в `result`
	// ...

	result := ""
	for i := 0; i < times; i++ {
		result += source
	}

	fmt.Println(result)
}
