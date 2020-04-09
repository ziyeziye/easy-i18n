package i18n

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

// Extract messages
func Extract(paths []string, outFile string) error {
	if len(paths) == 0 {
		paths = []string{"."}
	}
	messages := map[string]string{}
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
			buf, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, path, buf, parser.AllErrors)
			if err != nil {
				return err
			}

			fmt.Printf("Extract %+v ...\n", path)
			i18NPackName := i18nPackageName(file)
			// fmt.Printf("i18NPackName %T %+[1]v\n", i18NPackName)
			// ast.Print(fset, file)
			ast.Inspect(file, func(n ast.Node) bool {
				switch v := n.(type) {
				case *ast.CallExpr:
					if fn, ok := v.Fun.(*ast.SelectorExpr); ok {
						var packName string
						if pack, ok := fn.X.(*ast.Ident); ok {
							packName = pack.Name
						}
						funcName := fn.Sel.Name
						// 包名必须相等
						if i18NPackName == packName {
							// 函数名必须相等
							if funcName == "Printf" || funcName == "Sprintf" || funcName == "Fprintf" {
								// 找到字符串
								if str, ok := v.Args[0].(*ast.BasicLit); ok {
									id := strings.Trim(str.Value, `"`)
									if _, ok := messages[id]; !ok {
										messages[id] = id
									}
								}
							}
							if funcName == "Plural" {
								// 找到字符串
								for i := 0; i < len(v.Args); {
									if i++; i >= len(v.Args) {
										break
									}
									if str, ok := v.Args[i].(*ast.BasicLit); ok {
										id := strings.Trim(str.Value, `"`)
										if _, ok := messages[id]; !ok {
											messages[id] = id
										}
									}
									i++
								}
							}
						}
					}
				}
				return true
			})
			return nil
		}); err != nil {
			return err
		}
	}

	var content []byte
	var err error
	of := strings.ToLower(outFile)
	if strings.HasSuffix(of, ".json") {
		content, err = marshal(messages, "json")
	}
	if strings.HasSuffix(of, ".toml") {
		content, err = marshal(messages, "toml")
	}
	if strings.HasSuffix(of, ".yaml") {
		content, err = marshal(messages, "yaml")
	}
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outFile, content, 0664)
	if err != nil {
		return nil
	}
	return nil
}

func i18nPackageName(file *ast.File) string {
	for _, i := range file.Imports {
		if i.Path.Kind == token.STRING && i.Path.Value == `"github.com/mylukin/easy-i18n/i18n"` {
			if i.Name == nil {
				return "i18n"
			}
			return i.Name.Name
		}
	}
	return ""
}

func marshal(v interface{}, format string) ([]byte, error) {
	switch format {
	case "json":
		return json.MarshalIndent(v, "", "  ")
	case "toml":
		var buf bytes.Buffer
		enc := toml.NewEncoder(&buf)
		enc.Indent = ""
		err := enc.Encode(v)
		return buf.Bytes(), err
	case "yaml":
		return yaml.Marshal(v)
	}
	return nil, fmt.Errorf("unsupported format: %s", format)
}
