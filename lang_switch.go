package main

import (
	"fmt"
)

func main() {
	var code string
	fmt.Scan(&code)

	lang := ""
	switch code {
	case "en":
		lang = "English"
	case "fr":
		lang = "French"
	case "rus":
		fallthrough
	case "ru":
		lang = "Russian"
	default:
		lang = "Unknown"
	}

	fmt.Println(lang)
}
