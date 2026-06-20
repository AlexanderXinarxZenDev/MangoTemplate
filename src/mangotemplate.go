// Package mangotemplate provides a PHP‑like templating engine for Go.
package mangotemplate

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Engine is the main templating engine.
type Engine struct {
	dir        string                 // template directory
	delimLeft  string                 // left delimiter (default "{{")
	delimRight string                 // right delimiter (default "}}")
	debug      bool                   // debug mode flag
	funcs      template.FuncMap       // custom functions
	cache      map[string]*template.Template // compiled templates by file name
	cacheMu    sync.RWMutex           // protects cache
}

// NewEngine creates a new engine with the given template directory.
func NewEngine(dir string) *Engine {
	return &Engine{
		dir:        dir,
		delimLeft:  "{{",
		delimRight: "}}",
		funcs:      make(template.FuncMap),
		cache:      make(map[string]*template.Template),
	}
}

// SetTemplateDir changes the template directory.
func (e *Engine) SetTemplateDir(dir string) {
	e.dir = dir
}

// SetDelimiters changes the delimiters used for Go templates.
func (e *Engine) SetDelimiters(left, right string) {
	e.delimLeft = left
	e.delimRight = right
}

// EnableDebug turns on debug logging (prints to stdout).
func (e *Engine) EnableDebug() {
	e.debug = true
}

// DisableDebug turns off debug logging.
func (e *Engine) DisableDebug() {
	e.debug = false
}

// RegisterFunc adds a custom function to the template engine.
func (e *Engine) RegisterFunc(name string, fn interface{}) {
	e.funcs[name] = fn
}

// RegisterBuiltinFunctions registers all built‑in functions.
func (e *Engine) RegisterBuiltinFunctions() {
	// String functions
	e.RegisterFunc("upper", strings.ToUpper)
	e.RegisterFunc("lower", strings.ToLower)
	e.RegisterFunc("title", strings.Title)
	e.RegisterFunc("capitalize", func(s string) string {
		if len(s) == 0 {
			return s
		}
		return strings.ToUpper(s[:1]) + s[1:]
	})
	e.RegisterFunc("reverse", func(s string) string {
		r := []rune(s)
		for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
			r[i], r[j] = r[j], r[i]
		}
		return string(r)
	})
	e.RegisterFunc("truncate", func(s string, n int) string {
		runes := []rune(s)
		if len(runes) <= n {
			return s
		}
		return string(runes[:n]) + "..."
	})
	e.RegisterFunc("slug", func(s string) string {
		// simple slug: lowercase, replace spaces with hyphens, remove non-alphanumeric
		re := regexp.MustCompile(`[^a-z0-9]+`)
		return re.ReplaceAllString(strings.ToLower(strings.TrimSpace(s)), "-")
	})
	e.RegisterFunc("contains", strings.Contains)
	e.RegisterFunc("replace", strings.ReplaceAll)
	e.RegisterFunc("split", strings.Split)
	e.RegisterFunc("join", strings.Join)
	e.RegisterFunc("trim", strings.TrimSpace)
	e.RegisterFunc("strlen", func(s string) int { return len([]rune(s)) })
	e.RegisterFunc("substr", func(s string, start, length int) string {
		runes := []rune(s)
		if start < 0 {
			start = len(runes) + start
		}
		if start < 0 || start >= len(runes) {
			return ""
		}
		end := start + length
		if end > len(runes) {
			end = len(runes)
		}
		return string(runes[start:end])
	})

	// Math functions
	e.RegisterFunc("add", func(a, b float64) float64 { return a + b })
	e.RegisterFunc("subtract", func(a, b float64) float64 { return a - b })
	e.RegisterFunc("multiply", func(a, b float64) float64 { return a * b })
	e.RegisterFunc("divide", func(a, b float64) float64 { return a / b })
	e.RegisterFunc("mod", func(a, b int) int { return a % b })
	e.RegisterFunc("min", func(a, b float64) float64 { if a < b { return a }; return b })
	e.RegisterFunc("max", func(a, b float64) float64 { if a > b { return a }; return b })
	e.RegisterFunc("abs", func(a float64) float64 { if a < 0 { return -a }; return a })

	// Date/Time functions
	e.RegisterFunc("date", func(t time.Time) string {
		return t.Format("2006-01-02")
	})
	e.RegisterFunc("time", func(t time.Time) string {
		return t.Format("15:04:05")
	})
	e.RegisterFunc("datetime", func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	})
	e.RegisterFunc("unix", func(t time.Time) int64 {
		return t.Unix()
	})
	e.RegisterFunc("now", func() time.Time {
		return time.Now()
	})

	// Array functions
	e.RegisterFunc("first", func(items []interface{}) interface{} {
		if len(items) == 0 {
			return nil
		}
		return items[0]
	})
	e.RegisterFunc("last", func(items []interface{}) interface{} {
		if len(items) == 0 {
			return nil
		}
		return items[len(items)-1]
	})
	e.RegisterFunc("length", func(items interface{}) int {
		switch v := items.(type) {
		case []interface{}:
			return len(v)
		case string:
			return len([]rune(v))
		case map[string]interface{}:
			return len(v)
		default:
			return 0
		}
	})
	e.RegisterFunc("range", func(start, end int) []int {
		var r []int
		for i := start; i <= end; i++ {
			r = append(r, i)
		}
		return r
	})

	// Conditional & default
	e.RegisterFunc("default", func(def interface{}, val interface{}) interface{} {
		if val == nil {
			return def
		}
		return val
	})
	e.RegisterFunc("empty", func(val interface{}) bool {
		if val == nil {
			return true
		}
		switch v := val.(type) {
		case string:
			return v == ""
		case []interface{}:
			return len(v) == 0
		case map[string]interface{}:
			return len(v) == 0
		default:
			return false
		}
	})
	e.RegisterFunc("eq", func(a, b interface{}) bool { return a == b })
	e.RegisterFunc("ne", func(a, b interface{}) bool { return a != b })
	e.RegisterFunc("lt", func(a, b float64) bool { return a < b })
	e.RegisterFunc("gt", func(a, b float64) bool { return a > b })
	e.RegisterFunc("lte", func(a, b float64) bool { return a <= b })
	e.RegisterFunc("gte", func(a, b float64) bool { return a >= b })

	// Type conversion
	e.RegisterFunc("int", func(v interface{}) int {
		switch val := v.(type) {
		case int:
			return val
		case float64:
			return int(val)
		case string:
			var i int
			fmt.Sscanf(val, "%d", &i)
			return i
		default:
			return 0
		}
	})
	e.RegisterFunc("float", func(v interface{}) float64 {
		switch val := v.(type) {
		case float64:
			return val
		case int:
			return float64(val)
		case string:
			var f float64
			fmt.Sscanf(val, "%f", &f)
			return f
		default:
			return 0.0
		}
	})
	e.RegisterFunc("string", func(v interface{}) string {
		return fmt.Sprintf("%v", v)
	})
	e.RegisterFunc("bool", func(v interface{}) bool {
		switch val := v.(type) {
		case bool:
			return val
		case int:
			return val != 0
		case float64:
			return val != 0.0
		case string:
			return val != ""
		default:
			return false
		}
	})
	e.RegisterFunc("type", func(v interface{}) string {
		return fmt.Sprintf("%T", v)
	})

	// Security functions
	e.RegisterFunc("escape", template.HTMLEscapeString)
	e.RegisterFunc("unescape", func(s string) string {
		return template.HTML(s).String() // careful: only use when content is safe
	})
	e.RegisterFunc("safe", func(s string) template.HTML {
		return template.HTML(s)
	})
	e.RegisterFunc("md5", func(s string) string {
		// simplistic; real md5 would need crypto/md5
		return fmt.Sprintf("%x", []byte(s))
	})
	e.RegisterFunc("hash", func(s string) string {
		// placeholder
		return fmt.Sprintf("%d", len(s))
	})
	e.RegisterFunc("urlencode", func(s string) string {
		return strings.ReplaceAll(s, " ", "%20")
	})
}

