package mangotemplate

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ---------------------------------------------------------------------
// AST
// ---------------------------------------------------------------------

// Node is any element of a parsed template's syntax tree. Concrete types
// are *TextNode, *VarNode, *FuncNode, *IfNode, *ForNode and *DefaultNode.
type Node interface{}

// TextNode is raw, literal output copied as-is to the rendered result.
type TextNode struct {
	Text string
}

// VarNode renders the value found at a dotted path, e.g. <go>User.Name</go>.
type VarNode struct {
	Path string
}

// FuncNode renders the result of calling a registered function, e.g.
// <go upper="UserName"></go>. Extra holds any additional attributes for
// functions that take more than one argument.
type FuncNode struct {
	FuncName string
	Path     string
	Extra    []attr
}

// IfNode renders Body when Cond is truthy, otherwise Else (which may be
// nil if no matching <go else> block follows).
type IfNode struct {
	Cond string
	Body []Node
	Else []Node
}

// ForNode renders Body once per element of the list found at ListPath,
// binding each element to VarName.
type ForNode struct {
	VarName  string
	ListPath string
	Body     []Node
}

// DefaultNode renders the value at Path, or Fallback if that value is
// missing or empty, e.g. <go default="UserName" "Guest"></go>.
type DefaultNode struct {
	Path     string
	Fallback string
}

// ---------------------------------------------------------------------
// Parser
// ---------------------------------------------------------------------

