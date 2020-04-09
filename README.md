# Easy-i18n

Easy-i18n is a Go package and a command that helps you translate Go programs into multiple languages.

- Supports [pluralized strings](http://cldr.unicode.org/index/cldr-spec/plural-rules) for all 200+ languages in the [Unicode Common Locale Data Repository (CLDR)](http://www.unicode.org/cldr/charts/28/supplemental/language_plural_rules.html).
- Supports strings with similar to [fmt.Sprintf](https://golang.org/pkg/fmt/) format syntax.
- Supports message files of any format (e.g. JSON, TOML, YAML).

# Package i18n

The i18n package provides support for looking up messages according to a set of locale preferences.

```go
package main

import (
	"fmt"
	"os"

	"github.com/mylukin/easy-i18n/i18n"
	"golang.org/x/text/language"
)

func main() {

	i18n.SetLang(language.English)

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

	i18n.Fprintf(os.Stderr, "%s have %d apple.", name, 2, i18n.Plural(
		"%[2]d=1", "%s have an apple.",
		"%[2]d=2", "%s have two apples.",
		"%[2]d>2", "%s have %d apples.",
	))
	fmt.Println()
}
```

# Command easyi18n

The easyi18n command manages message files used by the i18n package.

```
go get -u github.com/mylukin/easy-i18n/easyi18n
easyi18n -h
```

### Extracting messages

Use `easyi18n extract` to extract all i18n.Sprintf function literals in Go source files to a message file for translation.

```json
// ./locales/en.json
{
  "hello world!": "hello world!",
  "hello %s!": "hello %s!",
  "%s has %d apple.": "%s has %d apple.",
  "%s has %d cat.": "%s has %d cat.",
  "%s has %d cats.": "%s has %d cats.",
  "%s have %d apples.": "%s have %d apples.",
  "%s have an apple.": "%s have an apple.",
  "%s have two apples.": "%s have two apples."
}
```

### Translating a new language

1. Create an empty message file for the language that you want to add (e.g. `zh-Hans.json`).
2. Run `easyi18n update ./locales/en.json ./locales/zh-Hans.json` to populate `zh-Hans.json` with the mesages to be translated.

    ```json
    // ./locales/zh-Hans.json
    {
    "hello world!": "hello world!",
    "hello %s!": "hello %s!",
    "%s has %d apple.": "%s has %d apple.",
    "%s has %d cat.": "%s has %d cat.",
    "%s has %d cats.": "%s has %d cats.",
    "%s have %d apples.": "%s have %d apples.",
    "%s have an apple.": "%s have an apple.",
    "%s have two apples.": "%s have two apples."
    }
    ```
3. After `zh-Hans.json` has been translated, run `easyi18n generate ./locales ./catalog.go --pkg main`.
