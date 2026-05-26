# MangoTemplate - Quick Reference Cheat Sheet

## ⚡ 60-Second Syntax Reference

### Setup
```go
engine := mangotemplate.NewEngine("./templates")
engine.RegisterBuiltinFunctions()
result, _ := engine.Render("file.mango", data)
```

### Variables
```html
<go>Name</go>              <!-- Output: Alice -->
<go>User.Email</go>        <!-- Nested -->
<go upper="Name"></go>     <!-- Functions -->
```

### Conditionals
```html
<go if="Condition">true</go>
<go if="!Condition">not true</go>
<go if="Age > 18">adult</go>
<go else>child</go>
```

### Loops
```html
<go for="item" in="Items">
  <go>item</go>
</go>
```

### Functions
```html
{{ upper .name }}          {{ lower .name }}
{{ title .name }}          {{ capitalize .name }}
{{ length .items }}        {{ first .items }}
{{ add .a .b }}            {{ date .timestamp }}
```

---

## 📚 String Functions

| Function | Usage | Result |
|----------|-------|--------|
| `upper` | `<go upper="abc"></go>` | ABC |
| `lower` | `<go lower="ABC"></go>` | abc |
| `title` | `<go title="hello world"></go>` | Hello World |
| `capitalize` | `<go capitalize="hello"></go>` | Hello |
| `reverse` | `<go reverse="abc"></go>` | cba |
| `truncate` | `<go truncate="text" 10></go>` | text... |
| `slug` | `<go slug="Hello World"></go>` | hello-world |
| `strlen` | `<go strlen="hello"></go>` | 5 |
| `substr` | `<go substr="hello" 0 3></go>` | hel |
| `contains` | `<go contains="hello" "ell"></go>` | true |
| `replace` | `<go replace="hello" "l" "x"></go>` | hexxo |
| `split` | `<go split="a,b,c" ","></go>` | [a b c] |
| `join` | `<go join="[a b c]" ","></go>` | a,b,c |
| `trim` | `<go trim="  hello  "></go>` | hello |

---

## 🔢 Math Functions

| Function | Usage | Result |
|----------|-------|--------|
| `add` | `<go add="10" "5"></go>` | 15 |
| `subtract` | `<go subtract="10" "5"></go>` | 5 |
| `multiply` | `<go multiply="10" "5"></go>` | 50 |
| `divide` | `<go divide="10" "5"></go>` | 2 |
| `mod` | `<go mod="10" "3"></go>` | 1 |
| `min` | `<go min="10" "5"></go>` | 5 |
| `max` | `<go max="10" "5"></go>` | 10 |
| `abs` | `<go abs="-5"></go>` | 5 |

---

## 📅 Date/Time Functions

| Function | Usage | Result |
|----------|-------|--------|
| `date` | `<go date="timestamp"></go>` | 2024-05-27 |
| `time` | `<go time="timestamp"></go>` | 14:30:00 |
| `datetime` | `<go datetime="timestamp"></go>` | 2024-05-27 14:30:00 |
| `unix` | `<go unix="timestamp"></go>` | 1234567890 |
| `now` | `<go>now</go>` | 2024-05-27 ... |

---

## 📦 Array Functions

| Function | Usage | Result |
|----------|-------|--------|
| `first` | `<go first="items"></go>` | Item 1 |
| `last` | `<go last="items"></go>` | Item N |
| `length` | `<go length="items"></go>` | 5 |
| `range` | `<go range="1" "5"></go>` | [1 2 3 4] |

---

## ✅ Conditional Functions

| Function | Usage | Result |
|----------|-------|--------|
| `empty` | `<go if="empty items">...</go>` | true/false |
| `eq` | `<go if="eq status 'active'">...</go>` | true/false |
| `ne` | `<go if="ne status 'active'">...</go>` | true/false |
| `lt` | `<go if="lt age 18">...</go>` | true/false |
| `gt` | `<go if="gt age 18">...</go>` | true/false |
| `lte` | `<go if="lte age 18">...</go>` | true/false |
| `gte` | `<go if="gte age 18">...</go>` | true/false |
| `default` | `<go default="name" "Unknown"></go>` | Name or Unknown |

---

## 🔒 Security Functions

| Function | Usage |
|----------|-------|
| `escape` | `<go escape="html"></go>` |
| `unescape` | `<go unescape="html"></go>` |
| `safe` | `<go safe="html"></go>` |
| `md5` / `hash` | `<go md5="text"></go>` |
| `urlencode` | `<go urlencode="text"></go>` |

---

## 🎨 Type Conversion

| Function | Usage |
|----------|-------|
| `int` | `<go int="value"></go>` |
| `float` | `<go float="value"></go>` |
| `string` | `<go string="value"></go>` |
| `bool` | `<go bool="value"></go>` |
| `type` | `<go type="value"></go>` |

---

## 🚀 Common Patterns

### List with Empty State
```html
<go if="empty items">
  <p>No items</p>
</go>
<go else>
  <ul>
    <go for="item" in="items">
      <li><go>item</go></li>
    </go>
  </ul>
</go>
```

### Conditional Rendering
```html
<go if="user.IsAdmin">
  <div>Admin Panel</div>
</go>
```