// attr is a single parsed tag attribute. A bare word or bare quoted
// string (no "key=") is represented with an empty Key.
type attr struct {
	key   string
	value string
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// parseAttrs parses the raw text between "<go" and the closing ">" of a
// tag into key="value" pairs, bare words (like "else"), and bare quoted
// strings (the fallback literal in <go default="UserName" "Guest"></go>).
func parseAttrs(s string) []attr {
	var attrs []attr
	i, n := 0, len(s)
	for i < n {
		for i < n && isSpace(s[i]) {
			i++
		}
		if i >= n {
			break
		}
		if s[i] == '"' {
			j := i + 1
			for j < n && s[j] != '"' {
				j++
			}
			attrs = append(attrs, attr{value: s[i+1 : minInt(j, n)]})
			i = j + 1
			continue
		}

		start := i
		for i < n && s[i] != '=' && !isSpace(s[i]) {
			i++
		}
		key := s[start:i]
		if key == "" {
			i++
			continue
		}

		save := i
		for i < n && isSpace(s[i]) {
			i++
		}
		if i < n && s[i] == '=' {
			i++
			for i < n && isSpace(s[i]) {
				i++
			}
			if i < n && s[i] == '"' {
				j := i + 1
				for j < n && s[j] != '"' {
					j++
				}
				attrs = append(attrs, attr{key: key, value: s[i+1 : minInt(j, n)]})
				i = j + 1
			} else {
				start = i
				for i < n && !isSpace(s[i]) {
					i++
				}
				attrs = append(attrs, attr{key: key, value: s[start:i]})
			}
		} else {
			// Bare word, e.g. "else". Don't consume the whitespace we
			// skipped looking for "=" since it separates this word from
			// whatever comes next.
			i = save
			attrs = append(attrs, attr{value: key})
		}
	}
	return attrs
}

func isElseTag(attrs []attr) bool {
	for _, a := range attrs {
		if a.key == "" && a.value == "else" {
			return true
		}
	}
	return false
}

// nodesToText concatenates the literal text content of a node list. Used
// to read the body of a plain <go>Path</go> variable tag.
func nodesToText(nodes []Node) string {
	var sb strings.Builder
	for _, n := range nodes {
		if t, ok := n.(*TextNode); ok {
			sb.WriteString(t.Text)
		}
	}
	return sb.String()
}

// buildNode turns a tag's parsed attributes plus its already-parsed inner
// content into a concrete Node.
func buildNode(attrs []attr, inner []Node) (Node, error) {
	var ifVal, forVal, inVal, defaultVal, fallback string
	var hasIf, hasFor, hasIn, hasDefault, hasFallback bool
	var funcAttrs []attr

	reserved := map[string]bool{"if": true, "for": true, "in": true, "default": true}

	for _, a := range attrs {
		switch {
		case a.key == "if":
			ifVal, hasIf = a.value, true
		case a.key == "for":
			forVal, hasFor = a.value, true
		case a.key == "in":
			inVal, hasIn = a.value, true
		case a.key == "default":
			defaultVal, hasDefault = a.value, true
		case a.key == "" && a.value == "else":
			// A stray <go else> with no preceding if; nothing to render.
			return &TextNode{Text: ""}, nil
		case a.key == "" && hasDefault && !hasFallback:
			fallback, hasFallback = a.value, true
		case a.key != "" && !reserved[a.key]:
			funcAttrs = append(funcAttrs, a)
		}
	}

	switch {
	case hasIf:
		return &IfNode{Cond: ifVal, Body: inner}, nil
	case hasFor && hasIn:
		return &ForNode{VarName: forVal, ListPath: inVal, Body: inner}, nil
	case hasDefault:
		return &DefaultNode{Path: defaultVal, Fallback: fallback}, nil
	case len(funcAttrs) > 0:
		primary := funcAttrs[0]
		return &FuncNode{FuncName: primary.key, Path: primary.value, Extra: funcAttrs[1:]}, nil
	default:
		return &VarNode{Path: strings.TrimSpace(nodesToText(inner))}, nil
	}
}

// parseTemplate parses a full .mango template into top-level nodes.
func parseTemplate(src string) ([]Node, error) {
	nodes, pos, err := parseSequence(src, 0)
	if err != nil {
		return nil, err
	}
	if pos != len(src) {
		return nil, fmt.Errorf("unexpected closing </go> tag near position %d", pos)
	}
	return nodes, nil
}

// parseSequence parses text and <go> tags starting at pos until either the
// end of src, or a closing "</go>" that belongs to an enclosing tag (in
// which case pos is left pointing at the start of that "</go>" so the
// caller can consume it).
func parseSequence(src string, pos int) ([]Node, int, error) {
	var nodes []Node
	var text strings.Builder

	flush := func() {
		if text.Len() > 0 {
			nodes = append(nodes, &TextNode{Text: text.String()})
			text.Reset()
		}
	}

	for pos < len(src) {
		rest := src[pos:]
		openRel := strings.Index(rest, "<go")
		closeRel := strings.Index(rest, "</go>")

		if openRel == -1 && closeRel == -1 {
			text.WriteString(rest)
			pos = len(src)
			break
		}

		if closeRel != -1 && (openRel == -1 || closeRel < openRel) {
			text.WriteString(rest[:closeRel])
			pos += closeRel
			flush()
			return nodes, pos, nil
		}

		// An opening "<go" tag comes next.
		text.WriteString(rest[:openRel])
		pos += openRel
		flush()

		gtRel := strings.IndexByte(src[pos:], '>')
		if gtRel == -1 {
			return nil, pos, fmt.Errorf("unterminated <go tag at position %d", pos)
		}
		header := src[pos+3 : pos+gtRel]
		pos += gtRel + 1

		innerNodes, newPos, err := parseSequence(src, pos)
		if err != nil {
			return nil, pos, err
		}
		pos = newPos
		if !strings.HasPrefix(src[pos:], "</go>") {
			return nil, pos, fmt.Errorf("missing closing </go> tag for <go%s>", header)
		}
		pos += len("</go>")

		node, err := buildNode(parseAttrs(header), innerNodes)
		if err != nil {
			return nil, pos, err
		}

		if ifNode, ok := node.(*IfNode); ok {
			if elseBody, newPos2, matched, err := tryParseElse(src, pos); err != nil {
				return nil, pos, err
			} else if matched {
				ifNode.Else = elseBody
				pos = newPos2
			}
		}

		nodes = append(nodes, node)
	}

	flush()
	return nodes, pos, nil
}

// tryParseElse looks ahead from pos, skipping insignificant whitespace,
// for a "<go else>...</go>" block immediately following an if block. If
// found, it returns the parsed body and the position right after it.
func tryParseElse(src string, pos int) (body []Node, newPos int, matched bool, err error) {
	p := pos
	for p < len(src) && isSpace(src[p]) {
		p++
	}
	if !strings.HasPrefix(src[p:], "<go") {
		return nil, pos, false, nil
	}
	gtRel := strings.IndexByte(src[p:], '>')
	if gtRel == -1 {
		return nil, pos, false, nil
	}
	header := src[p+3 : p+gtRel]
	if !isElseTag(parseAttrs(header)) {
		return nil, pos, false, nil
	}

	bodyStart := p + gtRel + 1
	innerNodes, afterBody, err := parseSequence(src, bodyStart)
	if err != nil {
		return nil, pos, false, err
	}
	if !strings.HasPrefix(src[afterBody:], "</go>") {
		return nil, pos, false, fmt.Errorf("missing closing </go> tag for <go else> at position %d", p)
	}
	afterBody += len("</go>")
	return innerNodes, afterBody, true, nil
}

// ---------------------------------------------------------------------
// Data resolution
// ---------------------------------------------------------------------

// scope is a chain of loop-variable bindings layered over the root data
// passed to Render. Lookups check the nearest enclosing loop variables
// first, then fall back to the root data.
type scope struct {
	vars   map[string]interface{}
	parent *scope
	root   interface{}
}

func newRootScope(data interface{}) *scope {
	return &scope{vars: map[string]interface{}{}, root: data}
}

func (s *scope) child() *scope {
	return &scope{vars: map[string]interface{}{}, parent: s, root: s.root}
}

// resolve looks up a dotted path such as "User.Address.City".
func (s *scope) resolve(path string) (interface{}, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}
	parts := strings.Split(path, ".")
	first := parts[0]

	var base interface{}
	found := false
	for cur := s; cur != nil; cur = cur.parent {
		if v, ok := cur.vars[first]; ok {
			base, found = v, true
			break
		}
	}
	if !found {
		v, err := fieldAccess(s.root, first)
		if err != nil {
			return nil, fmt.Errorf("%q: %w", path, err)
		}
		base = v
	}

	for _, p := range parts[1:] {
		v, err := fieldAccess(base, p)
		if err != nil {
			return nil, fmt.Errorf("%q: %w", path, err)
		}
		base = v
	}
	return base, nil
}

