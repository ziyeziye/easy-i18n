package main

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// init
func init() {
	initEn(language.Make("en"))
	initZhHans(language.Make("zh-Hans"))
}
// initEn will init en support.
func initEn(tag language.Tag) {
	message.SetString(tag, "%s extract [path] [outfile]", "%s extract [path] [outfile]")
	message.SetString(tag, "%s generate [path] [outfile]", "%s generate [path] [outfile]")
	message.SetString(tag, "a tool for managing message translations.", "a tool for managing message translations.")
	message.SetString(tag, "extracts strings to be translated from code", "extracts strings to be translated from code")
	message.SetString(tag, "generates code to insert translated messages", "generates code to insert translated messages")
	message.SetString(tag, "merge translations and generate catalog", "merge translations and generate catalog")
}
// initZhHans will init zh-Hans support.
func initZhHans(tag language.Tag) {
	message.SetString(tag, "%s extract [path] [outfile]", "%s 提取 [路径] [输出文件]")
	message.SetString(tag, "%s generate [path] [outfile]", "%s 生成 [路径] [输出文件]")
	message.SetString(tag, "a tool for managing message translations.", "用于管理消息翻译的工具。")
	message.SetString(tag, "extracts strings to be translated from code", "从代码中提取要翻译的字符串")
	message.SetString(tag, "generates code to insert translated messages", "生成代码以插入翻译后的消息")
	message.SetString(tag, "merge translations and generate catalog", "合并翻译并生成目录")
}
