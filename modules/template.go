package modules

import (
	"bytes"
	"os"
	"text/template"

	"github.com/pkg/errors"
)

// NewTemplateModule create a new TemplateModule for the target with the given template and context
func NewTemplateModule(target string, template string, context interface{}) *TemplateModule {
	module := &TemplateModule{
		Target:   target,
		Template: template,
		Context:  context,
	}
	module.FileMode = os.FileMode(0644)
	module.UID = os.Getuid()
	module.GID = os.Getegid()
	return module
}

// TemplateModule evaluates the given template, with the context object and ensures that the target file exists with the
// content of the evaluated template.
type TemplateModule struct {
	permissions
	Target   string
	Template string
	Context  interface{}
}

func (module *TemplateModule) Run() (bool, error) {
	tpl, err := template.New(module.Target).Parse(module.Template)
	if err != nil {
		return false, errors.Wrap(err, "failed to parse template")
	}

	var buffer bytes.Buffer
	err = tpl.Execute(&buffer, module.Context)
	if err != nil {
		return false, errors.Wrap(err, "failed to execute template")
	}

	target, err := collectFileInfo(module.Target)
	if err != nil {
		return false, err
	}

	contentChanged, err := ensureContent(target, buffer.String(), module.FileMode)
	if err != nil {
		return false, err
	}

	permissionsChanged, err := ensurePermissions(module.permissions, target)
	if err != nil {
		return false, err
	}

	return contentChanged || permissionsChanged, nil
}