// fieldAccess resolves a single path segment against a map key or struct
// field, transparently dereferencing pointers and interfaces.
func fieldAccess(obj interface{}, name string) (interface{}, error) {
	if obj == nil {
		return nil, fmt.Errorf("cannot access %q on nil value", name)
	}
	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil, fmt.Errorf("cannot access %q on nil value", name)
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Map:
		mv := v.MapIndex(reflect.ValueOf(name))
		if !mv.IsValid() {
			return nil, fmt.Errorf("key %q not found", name)
		}
		return mv.Interface(), nil
	case reflect.Struct:
		fv := v.FieldByName(name)
		if !fv.IsValid() {
			return nil, fmt.Errorf("field %q not found", name)
		}
		if !fv.CanInterface() {
			return nil, fmt.Errorf("field %q is unexported", name)
		}
		return fv.Interface(), nil
	default:
		return nil, fmt.Errorf("cannot access field %q on %s", name, v.Kind())
	}
}

func isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return t != ""
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return false
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String, reflect.Chan:
		return rv.Len() > 0
	case reflect.Bool:
		return rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() != 0
	default:
		return true
	}
}

func isEmptyValue(v interface{}) bool {
	if v == nil {
		return true
	}
	if s, ok := v.(string); ok {
		return s == ""
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return true
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String, reflect.Chan:
		return rv.Len() == 0
	default:
		return false
	}
}

func formatValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case fmt.Stringer:
		return t.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ---------------------------------------------------------------------
// Built-in & custom functions
// ---------------------------------------------------------------------

var builtinFuncs = map[string]Func{
	"upper":  funcUpper,
	"lower":  funcLower,
	"length": funcLength,
	"date":   funcDate,
}

func funcUpper(v interface{}) (string, error) {
	return strings.ToUpper(formatValue(v)), nil
}

func funcLower(v interface{}) (string, error) {
	return strings.ToLower(formatValue(v)), nil
}

func funcLength(v interface{}) (int, error) {
	if v == nil {
		return 0, nil
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return 0, nil
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String, reflect.Chan:
		return rv.Len(), nil
	default:
		return 0, fmt.Errorf("length: unsupported type %s", rv.Kind())
	}
}

// dateLayout is the format used by the built-in "date" function: 2006-01-02.
const dateLayout = "2006-01-02"

func funcDate(v interface{}) (string, error) {
	switch t := v.(type) {
	case time.Time:
		return t.Format(dateLayout), nil
	case *time.Time:
		if t == nil {
			return "", nil
		}
		return t.Format(dateLayout), nil
	default:
		return "", fmt.Errorf("date: expected time.Time, got %T", v)
	}
}

var errType = reflect.TypeOf((*error)(nil)).Elem()

// callFunc invokes fn (a builtin or a value registered via
// Engine.RegisterFunc) with args, converting each argument to the
// function's declared parameter type where possible.
func callFunc(fn Func, args []interface{}) (interface{}, error) {
	fv := reflect.ValueOf(fn)
	if fv.Kind() != reflect.Func {
		return nil, fmt.Errorf("registered value is not a function")
	}
	ft := fv.Type()
	numIn := ft.NumIn()

	in := make([]reflect.Value, 0, numIn)
	for i := 0; i < numIn; i++ {
		paramType := ft.In(i)
		if i < len(args) {
			argVal, err := convertArg(args[i], paramType)
			if err != nil {
				return nil, err
			}
			in = append(in, argVal)
		} else {
			in = append(in, reflect.Zero(paramType))
		}
	}

	out := fv.Call(in)

	var result interface{}
	var callErr error
	for _, o := range out {
		if o.Type().Implements(errType) {
			if !o.IsNil() {
				callErr, _ = o.Interface().(error)
			}
		} else {
			result = o.Interface()
		}
	}
	return result, callErr
}

