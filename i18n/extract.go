package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Extract(paths []string) error {

	if len(paths) == 0 {
		paths = []string{"."}
	}

	for _, path := range paths {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ".go" {
				return nil
			}

			// Don't extract from test files.
			if strings.HasSuffix(path, "_test.go") {
				return nil
			}
			fmt.Printf("%T %+[1]v\n", path)
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}