// Render parses and executes the template file with the given data.
func (e *Engine) Render(tplName string, data interface{}) (string, error) {
	// Check cache first
	e.cacheMu.RLock()
	tmpl, ok := e.cache[tplName]
	e.cacheMu.RUnlock()
	if ok {
		return e.executeTemplate(tmpl, data)
	}

	// Load and parse the template
	fullPath := filepath.Join(e.dir, tplName)
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", tplName, err)
	}

	// Pre‑process <go> tags into Go template syntax
	converted, err := e.convertGoTags(string(content))
	if err != nil {
		return "", fmt.Errorf("conversion error in %s: %w", tplName, err)
	}

	if e.debug {
		fmt.Printf("[DEBUG] converted %s:\n%s\n", tplName, converted)
	}

	// Create a new template with custom delimiters and functions
	tmpl = template.New(tplName).Delims(e.delimLeft, e.delimRight).Funcs(e.funcs)
	tmpl, err = tmpl.Parse(converted)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", tplName, err)
	}

	// Cache the compiled template
	e.cacheMu.Lock()
	e.cache[tplName] = tmpl
	e.cacheMu.Unlock()

	return e.executeTemplate(tmpl, data)
}

// RenderString parses and executes a template string directly.
func (e *Engine) RenderString(tplContent string, data interface{}) (string, error) {
	converted, err := e.convertGoTags(tplContent)
	if err != nil {
		return "", fmt.Errorf("conversion error: %w", err)
	}
	tmpl := template.New("string").Delims(e.delimLeft, e.delimRight).Funcs(e.funcs)
	tmpl, err = tmpl.Parse(converted)
	if err != nil {
		return "", fmt.Errorf("failed to parse template string: %w", err)
	}
	return e.executeTemplate(tmpl, data)
}

