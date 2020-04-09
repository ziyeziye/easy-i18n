package main

import (
	"log"
	"os"

	"github.com/Xuanwo/go-locale"
	"github.com/mylukin/easy-i18n/i18n"
	"github.com/urfave/cli/v2"
)

func main() {
	// 探测操作系统语言
	tag, _ := locale.Detect()

	// 设置语言包
	i18n.SetLang(tag)

	appName := "easyi18n"

	app := &cli.App{
		HelpName: appName,
		Name:     appName,
		Usage:    i18n.Sprintf("a tool for managing message translations."),
		Action: func(c *cli.Context) error {
			cli.ShowAppHelp(c)
			return nil
		},

		Commands: []*cli.Command{
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   i18n.Sprintf("merge translations and generate catalog"),
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:      "extract",
				Aliases:   []string{"e"},
				Usage:     i18n.Sprintf("extracts strings to be translated from code"),
				UsageText: i18n.Sprintf("%s extract [path] [outfile]", appName),
				Action: func(c *cli.Context) error {
					path := c.Args().Get(0)
					if len(path) == 0 {
						path = "."
					}
					outFile := c.Args().Get(1)
					if len(outFile) == 0 {
						outFile = "./locales/en.json"
					}
					err := i18n.Extract([]string{
						path,
					}, outFile)
					return err
				},
			},
			{
				Name:      "generate",
				Aliases:   []string{"g"},
				Usage:     i18n.Sprintf("generates code to insert translated messages"),
				UsageText: i18n.Sprintf("%s generate [path] [outfile]", appName),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "pkg",
						Value: "main",
						Usage: i18n.Sprintf("generated go file package name"),
					},
				},
				Action: func(c *cli.Context) error {
					path := c.Args().Get(0)
					if len(path) == 0 {
						path = "./locales"
					}
					outFile := c.Args().Get(1)
					if len(outFile) == 0 {
						outFile = "./catalog.go"
					}
					pkgName := c.String("pkg")
					err := i18n.Generate(
						pkgName,
						[]string{
							path,
						}, outFile)
					return err
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
