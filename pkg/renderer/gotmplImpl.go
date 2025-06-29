package renderer

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Phillezi/gaspecgen/util"
)

type GoTmplRenderer struct{}

func NewGoTemplateRenderer() *GoTmplRenderer {
	return &GoTmplRenderer{}
}

func (r *GoTmplRenderer) Render(templateContent string, data QueryData) (string, error) {
	templ := template.New("sql").Funcs(getTemplateFuncs())
	fields, err := ExtractFields(templateContent, templ)
	if err != nil {
		return "", err
	}

	fieldsToCheck := util.FilterAndTrimPrefix(fields, ".Rows[].")

	var missing string
	if len(data.Rows) > 0 {
		for _, field := range fieldsToCheck {
			if _, exists := data.Rows[0][field]; !exists {
				if missing == "" {
					missing = fmt.Sprintf("Missing from input data, but specified in query:\n%q", field)
				} else {
					missing = fmt.Sprintf("%s,\n%q", missing, field)
				}
			}
		}
	} else if len(fieldsToCheck) > 0 {
		for _, field := range fieldsToCheck {
			if missing == "" {
				missing = fmt.Sprintf("Missing from input data, but specified in query:\n%q", field)
			} else {
				missing = fmt.Sprintf("%s,\n%q", missing, field)
			}
		}
	}
	if missing != "" {
		return "", fmt.Errorf("%s\n[NOTE]: the input data headers gets camelCased", missing)
	}

	tmpl, err := templ.Option("missingkey=error").Parse(templateContent)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
