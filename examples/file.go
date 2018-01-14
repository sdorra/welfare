package main

import (
	"fmt"

	"github.com/sdorra/welfare"
	"github.com/sdorra/welfare/files"
)

const configTemplate = `
one: {{.one}}
two: {{.two}}
three: {{.three}}
`

func main() {
	file := files.NewFileModule("/etc/issue.net", files.Absent)
	run(file, "remove file")

	file = files.NewFileModule("/etc/welfare", files.Directory)
	file.FileMode = 0700
	run(file, "create directory")

	file = files.NewFileModule("/etc/welfare/message", files.File)
	file.Content = "# welfare message file"
	run(file, "create message file")

	copy := files.NewCopyModule("/etc/issue", "/etc/welfare/issue")
	run(copy, "copy issue file")

	context := make(map[string]string)
	context["one"] = "1"
	context["two"] = "1+1"
	context["three"] = "3 (three|drei)"

	template := files.NewTemplateModule("/etc/welfare/config", configTemplate, context)
	run(template, "create config file")
}

func run(module welfare.Module, action string) {
	changed, err := module.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(action, ":", changed)
}
