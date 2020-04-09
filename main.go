package main

import (
	"fmt"

	goi18n "github.com/mylukin/easy-i18n/i18n"
)

func main() {

	test("a")

	goi18n.Printf("hello world!")
	fmt.Println()

	name := "Lukin"

	goi18n.Printf("hello %s!", name)
	fmt.Println()

	goi18n.Printf("%s has %d apple.", name, 1)
	fmt.Println()

	goi18n.Printf("%s has %d cat.", name, 2, goi18n.Plural(
		"%[2]d=1", "%s has %d cat.",
		"%[2]d>1", "%s has %d cats.",
	))
	fmt.Println()

	goi18n.Extract([]string{
		".",
	}, "./en.json")
}

func test(a string) {
	fmt.Print(a)
}
