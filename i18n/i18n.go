package i18n

import (
	"io"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var p *message.Printer

// PluralRule is rule
type PluralRule struct {
	Pos   int
	Expr  string
	Value int
	Text  string
}

func init() {
	// init use English
	p = message.NewPrinter(language.English)
}

// SetLang set language
func SetLang(lang language.Tag) {
	p = message.NewPrinter(lang)
}

// Printf is like fmt.Printf, but using language-specific formatting.
func Printf(format string, args ...interface{}) {
	format, args = preArgs(format, args...)
	p.Printf(format, args...)
}

// Sprintf formats according to a format specifier and returns the resulting string.
func Sprintf(format string, args ...interface{}) string {
	format, args = preArgs(format, args...)
	return p.Sprintf(format, args...)
}

// Fprintf is like fmt.Fprintf, but using language-specific formatting.
func Fprintf(w io.Writer, key message.Reference, a ...interface{}) (n int, err error) {
	format, args := preArgs(key.(string), a...)
	key = message.Reference(format)
	return p.Fprintf(w, key, args...)
}

// preArgs
func preArgs(format string, args ...interface{}) (string, []interface{}) {
	length := len(args)
	if length > 0 {
		lastArg := args[length-1]
		switch lastArg.(type) {
		case []PluralRule:
			rules := lastArg.([]PluralRule)
			// parse rule
			for _, rule := range rules {
				curPosVal := args[rule.Pos-1].(int)
				if (rule.Expr == "=" && curPosVal == rule.Value) || (rule.Expr == ">" && curPosVal > rule.Value) {
					format = rule.Text
					break
				}
			}
			args = args[0:strings.Count(format, "%")]
		}
	}
	return format, args
}

// Plural is parse to rule
func Plural(cases ...interface{}) []PluralRule {
	rules := []PluralRule{}
	// %[1]d=1, %[1]d>1
	re := regexp.MustCompile(`\[(\d+)\][^=>]\s*(\=|\>)\s*(\d+)$`)
	for i := 0; i < len(cases); {
		expr := cases[i].(string)
		if i++; i >= len(cases) {
			return rules
		}
		text := cases[i].(string)
		// cannot match continue
		if !re.MatchString(expr) {
			continue
		}
		matches := re.FindStringSubmatch(expr)
		pos, _ := strconv.Atoi(matches[1])
		value, _ := strconv.Atoi(matches[3])
		rules = append(rules, PluralRule{
			Pos:   pos,
			Expr:  matches[2],
			Value: value,
			Text:  text,
		})
		i++
	}
	return rules
}
