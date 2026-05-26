package mangotemplate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// Template represents a Mango template
type Template struct {
	Name     string
	Content  string
	Compiled *template.Template
	Engine   *Engine
}

// Engine is the MangoTemplate rendering engine
type Engine struct {
	TemplateDir string
	CachePath   string
	Templates   map[string]*Template
	Delims      Delimiters
	Debug       bool
	FuncMap     template.FuncMap
}

// Delimiters defines custom template delimiters
type Delimiters struct {
	Open  string
	Close string
}

// RenderContext holds the data passed to templates
type RenderContext struct {
	Data map[string]interface{}
	Vars map[string]interface{}
}

// NewEngine creates a new MangoTemplate engine
func NewEngine(templateDir string) *Engine {
	if templateDir == "" {
		templateDir = "./templates"
	}

	engine := &Engine{
		TemplateDir: templateDir,
		CachePath:   filepath.Join(templateDir, ".mango_cache"),
		Templates:   make(map[string]*Template),
		Delims: Delimiters{
			Open:  "{{",
			Close: "}}",
		},
		Debug:   false,
		FuncMap: make(template.FuncMap),
	}

	// Ensure template directory exists
	os.MkdirAll(templateDir, 0755)
	os.MkdirAll(engine.CachePath, 0755)

	return engine
}

// SetDelimiters sets custom template delimiters
func (e *Engine) SetDelimiters(open, close string) {
	e.Delims.Open = open
	e.Delims.Close = close
}

// RegisterFunc registers a custom template function
func (e *Engine) RegisterFunc(name string, fn interface{}) {
	e.FuncMap[name] = fn
}

// LoadTemplate loads a .mango or .html template
func (e *Engine) LoadTemplate(filename string) (*Template, error) {
	// Support both .mango and .html extensions
	validExtensions := []string{".mango", ".html"}
	hasValidExt := false
	for _, ext := range validExtensions {
		if strings.HasSuffix(filename, ext) {
			hasValidExt = true
			break
		}
	}

	if !hasValidExt {
		filename = strings.TrimSuffix(filename, filepath.Ext(filename)) + ".mango"
	}

	// Check cache first
	if tmpl, exists := e.Templates[filename]; exists {
		return tmpl, nil
	}

	// Load from file
	filepath := filepath.Join(e.TemplateDir, filename)
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template: %w", err)
	}

	tmpl := &Template{
		Name:    filename,
		Content: string(content),
		Engine:  e,
	}

	// Parse and compile the template
	if err := tmpl.Parse(); err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	// Cache the template
	e.Templates[filename] = tmpl

	if e.Debug {
		fmt.Printf("[MangoTemplate] Loaded: %s\n", filename)
	}

	return tmpl, nil
}

// Parse parses a Mango template and converts it to Go template
func (t *Template) Parse() error {
	// Process the content
	processed := t.processMangoSyntax(t.Content)

	// Create compiled template
	tmpl, err := template.New(t.Name).
		Delims(t.Engine.Delims.Open, t.Engine.Delims.Close).
		Funcs(t.Engine.FuncMap).
		Parse(processed)

	if err != nil {
		return fmt.Errorf("template parsing error: %w", err)
	}

	t.Compiled = tmpl
	return nil
}

// processMangoSyntax processes Mango-specific syntax
func (t *Template) processMangoSyntax(content string) string {
	// Extract and process <go> tags
	goTagRegex := regexp.MustCompile(`(?s)<go>(.*?)</go>`)
	processed := content

	// Replace <go> tags with Go template syntax
	processed = goTagRegex.ReplaceAllStringFunc(processed, func(match string) string {
		// Extract content between tags
		content := goTagRegex.FindStringSubmatch(match)[1]
		content = strings.TrimSpace(content)

		// Determine if it's an expression, statement, or block
		if strings.HasPrefix(content, "if ") || strings.HasPrefix(content, "for ") || strings.HasPrefix(content, "range ") {
			return fmt.Sprintf("{{ %s }}", content)
		}

		if strings.HasPrefix(content, "end") {
			return "{{ end }}"
		}

		if strings.HasPrefix(content, "else") {
			return "{{ else }}"
		}

		if strings.HasPrefix(content, ":=") || strings.Contains(content, "=") {
			return fmt.Sprintf("{{ %s }}", content)
		}

		// Default: treat as expression output
		return fmt.Sprintf("{{ %s }}", content)
	})

	// Handle short echo syntax: <go echo="variable">
	echoRegex := regexp.MustCompile(`<go\s+echo="([^"]+)"\s*/>`)
	processed = echoRegex.ReplaceAllString(processed, "{{ .$1 }}")

	// Handle variable output: <go print="variable">
	printRegex := regexp.MustCompile(`<go\s+print="([^"]+)"\s*/>`)
	processed = printRegex.ReplaceAllString(processed, "{{ .$1 }}")

	// Handle if statements with attributes
	ifRegex := regexp.MustCompile(`<go\s+if="([^"]+)">`)
	processed = ifRegex.ReplaceAllString(processed, "{{ if .$1 }}")

	// Handle for loops with attributes
	forRegex := regexp.MustCompile(`<go\s+for="([^"]+)"\s+in="([^"]+)">`)
	processed = forRegex.ReplaceAllString(processed, "{{ range .$2 }}")

	// Handle variable assignment
	assignRegex := regexp.MustCompile(`<go\s+var="([^"]+)"\s+value="([^"]+)"\s*/>`)
	processed = assignRegex.ReplaceAllString(processed, "{{ $.$1 := .$2 }}")

	return processed
}

