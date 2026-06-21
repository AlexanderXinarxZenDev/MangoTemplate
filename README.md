# MangoTemplate 🥭 [1.9]

A powerful, PHP-like templating engine for Go that brings familiar web templating paradigms to the Go ecosystem.

[![Go Version](https://img.shields.io/badge/Go-1.16%2B-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Active-brightgreen)]()

---

## 🎯 What is MangoTemplate?

MangoTemplate is a templating library for Go that supports `.mango` and `.html` files with **PHP-like syntax**. It allows you to embed Go code directly in your templates using `<go>` tags, making it feel natural for developers coming from dynamic web languages.

```html
<!-- MangoTemplate .mango file -->
<!DOCTYPE html>
<html>
<body>
    <h1>Hello, <go>Name</go>!</h1>
    
    <go if="User.IsAdmin">
        <p>Welcome Administrator</p>
    </go>
    
    <ul>
        <go for="item" in="Items">
            <li><go>item.Name</go> - $<go>item.Price</go></li>
        </go>
    </ul>
</body>
</html>
```

---

## ✨ Key Features

- ✅ **PHP-like Syntax**: Familiar to web developers
- ✅ **`<go>` Tags**: Embed Go code directly in templates
- ✅ **50+ Built-in Functions**: String, math, date, array operations
- ✅ **Smart Caching**: Automatic template caching for performance
- ✅ **Layouts & Includes**: Template composition and inheritance
- ✅ **Custom Functions**: Register your own template functions
- ✅ **Error Friendly**: Clear error messages for debugging
- ✅ **Debug Mode**: Enhanced logging during development
- ✅ **Type Safe**: Full Go type support
- ✅ **Web Ready**: Perfect for HTTP servers

---

## 🚀 Quick Start

### 1. Install

```bash
go get github.com/AlexanderXinaxrZenDev/mango_template/src
```

### 2. Create a Template

**templates/index.mango**
```html
<!DOCTYPE html>
<html>
<head>
    <title><go>Title</go></title>
</head>
<body>
    <h1>Welcome, <go upper="Name"></go>!</h1>
</body>
</html>
```

### 3. Use in Go

```go
package main

import (
	"fmt"
	"github.com/AlexanderXinaxrZenDev/mangotemplate"
)

func main() {
	engine := mangotemplate.NewEngine("./templates")
	engine.RegisterBuiltinFunctions()

	data := map[string]interface{}{
		"Title": "Home",
		"Name":  "alice",
	}

	result, _ := engine.Render("index.mango", data)
	fmt.Println(result)
}
```

### Output

```html
<!DOCTYPE html>
<html>
<head>
    <title>Home</title>
</head>
<body>
    <h1>Welcome, ALICE!</h1>
</body>
</html>
```

---

## 📚 Syntax Overview

### Variables

```html
<go>VariableName</go>              <!-- Output variable -->
<go>User.Name</go>                  <!-- Nested access -->
<go>User.Address.City</go>          <!-- Deep nesting -->
```

### Functions

```html
<go upper="name"></go>              <!-- UPPERCASE -->
<go lower="name"></go>              <!-- lowercase -->
<go title="name"></go>              <!-- Title Case -->
<go truncate="text" 50></go>        <!-- Truncate -->
<go date="timestamp"></go>          <!-- Format date -->
<go length="items"></go>            <!-- Array length -->
```

### Conditionals

```html
<go if="Condition">
    Content if true
</go>

<go if="User.IsAdmin">
    Admin area
</go>
<go else>
    User area
</go>
```

### Loops

```html
<go for="item" in="Items">
    <div>
        <h3><go>item.Name</go></h3>
        <p>$<go>item.Price</go></p>
    </div>
</go>

<!-- Nested loops -->
<go for="category" in="Categories">
    <h2><go>category.Name</go></h2>
    <go for="product" in="category.Products">
        <span><go>product.Name</go></span>
    </go>
</go>
```

---

## 🛠️ Built-in Functions (50+)

### String Functions
- `upper`, `lower`, `title`, `capitalize`
- `reverse`, `truncate`, `slug`
- `contains`, `replace`, `split`, `join`
- `trim`, `strlen`, `substr`

### Math Functions
- `add`, `subtract`, `multiply`, `divide`, `mod`
- `min`, `max`, `abs`

### Date/Time Functions
- `date`, `time`, `datetime`, `unix`, `now`

### Array Functions
- `first`, `last`, `length`, `range`

### Conditional Functions
- `default`, `empty`, `eq`, `ne`, `lt`, `gt`, `lte`, `gte`

### Type Conversion
- `int`, `float`, `string`, `bool`, `type`

### Security
- `escape`, `unescape`, `safe`, `md5`, `hash`, `urlencode`

---

## 📖 Examples

### Blog Post Template

```html
<!DOCTYPE html>
<html>
<head>
    <title><go>post.Title</go></title>
</head>
<body>
    <article>
        <h1><go>post.Title</go></h1>
        <p class="meta">
            By <go>post.Author</go> on <go date="post.PublishedDate"></go>
        </p>
        <div class="content">
            <go>post.Content</go>
        </div>
        <footer>
            Tags:
            <go for="tag" in="post.Tags">
                <a href="/tags/<go slug="tag"></go>"><go>tag</go></a>
            </go>
        </footer>
    </article>
</body>
</html>
```

### E-commerce Product Page

```html
<div class="product">
    <h2><go>product.Name</go></h2>
    <p><go>product.Description</go></p>
    
    <div class="price">
        <go if="product.OnSale">
            <span class="sale">$<go>product.SalePrice</go></span>
            <span class="original">$<go>product.Price</go></span>
        </go>
        <go else>
            $<go>product.Price</go>
        </go>
    </div>

    <go if="product.InStock">
        <button>Add to Cart</button>
    </go>
    <go else>
        <p>Out of Stock</p>
    </go>
</div>
```

### Dashboard with Conditionals

```html
<h1>Welcome, <go>user.Name</go>!</h1>

<go if="user.Premium">
    <span class="badge">Premium Member</span>
</go>

<div class="orders">
    <h2>Your Orders</h2>
    <go if="empty user.Orders">
        <p>You haven't placed any orders yet.</p>
    </go>
    <go else>
        <table>
            <tr>
                <th>ID</th>
                <th>Date</th>
                <th>Total</th>
            </tr>
            <go for="order" in="user.Orders">
                <tr>
                    <td><go>order.ID</go></td>
                    <td><go date="order.Date"></go></td>
                    <td>$<go>order.Total</go></td>
                </tr>
            </go>
        </table>
    </go>
</div>
```

---

## 🔧 API Reference

### Creating Engine

```go
engine := mangotemplate.NewEngine("./templates")
engine.RegisterBuiltinFunctions()
```

### Rendering

```go
// Render template from file
result, err := engine.Render("index.mango", data)

// Render from string
result, err := engine.RenderString(templateContent, data)

// Render with layout
result, err := engine.RenderLayout("layout.mango", "page.mango", data)
```

### Custom Functions

```go
engine.RegisterFunc("discount", func(price, percent float64) float64 {
    return price * (1 - percent/100)
})

// Usage: {{ discount .price 20 }}
```

### Configuration

```go
engine.SetTemplateDir("./new_dir")
engine.SetDelimiters("{{", "}}")
engine.EnableDebug()
engine.DisableDebug()
```

### Cache Management

```go
engine.ClearCache()                    // Clear all
engine.ClearCacheFile("index.mango")   // Clear one
stats := engine.GetCacheStats()        // Get statistics
```

---

## 💡 Use Cases

- 🌐 **Web Applications**: Full-stack Go web apps
- 📧 **Email Templates**: HTML emails with dynamic content
- 📄 **Document Generation**: HTML/PDF reports
- 🎨 **Static Site Generators**: Dynamic HTML generation
- 📱 **API Response Templates**: Format API responses
- 🔔 **Notification Templates**: SMS/Push notifications
- 📊 **Dashboard Templates**: Admin panels
- 🎯 **Landing Pages**: Marketing pages

---

## 📊 Performance

- **Automatic Caching**: Templates cached after first render
- **Zero Dependencies**: Pure Go implementation
- **Fast Parsing**: Efficient regex-based parser
- **Minimal Overhead**: Leverages Go's `text/template`

---

## 🛡️ Security

- **HTML Escaping**: Automatic escaping of output
- `escape` / `unescape` functions for control
- `safe` function for trusted content
- No code injection vulnerabilities
- Type-safe Go integration

---

## 📦 File Structure

```
mangotemplate/
├── mangotemplate.go              # Core engine
├── mangotemplate_functions.go    # Built-in functions
├── MANGOTEMPLATE_DOCUMENTATION.md # Full documentation
├── QUICKSTART_MANGO.md           # Quick start guide
├── examples.go                   # Practical examples
├── example_templates.mango       # Template examples
└── README.md                     # This file
```

---

## 🚦 Getting Started

1. **[Quick Start Guide](QUICKSTART_MANGO.md)** - 5-minute setup
2. **[Full Documentation](MANGOTEMPLATE_DOCUMENTATION.md)** - Complete reference
3. **[Examples](examples.go)** - 10 practical examples
4. **[Template Examples](example_templates.mango)** - Template samples

---

## 🤝 Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

---

## 📝 License

MIT License - feel free to use in any project

---

## 🎉 What's Next?

- [ ] Support for template inheritance
- [ ] Caching strategies (memory, file, Redis)
- [ ] Whitespace control
- [ ] Macro/snippet support
- [ ] Template includes with parameters
- [ ] Better error messages with line numbers
- [ ] Performance optimizations
- [ ] CLI tool for template testing

---

## ❓ FAQ

### Q: Can I use Go code in templates?
**A:** Yes! Use `<go>` tags for Go code and template syntax.

### Q: Is it production-ready?
**A:** Yes! MangoTemplate is designed for production use.

### Q: How does it compare to Go's `html/template`?
**A:** MangoTemplate provides a PHP-like interface over `text/template`, making it more familiar to web developers.

### Q: Can I register custom functions?
**A:** Absolutely! Use `RegisterFunc()` to add your own functions.

### Q: What about security?
**A:** Built-in HTML escaping and type safety. Use `safe` for trusted content.

### Q: How's performance?
**A:** Excellent! Automatic caching and minimal overhead.

---

## 📞 Support

For issues, questions, or suggestions:
- Open an issue on GitHub
- Check documentation first
- See examples for common patterns

---

## 🌟 Why MangoTemplate?

```
┌─────────────────────────────────────────────┐
│  Feature          │  MangoTemplate  │  Text/Template
├──────────────────┼─────────────────┼─────────────┤
│  Familiar syntax  │  ✅ PHP-like    │  ✗ Complex
│  <go> tags       │  ✅ Yes         │  ✗ No
│  Built-in funcs  │  ✅ 50+         │  ✓ Basic
│  Caching         │  ✅ Auto        │  ✗ Manual
│  Easy to learn   │  ✅ Yes         │  ✗ Steep
│  Production use  │  ✅ Yes         │  ✓ Yes
└─────────────────────────────────────────────┘
```

---

## 🚀 Ready?

```go
import "github.com/AlexanderXinaxrZenDev/mangotemplate"

engine := mangotemplate.NewEngine("./templates")
engine.RegisterBuiltinFunctions()

// Start building!
```

---

**Happy templating! 🥭✨**

Made with ❤️ for Go developers who love the web