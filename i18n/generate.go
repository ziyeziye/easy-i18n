package i18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

// Generate messages
func Generate(pkgName string, paths []string, outFile string) error {
	if len(paths) == 0 {
		paths = []string{"."}
	}

	goFile, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]*map[string]string{}
	for _, path := range paths {
		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			fileExt := strings.ToLower(filepath.Ext(path))
			if fileExt != ".toml" && fileExt != ".json" && fileExt != ".yaml" {
				return nil
			}

			buf, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			lang := info.Name()[0 : len(info.Name())-5]
			data[lang] = new(map[string]string)
			fmt.Printf("Generate %+v ...\n", path)

			if strings.HasSuffix(fileExt, ".json") {
				err := json.Unmarshal(buf, data[lang])
				if err != nil {
					return err
				}
			}

			if strings.HasSuffix(fileExt, ".yaml") {
				err := yaml.Unmarshal(buf, data[lang])
				if err != nil {
					return err
				}
			}

			if strings.HasSuffix(fileExt, ".toml") {
				_, err := toml.Decode(string(buf), data[lang])
				if err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	err = i18nTmpl.Execute(goFile, struct {
		Data      map[string]*map[string]string
		BackQuote string
		Package   string
	}{
		data,
		"`",
		pkgName,
	})

	return err
}

var funcs = template.FuncMap{
	"funcName": func(lang string) string {
		lang = strings.ReplaceAll(lang, "_", "")
		lang = strings.ReplaceAll(lang, "-", "")
		lang = strings.ToUpper(lang[:1]) + lang[1:]
		return lang
	},
}

var i18nTmpl = template.Must(template.New("i18n").Funcs(funcs).Parse(`package {{.Package}}

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// init
func init() {
	{{- range $k, $v := .Data }}
	init{{ funcName $k }}(language.Make("{{ $k }}"))
	{{- end }}
}

{{- range $k, $v := .Data }}
// init{{ funcName $k }} will init {{ $k }} support.
func init{{ funcName $k }}(tag language.Tag) {
	{{- range $k, $v := $v }}
	message.SetString(tag, "{{$k}}", "{{$v}}")
	{{- end }}
}
{{- end }}
`))