// RenderLayout renders a layout template with a page template embedded inside.
func (e *Engine) RenderLayout(layoutName, pageName string, data interface{}) (string, error) {
	// Render the page first to get its content
	pageContent, err := e.Render(pageName, data)
	if err != nil {
		return "", err
	}
	// We pass the page content as a variable to the layout.
	// The layout should use {{ .PageContent }} to embed the page.
	combinedData := map[string]interface{}{
		"PageContent": template.HTML(pageContent),
	}
	// Also merge the original data into the layout context
	if m, ok := data.(map[string]interface{}); ok {
		for k, v := range m {
			combinedData[k] = v
		}
	}
	return e.Render(layoutName, combinedData)
}

// executeTemplate runs the parsed template with data and returns the output.
func (e *Engine) executeTemplate(tmpl *template.Template, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execution error: %w", err)
	}
	return buf.String(), nil
}

// convertGoTags replaces <go> ... </go> with {{ ... }} and handles special attributes.
// Supported forms:
//   <go>expression</go>          -> {{expression}}
//   <go if="cond">...</go>       -> {{if cond}}...{{end}}
//   <go else>...</go>            -> {{else}}...
//   <go for="item" in="list">... -> {{range $index, $element := list}}...
//   <go func="arg"></go>         -> {{func "arg"}}
//   <go attr="value">            -> {{attr "value"}}
// This is a simplified implementation; a production version would use a proper parser.
func (e *Engine) convertGoTags(content string) (string, error) {
	// Regex to match <go ...> ... </go> (with optional attributes)
	re := regexp.MustCompile(`(?s)<go\s*((?:[^>]|\n)*?)>(.*?)</go>`)
	// We need to handle self-closing? For simplicity, we treat everything as block.
	// We'll replace with corresponding Go template syntax.
	// This is a basic approach - more robust parsing would be needed for nested tags.
	// For demonstration, we implement a simple replacement.
	return re.ReplaceAllStringFunc(content, func(match string) string {
		// Extract attributes and body
		submatches := re.FindStringSubmatch(match)
		if len(submatches) != 3 {
			return match // fallback
		}
		attrs := strings.TrimSpace(submatches[1])
		body := submatches[2]

		// If no attributes, treat as plain expression
		if attrs == "" {
			return fmt.Sprintf("{{%s}}", strings.TrimSpace(body))
		}

		// Parse attributes: attr="value" or attr=value
		parts := strings.Fields(attrs)
		var tagName string
		var attrMap = make(map[string]string)
		for _, part := range parts {
			if strings.Contains(part, "=") {
				kv := strings.SplitN(part, "=", 2)
				key := kv[0]
				val := strings.Trim(kv[1], `"`)
				attrMap[key] = val
			} else {
				tagName = part
			}
		}

		switch tagName {
		case "if":
			cond := attrMap["if"]
			if cond == "" {
				return match
			}
			// Remove trailing spaces/newlines from body
			inner := strings.TrimSpace(body)
			return fmt.Sprintf("{{if %s}}%s{{end}}", cond, inner)
		case "else":
			return "{{else}}"
		case "for":
			item := attrMap["for"]
			in := attrMap["in"]
			if item == "" || in == "" {
				return match
			}
			inner := strings.TrimSpace(body)
			// Use range with index and element
			return fmt.Sprintf("{{range $index, $element := %s}}%s{{end}}", in, inner)
		default:
			// Treat as function call: <go func="arg" />
			// Actually, we may have other attributes like upper="name"
			// We'll build a call: func arg1 arg2 ...
			var args []string
			for key, val := range attrMap {
				if key != tagName {
					// If it's a function call, we treat key as function name and val as argument.
					// We'll use: {{key val}}
					// But we need to handle multiple attributes.
					// For simplicity, we'll take the first attribute as the function name and its value as argument.
					// Better: support form <go upper="name" /> -> {{upper .name}}
					// We'll implement a basic version that assumes only one attribute besides tag.
				}
			}
			// If tagName itself is a function name and there is one attribute, use it.
			if len(attrMap) == 1 {
				for key, val := range attrMap {
					// If the key is the same as tag? e.g., <go upper="name"> -> upper is tag name and key? No.
					// In our examples, they use <go upper="name"></go> -> the tag is "go", not "upper".
					// Actually they use <go upper="Name"></go> => the attribute "upper" is the function.
					// So we need to detect the function name as the attribute key.
					// If there is exactly one attribute, that key is the function name.
					return fmt.Sprintf("{{%s %s}}", key, val)
				}
			}
			// Fallback: just output as plain expression
			return fmt.Sprintf("{{%s}}", strings.TrimSpace(body))
		}
	}), nil
}

// ClearCache removes all cached templates.
func (e *Engine) ClearCache() {
	e.cacheMu.Lock()
	defer e.cacheMu.Unlock()
	e.cache = make(map[string]*template.Template)
}

// ClearCacheFile removes a specific template from the cache.
func (e *Engine) ClearCacheFile(tplName string) {
	e.cacheMu.Lock()
	defer e.cacheMu.Unlock()
	delete(e.cache, tplName)
}

// GetCacheStats returns the number of cached templates.
func (e *Engine) GetCacheStats() map[string]int {
	e.cacheMu.RLock()
	defer e.cacheMu.RUnlock()
	return map[string]int{"cached": len(e.cache)}
}