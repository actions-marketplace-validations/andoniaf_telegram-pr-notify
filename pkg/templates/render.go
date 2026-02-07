package templates

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/andoniaf/telegram-pr-notify/pkg/events"
)

var funcMap = template.FuncMap{
	"truncate": func(s string, max int) string {
		runes := []rune(s)
		if len(runes) <= max {
			return s
		}
		return string(runes[:max]) + "..."
	},
}

// Render executes a template against the given data.
// If customTpl is non-empty, it is used as the template string.
// Otherwise, a default template is selected based on event type and action.
//
// Note: html/template applies URL-context escaping inside href attributes.
// GitHub URLs are clean ASCII so this is safe for default templates.
// Custom templates with arbitrary URLs containing query params may see URL mangling.
func Render(data *events.TemplateData, customTpl string) (string, error) {
	tplStr := customTpl
	if tplStr == "" {
		tplStr = selectDefault(data)
		if tplStr == "" {
			return "", fmt.Errorf("no template for event %s action %q", data.EventName, data.Action)
		}
	}

	tpl, err := template.New("msg").Funcs(funcMap).Parse(tplStr)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

func selectDefault(data *events.TemplateData) string {
	if data.IsMerged() {
		return defaultTemplates["pull_request:merged"]
	}

	key := data.EventName + ":" + data.Action
	return defaultTemplates[key]
}
