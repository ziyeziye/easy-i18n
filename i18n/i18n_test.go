package i18n

import (
	"fmt"
	"testing"
)

func TestPlural(t *testing.T) {
	p := Plural(
		"%[2]d=1", "%s has %d cat.",
		"%[2]d>1", "%s has %d cats.",
	)

	fmt.Printf("%T => %+[1]v\n", p)

}
