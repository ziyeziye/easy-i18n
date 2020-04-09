package main

import (
	"fmt"

	"github.com/mylukin/easy-i18n/i18n"
)

func main() {

	test("a")

	i18n.Printf("hello world!")
	fmt.Println()

	name := "Lukin"

	i18n.Printf("hello %s!", name)
	fmt.Println()

	i18n.Printf("%s has %d apple.", name, 1)
	fmt.Println()

	i18n.Printf("%s has %d cat.", name, 2, i18n.Plural(
		"%[2]d=1", "%s has %d cat.",
		"%[2]d>1", "%s has %d cats.",
	))
	fmt.Println()

	i18n.Extract([]string{
		".",
	}, "./en.json")
}

func test(a string) {
	fmt.Print(a)
}