### Formatted Display
```html
<p>Posted: <go date="date"></go></p>
<p>Price: $<go>price</go></p>
<p><go truncate="text" 50></go></p>
```

### Table Rows
```html
<table>
  <go for="row" in="rows">
    <tr>
      <td><go>row.Name</go></td>
      <td><go>row.Value</go></td>
    </tr>
  </go>
</table>
```

### Status Badge
```html
<go if="active">
  <span class="badge-success">Active</span>
</go>
<go else>
  <span class="badge-danger">Inactive</span>
</go>
```

### Price with Sale
```html
<go if="onSale">
  <s>$<go>originalPrice</go></s>
  <b>$<go>salePrice</go></b>
</go>
<go else>
  $<go>price</go>
</go>
```

### User Greeting
```html
<h1>
  Hello, <go title="firstName"></go>!
</h1>
```

### Pagination Links
```html
<a href="?page=<go subtract="page" 1></go>">
  Previous
</a>
<a href="?page=<go add="page" 1></go>">
  Next
</a>
```

---

## 🔧 Engine Configuration

```go
// Create engine
engine := mangotemplate.NewEngine("./templates")

// Register functions
engine.RegisterBuiltinFunctions()
engine.RegisterFunc("myFunc", myFunction)

// Configuration
engine.SetTemplateDir("./templates")
engine.SetDelimiters("{{", "}}")
engine.EnableDebug()
engine.DisableDebug()

// Rendering
result, err := engine.Render("file.mango", data)
result, err := engine.RenderString(content, data)
result, err := engine.RenderLayout("layout.mango", "page.mango", data)

// Cache
engine.ClearCache()
engine.ClearCacheFile("file.mango")
stats := engine.GetCacheStats()

// Utilities
engine.TemplateExists("file.mango")
engine.Include("other.mango", data)
```

---

## 🏗️ File Structure

```
project/
├── main.go
├── templates/
│   ├── layouts/
│   │   └── base.mango
│   ├── pages/
│   │   ├── home.mango
│   │   └── products.mango
│   └── components/
│       ├── header.mango
│       └── footer.mango
└── public/
    ├── css/
    └── js/
```

---

## 💻 Go Code Examples

### Simple Template
```go
engine := mangotemplate.NewEngine("./templates")
engine.RegisterBuiltinFunctions()

data := map[string]interface{}{
    "Name": "Alice",
    "Age":  25,
}

result, _ := engine.Render("index.mango", data)
```

### Custom Function
```go
engine.RegisterFunc("discount", func(price, percent float64) float64 {
    return price * (1 - percent/100)
})
// Use: <go>discount .price 20</go>
```

### Web Handler
```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    result, _ := engine.Render("index.mango", data)
    fmt.Fprint(w, result)
})
```

### Error Handling
```go
result, err := engine.Render("file.mango", data)
if err != nil {
    fmt.Printf("Error: %v\n", err)
}
```

---

## 🎯 Quick Debugging

```go
// Enable debug mode
engine.EnableDebug()

// Check cache
stats := engine.GetCacheStats()
fmt.Printf("Cached: %v\n", stats["cached_templates"])

// Clear cache
engine.ClearCache()

// Test template parsing
tmpl, err := engine.LoadTemplate("file.mango")
if err != nil {
    fmt.Printf("Error: %v\n", err)
}
```

---

## 📊 Comparison with Other Engines

| Feature | MangoTemplate | html/template | Pongo2 |
|---------|---------------|---------------|--------|
| Syntax familiarity | PHP-like | Custom | Django-like |
| Functions | 50+ | Few | Many |
| <go> tags | ✅ | ❌ | ❌ |
| Caching | Auto | Manual | Manual |
| Learning curve | Easy | Medium | Medium |
| Performance | Fast | Fast | Medium |

---

## ⚠️ Common Mistakes

### ❌ Wrong
```html
<go>name</go>       <!-- Variable must be capitalized -->
<go if User.Admin> <!-- Missing quotes -->
<go>for item</go>   <!-- Wrong syntax -->
```

### ✅ Correct
```html
<go>Name</go>
<go if="User.Admin">
<go for="item" in="Items">
```

---

## 🆘 Troubleshooting

| Problem | Solution |
|---------|----------|
| Template not found | Check file path and .mango extension |
| Variable not showing | Ensure key name matches in data map |
| Function not found | Register function with `RegisterFunc()` |
| Wrong output type | Use type conversion functions |
| Slow rendering | Check cache with `GetCacheStats()` |

---

## 📖 Documentation Links

- **Full Docs**: `MANGOTEMPLATE_DOCUMENTATION.md`
- **Quick Start**: `QUICKSTART_MANGO.md`
- **Examples**: `examples.go` and `example_templates.mango`
- **README**: `README_MANGO.md`

---

## 🚀 Next Steps

1. Create `templates/` directory
2. Write your first `.mango` file
3. Create Go code to render
4. Register custom functions as needed
5. Build your app!

---

## 💡 Pro Tips

1. Use meaningful variable names
2. Enable debug during development
3. Organize templates in subdirectories
4. Register all functions at startup
5. Let automatic caching work for you
6. Use layouts for consistent design
7. Test with `engine.RenderString()` first

---

**Bookmark this page! 🥭**

MangoTemplate brings the ease of PHP to Go templating.
Start building today!