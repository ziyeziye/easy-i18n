package i18n

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Extract messages
func Extract(packName string, paths []string, outFile string) error {
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
			// ignore easy-i18n
			if strings.Index(path, "github.com/mylukin/easy-i18n") > -1 {
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

			// fmt.Printf("Extract %+v ...\n", path)
			i18NPackName := i18nPackageName(file)
			if i18NPackName == "" {
				i18NPackName = packName
			}

			//ast.Print(fset, file)
			ast.Inspect(file, func(n ast.Node) bool {
				switch v := n.(type) {
				case *ast.CallExpr:
					if fn, ok := v.Fun.(*ast.SelectorExpr); ok {
						var packName string
						if pack, ok := fn.X.(*ast.Ident); ok {
							if pack.Obj != nil {
								if as, ok := pack.Obj.Decl.(*ast.AssignStmt); ok {
									if vv, ok := as.Rhs[0].(*ast.CallExpr); ok {
										if vfn, ok := vv.Fun.(*ast.SelectorExpr); ok {
											if vpack, ok := vfn.X.(*ast.Ident); ok {
												packName = vpack.Name
											}
										}
									}
								}
							} else {
								packName = pack.Name
							}

						}
						funcName := fn.Sel.Name
						namePos := fset.Position(fn.Sel.NamePos)

						// Package name must be equal
						if len(packName) > 0 && i18NPackName == packName {
							// Function name must be equal
							if funcName == "Printf" || funcName == "Sprintf" || funcName == "Fprintf" || funcName == "GenerateLocales" {
								id := ""
								domain := ""
								// get domain
								for _, expr := range v.Args {
									if cv, ok := expr.(*ast.CompositeLit); ok {
										if cvt, ok := cv.Type.(*ast.SelectorExpr); ok {
											if pack, ok := cvt.X.(*ast.Ident); ok {
												if pack.Name == "i18n" && cvt.Sel.Name == "Domain" {
													// 读取 domain
													if dv, ok := cv.Elts[0].(*ast.BasicLit); ok {
														domain = strings.Trim(dv.Value, "`")
														domain = strings.Trim(domain, `"`)
													} else if kv, ok := cv.Elts[0].(*ast.KeyValueExpr); ok {
														if dv, ok := kv.Value.(*ast.BasicLit); ok {
															domain = strings.Trim(dv.Value, "`")
															domain = strings.Trim(domain, `"`)
														}
													}
												}
											}
										}
									}
								}
								var fn func(arg ast.Expr) string
								fn = func(arg ast.Expr) string {
									switch value := arg.(type) {
									case *ast.BasicLit:
										id = trim(value.Value)
									case *ast.Ident:
										if value.Obj.Kind == ast.Con {
											if spec, ok := value.Obj.Decl.(*ast.ValueSpec); ok {
												for _, v := range spec.Values {
													val := fn(v)
													if val != "" {
														return val
													}
												}
											}
										}
									}
									return ""
								}
								// Find the string to be translated
								for _, arg := range v.Args {
									val := fn(arg)
									if val != "" {
										id = val
										break
									}
								}

								if id != "" {
									value := id
									if domain != "" {
										id = fmt.Sprintf("%s.%s", domain, id)
									}
									if _, ok := messages[id]; !ok {
										messages[id] = value
									}
									fmt.Printf("Extract %+v %v.%v => %s\n", namePos, packName, funcName, id)
								}
							}
							if funcName == "Plural" {
								// Find the string to be translated
								for i := 0; i < len(v.Args); {
									if i++; i >= len(v.Args) {
										break
									}
									if str, ok := v.Args[i].(*ast.BasicLit); ok {
										id := trim(str.Value)
										if _, ok := messages[id]; !ok {
											messages[id] = id
										}
										fmt.Printf("Extract %+v %v.%v => %s\n", namePos, packName, funcName, id)
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
			fmt.Printf("Extract error: %s\n", err)
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
	err = os.MkdirAll(path.Dir(outFile), os.ModePerm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outFile, content, os.ModePerm)
	if err != nil {
		return err
	}
	fmt.Printf("Extract to %v ...\n", outFile)
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

func trim(text string) string {
	if len(text) > 2 &&
		(text[0] == '"' && text[len(text)-1] == '"') ||
		(text[0] == '`' && text[len(text)-1] == '`') {
		return text[1 : len(text)-1]
	}
	return text
}