// convertArg converts v to target's type, supporting direct assignment,
// reflect-convertible numeric/string/bool kinds, parsing numeric strings,
// and passing through to interface{} parameters.
func convertArg(v interface{}, target reflect.Type) (reflect.Value, error) {
	if v == nil {
		return reflect.Zero(target), nil
	}
	if target.Kind() == reflect.Interface {
		return reflect.ValueOf(v), nil
	}

	rv := reflect.ValueOf(v)
	if rv.Type().AssignableTo(target) {
		return rv, nil
	}
	if rv.Type().ConvertibleTo(target) {
		switch target.Kind() {
		case reflect.String, reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			return rv.Convert(target), nil
		}
	}
	if s, ok := v.(string); ok {
		switch target.Kind() {
		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("cannot convert %q to %s", s, target)
			}
			return reflect.ValueOf(f).Convert(target), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("cannot convert %q to %s", s, target)
			}
			return reflect.ValueOf(n).Convert(target), nil
		case reflect.Bool:
			b, err := strconv.ParseBool(s)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("cannot convert %q to %s", s, target)
			}
			return reflect.ValueOf(b), nil
		}
	}
	return reflect.Value{}, fmt.Errorf("cannot convert %T to %s", v, target)
}

func lookupFunc(e *Engine, name string) (Func, bool) {
	e.funcsMu.RLock()
	fn, ok := e.funcs[name]
	e.funcsMu.RUnlock()
	if ok {
		return fn, true
	}
	fn, ok = builtinFuncs[name]
	return fn, ok
}

// ---------------------------------------------------------------------
// Rendering
// ---------------------------------------------------------------------

func renderNodes(sb *strings.Builder, nodes []Node, ctx *scope, e *Engine) error {
	for _, node := range nodes {
		switch n := node.(type) {
		case *TextNode:
			sb.WriteString(n.Text)

		case *VarNode:
			v, err := ctx.resolve(n.Path)
			if err != nil {
				e.debugf("variable %q: %v", n.Path, err)
				continue
			}
			sb.WriteString(formatValue(v))

		case *FuncNode:
			result, err := renderFuncNode(n, ctx, e)
			if err != nil {
				e.debugf("function %q: %v", n.FuncName, err)
				continue
			}
			sb.WriteString(formatValue(result))

		case *IfNode:
			condVal, _ := ctx.resolve(n.Cond)
			if isTruthy(condVal) {
				if err := renderNodes(sb, n.Body, ctx, e); err != nil {
					return err
				}
			} else if n.Else != nil {
				if err := renderNodes(sb, n.Else, ctx, e); err != nil {
					return err
				}
			}

		case *ForNode:
			if err := renderForNode(sb, n, ctx, e); err != nil {
				return err
			}

		case *DefaultNode:
			v, err := ctx.resolve(n.Path)
			if err != nil || isEmptyValue(v) {
				sb.WriteString(n.Fallback)
			} else {
				sb.WriteString(formatValue(v))
			}

		default:
			return fmt.Errorf("unknown node type %T", node)
		}
	}
	return nil
}

func renderFuncNode(n *FuncNode, ctx *scope, e *Engine) (interface{}, error) {
	fn, ok := lookupFunc(e, n.FuncName)
	if !ok {
		return nil, fmt.Errorf("unknown function %q", n.FuncName)
	}

	args := make([]interface{}, 0, len(n.Extra)+1)
	primary, err := ctx.resolve(n.Path)
	if err != nil {
		primary = n.Path
	}
	args = append(args, primary)

	for _, extra := range n.Extra {
		v, err := ctx.resolve(extra.value)
		if err != nil {
			v = extra.value
		}
		args = append(args, v)
	}

	return callFunc(fn, args)
}

func renderForNode(sb *strings.Builder, n *ForNode, ctx *scope, e *Engine) error {
	listVal, err := ctx.resolve(n.ListPath)
	if err != nil {
		e.debugf("for loop %q: %v", n.ListPath, err)
		return nil
	}
	if listVal == nil {
		return nil
	}

	rv := reflect.ValueOf(listVal)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			child := ctx.child()
			child.vars[n.VarName] = rv.Index(i).Interface()
			if err := renderNodes(sb, n.Body, child, e); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, k := range rv.MapKeys() {
			child := ctx.child()
			child.vars[n.VarName] = rv.MapIndex(k).Interface()
			if err := renderNodes(sb, n.Body, child, e); err != nil {
				return err
			}
		}
	default:
		e.debugf("for loop %q: not iterable (%s)", n.ListPath, rv.Kind())
	}
	return nil
}