// Render renders the template with the given data
func (t *Template) Render(data interface{}) (string, error) {
	if t.Compiled == nil {
		return "", fmt.Errorf("template not compiled")
	}

	var buf bytes.Buffer
	if err := t.Compiled.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execution error: %w", err)
	}

	return buf.String(), nil
}

// RenderToWriter renders the template to a writer
func (t *Template) RenderToWriter(w bytes.Buffer, data interface{}) error {
	if t.Compiled == nil {
		return fmt.Errorf("template not compiled")
	}

	return t.Compiled.Execute(&w, data)
}

// Render is a convenience method on Engine to load and render a template
func (e *Engine) Render(filename string, data interface{}) (string, error) {
	tmpl, err := e.LoadTemplate(filename)
	if err != nil {
		return "", err
	}

	return tmpl.Render(data)
}

// RenderString renders a template from a string
func (e *Engine) RenderString(content string, data interface{}) (string, error) {
	tmpl := &Template{
		Name:    "inline",
		Content: content,
		Engine:  e,
	}

	if err := tmpl.Parse(); err != nil {
		return "", err
	}

	return tmpl.Render(data)
}

// RenderLayout renders a template with a layout
func (e *Engine) RenderLayout(layoutFile, contentFile string, data interface{}) (string, error) {
	content, err := e.Render(contentFile, data)
	if err != nil {
		return "", err
	}

	// Add content to data for layout
	contextData := map[string]interface{}{
		"content": template.HTML(content),
	}

	// Merge with existing data if it's a map
	if dataMap, ok := data.(map[string]interface{}); ok {
		for k, v := range dataMap {
			contextData[k] = v
		}
	}

	return e.Render(layoutFile, contextData)
}

// Include includes another template within current template
func (e *Engine) Include(filename string, data interface{}) (string, error) {
	return e.Render(filename, data)
}

// ClearCache clears the template cache
func (e *Engine) ClearCache() {
	e.Templates = make(map[string]*Template)
	if e.Debug {
		fmt.Println("[MangoTemplate] Cache cleared")
	}
}

// ClearCacheFile clears a specific template from cache
func (e *Engine) ClearCacheFile(filename string) {
	delete(e.Templates, filename)
	if e.Debug {
		fmt.Printf("[MangoTemplate] Cache cleared for: %s\n", filename)
	}
}

// GetCacheStats returns cache statistics
func (e *Engine) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"cached_templates": len(e.Templates),
		"template_names":   getTemplateNames(e.Templates),
	}
}

// EnableDebug enables debug mode
func (e *Engine) EnableDebug() {
	e.Debug = true
}

// DisableDebug disables debug mode
func (e *Engine) DisableDebug() {
	e.Debug = false
}

// Helper function to get template names
func getTemplateNames(templates map[string]*Template) []string {
	names := make([]string, 0, len(templates))
	for name := range templates {
		names = append(names, name)
	}
	return names
}

// SetTemplateDir sets the template directory
func (e *Engine) SetTemplateDir(dir string) error {
	e.TemplateDir = dir
	e.CachePath = filepath.Join(dir, ".mango_cache")

	// Create directories
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(e.CachePath, 0755); err != nil {
		return err
	}

	e.ClearCache()
	return nil
}

// GetTemplate returns a loaded template
func (e *Engine) GetTemplate(filename string) (*Template, error) {
	return e.LoadTemplate(filename)
}

// TemplateExists checks if a template file exists
func (e *Engine) TemplateExists(filename string) bool {
	path := filepath.Join(e.TemplateDir, filename)
	_, err := os.Stat(path)
	return err == nil
}