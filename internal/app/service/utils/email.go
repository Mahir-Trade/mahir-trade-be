package utils

import (
	"bytes"
	html "html/template"
)

func MappingValuesToTemplate(values map[string]interface{}, template string) (parsedHtml string, err error) {
	buf := new(bytes.Buffer)
	t, err := html.New("tmpl").Parse(template)
	if err != nil {
		return
	}
	err = t.Execute(buf, values)
	if err != nil {
		return
	}
	parsedHtml = buf.String()
	return
}
