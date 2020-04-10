package main

//go:generate go run main.go extract . ./locales/en.json
//go:generate go run main.go generate ./locales ./catalog.go --pkg=main
//go:generate go build -o example

import (
	"fmt"
	"os"

	"github.com/mylukin/easy-i18n/i18n"
	"golang.org/x/text/language"
)

func main() {

	i18n.SetLang(language.SimplifiedChinese)

	i18n.Printf("hello world!")
	fmt.Println()

	name := "Lukin"

	i18n.Printf("hello %s!", name)
	fmt.Println()

	i18n.Printf("%s has %d cat.", name, 1)
	fmt.Println()

	i18n.Printf("%s has %d cat.", name, 2, i18n.Plural(
		"%[2]d=1", "%s has %d cat.",
		"%[2]d>1", "%s has %d cats.",
	))
	fmt.Println()

	i18n.Fprintf(os.Stderr, "%s have %d apple.", name, 2, i18n.Plural(
		"%[2]d=1", "%s have an apple.",
		"%[2]d=2", "%s have two apples.",
		"%[2]d>2", "%s have %d apples.",
	))
	fmt.Println()
}
