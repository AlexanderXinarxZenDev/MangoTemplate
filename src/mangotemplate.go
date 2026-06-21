// Package mangotemplate is a tiny HTML-friendly template engine built around
// a single custom element, <go>, which is used for variable interpolation,
// built-in/custom function calls, conditionals, loops and default values.
//
// See render.go for the parser and renderer implementation.
package mangotemplate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Func is the signature accepted by RegisterFunc. Custom functions are
// invoked via reflection, so any function with a fixed, non-variadic
// signature is allowed (see callFunc in render.go for the exact rules).
type Func interface{}

// compiledTemplate is a parsed template ready to be rendered repeatedly
// without re-parsing.
type compiledTemplate struct {
	name string
	root []Node
}

// Engine renders .mango templates loaded from a directory on disk. Create
// one Engine per application and reuse it; parsed templates are cached
// automatically.
type Engine struct {
	templatesDir string

	mu    sync.RWMutex
	cache map[string]*compiledTemplate

	funcsMu sync.RWMutex
	funcs   map[string]Func

	debug bool

	statsMu      sync.Mutex
	renderCount  int
	cacheHits    int
	cacheMisses  int
}

// NewEngine creates a new template engine that loads templates relative to
// templatesDir. Templates are cached the first time they are rendered.
func NewEngine(templatesDir string) *Engine {
	return &Engine{
		templatesDir: templatesDir,
		cache:        make(map[string]*compiledTemplate),
		funcs:        make(map[string]Func),
	}
}

// RegisterBuiltinFunctions registers the engine's built-in template
// functions (upper, lower, length, date) under the names used by <go>
// attributes such as <go upper="name"></go>.
//
// Calling this is optional: the built-ins are always available even
// without calling it, but it mirrors the documented quick-start API and
// gives a place to hook future opt-in behavior.
func (e *Engine) RegisterBuiltinFunctions() {
	e.funcsMu.Lock()
	defer e.funcsMu.Unlock()
	for name, fn := range builtinFuncs {
		e.funcs[name] = fn
	}
}

// RegisterFunc registers a custom function that can be used as a <go>
// attribute, e.g. engine.RegisterFunc("discount", fn) enables
// <go discount="price"></go> style usage for single-argument functions,
// or via custom AST extensions for multi-argument ones.
func (e *Engine) RegisterFunc(name string, fn Func) {
	e.funcsMu.Lock()
	defer e.funcsMu.Unlock()
	e.funcs[name] = fn
}

// EnableDebug turns on verbose logging of template loads, cache
// hits/misses, and render errors to stderr.
func (e *Engine) EnableDebug() {
	e.debug = true
}

func (e *Engine) debugf(format string, args ...interface{}) {
	if e.debug {
		fmt.Fprintf(os.Stderr, "[mangotemplate] "+format+"\n", args...)
	}
}

// GetCacheStats returns simple counters about template caching and
// rendering, useful for monitoring in production as described in the
// quick start guide.
func (e *Engine) GetCacheStats() map[string]interface{} {
	e.mu.RLock()
	cached := len(e.cache)
	e.mu.RUnlock()

	e.statsMu.Lock()
	defer e.statsMu.Unlock()

	return map[string]interface{}{
		"cached_templates": cached,
		"render_count":     e.renderCount,
		"cache_hits":       e.cacheHits,
		"cache_misses":     e.cacheMisses,
	}
}

// Render loads (or reuses a cached copy of) the named template and
// executes it against data, returning the resulting string.
//
// name is resolved relative to the engine's templates directory, e.g.
// Render("index.mango", data) loads "<templatesDir>/index.mango".
func (e *Engine) Render(name string, data interface{}) (string, error) {
	e.statsMu.Lock()
	e.renderCount++
	e.statsMu.Unlock()

	tmpl, err := e.load(name)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	ctx := newRootScope(data)
	if err := renderNodes(&sb, tmpl.root, ctx, e); err != nil {
		e.debugf("render error in %s: %v", name, err)
		return "", fmt.Errorf("mangotemplate: render %s: %w", name, err)
	}

	return sb.String(), nil
}

// load returns the compiled template for name, parsing and caching it on
// first use.
func (e *Engine) load(name string) (*compiledTemplate, error) {
	e.mu.RLock()
	tmpl, ok := e.cache[name]
	e.mu.RUnlock()
	if ok {
		e.statsMu.Lock()
		e.cacheHits++
		e.statsMu.Unlock()
		e.debugf("cache hit: %s", name)
		return tmpl, nil
	}

	e.statsMu.Lock()
	e.cacheMisses++
	e.statsMu.Unlock()

	path := filepath.Join(e.templatesDir, name)
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("mangotemplate: failed to load template: %w", err)
	}

	e.debugf("parsing template: %s", name)
	nodes, err := parseTemplate(string(raw))
	if err != nil {
		return nil, fmt.Errorf("mangotemplate: template parsing error in %s: %w", name, err)
	}

	tmpl = &compiledTemplate{name: name, root: nodes}

	e.mu.Lock()
	e.cache[name] = tmpl
	e.mu.Unlock()

	return tmpl, nil
}

// InvalidateCache clears all cached compiled templates, forcing the next
// Render call for each template to re-read and re-parse it from disk.
func (e *Engine) InvalidateCache() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cache = make(map[string]*compiledTemplate)
}