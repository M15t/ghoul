package email

import (
	"bytes"
	"html/template"
	"io/ioutil"
)

// ParseTemplate embeds data into the input template
func ParseTemplate(name string, inputTpl string, data interface{}) (string, error) {
	tpl, err := template.New(name).Parse(inputTpl)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ParseFromPathTemplate embeds data into the input template from path of template file
func ParseFromPathTemplate(path string, data interface{}) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	tplStr := string(b)

	buf := new(bytes.Buffer)

	funcMap := template.FuncMap{
		"safeHTMLAttr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"safeURL": func(s string) template.URL {
			return template.URL(s)
		},
	}

	tpl, err := template.New("").Funcs(funcMap).Option("missingkey=zero").Parse(tplStr)
	if err != nil {
		return "", err
	}

	if err := tpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
