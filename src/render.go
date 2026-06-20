package mangotemplate

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

// convertGoTags transforms <go>...</go> blocks into Go's {{ }} syntax.
// It handles:
//   - <go>expr</go>              → {{expr}}
//   - <go if="cond">...</go>     → {{if cond}}...{{end}}
//   - <go else>...</go>          → {{else}}
//   - <go for="item" in="list">  → {{range $index, $element := list}}
//   - <go func="arg"></go>       → {{func "arg"}} (attribute as function)
func convertGoTags(content string) (string, error) {
	re := regexp.MustCompile(`(?s)<go\s*((?:[^>]|\n)*?)>(.*?)</go>`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		sub := re.FindStringSubmatch(match)
		if len(sub) != 3 {
			return match
		}
		attrs := strings.TrimSpace(sub[1])
		body := strings.TrimSpace(sub[2])

		// No attributes → plain expression
		if attrs == "" {
			return fmt.Sprintf("{{%s}}", body)
		}

		// Parse attributes
		parts := strings.Fields(attrs)
		attrMap := make(map[string]string)
		var tag string
		for _, p := range parts {
			if strings.Contains(p, "=") {
				kv := strings.SplitN(p, "=", 2)
				key := kv[0]
				val := strings.Trim(kv[1], `"`)
				attrMap[key] = val
			} else {
				tag = p
			}
		}

		switch tag {
		case "if":
			if cond, ok := attrMap["if"]; ok {
				return fmt.Sprintf("{{if %s}}%s{{end}}", cond, body)
			}
		case "else":
			return "{{else}}"
		case "for":
			item, ok1 := attrMap["for"]
			in, ok2 := attrMap["in"]
			if ok1 && ok2 {
				return fmt.Sprintf("{{range $index, $element := %s}}%s{{end}}", in, body)
			}
		default:
			// Function call via attribute, e.g. <go upper="name">
			// If exactly one attribute, treat as function call.
			if len(attrMap) == 1 {
				for fname, arg := range attrMap {
					return fmt.Sprintf("{{%s %s}}", fname, arg)
				}
			}
		}
		// fallback: treat body as expression
		return fmt.Sprintf("{{%s}}", body)
	}), nil
}

// executeTemplate runs a compiled template with data and returns the output.
func executeTemplate(tmpl *template.Template, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execution error: %w", err)
	}
	return buf.String(), nil
}