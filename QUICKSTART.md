# MangoTemplate - Go Templating Engine

A powerful, PHP-like templating engine for Go that supports `.mango` and `.html` files with seamless Go integration.

## Table of Contents

1. [Installation](#installation)
2. [Quick Start](#quick-start)
3. [Syntax](#syntax)
4. [Built-in Functions](#built-in-functions)
5. [Examples](#examples)
6. [API Reference](#api-reference)
7. [Best Practices](#best-practices)

---

## Installation

```bash
go get github.com/AlexanderXinaxrZenDev/mangotemplate
```

---

## Quick Start

### Basic Setup

```go
package main

import (
	"fmt"
	"log"
	"github.com/AlexanderXinaxrZenDev/mangotemplate"
)

func main() {
	// Create engine with template directory
	engine := mangotemplate.NewEngine("./templates")
	
	// Register built-in functions
	engine.RegisterBuiltinFunctions()
	
	// Render a template
	data := map[string]interface{}{
		"Title": "Hello World",
		"Name":  "Alice",
		"Age":   25,
	}
	
	result, err := engine.Render("index.mango", data)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(result)
}
```

### Simple Template (index.mango)

```html
<!DOCTYPE html>
<html>
<head>
	<title><go>Title</go></title>
</head>
<body>
	<h1>Hello, <go>Name</go>!</h1>
	<p>You are <go>Age</go> years old.</p>
</body>
</html>
```

### Output

```html
<!DOCTYPE html>
<html>
<head>
	<title>Hello World</title>
</head>
<body>
	<h1>Hello, Alice!</h1>
	<p>You are 25 years old.</p>
</body>
</html>
```

---

## Syntax

### 1. Variable Output

```html
<!-- Output a variable -->
<go>VariableName</go>

<!-- With dot notation for nested access -->
<go>User.Name</go>
<go>User.Address.City</go>

<!-- With functions -->
<go upper="Name"></go>     <!-- Uppercase -->
<go lower="Name"></go>     <!-- Lowercase -->
<go title="Name"></go>     <!-- Title case -->
<go truncate="Name" 10></go>  <!-- Truncate to 10 chars -->
```

### 2. Conditionals

```html
<!-- If statement -->
<go if="Condition">
	Content if true
</go>

<!-- If-else -->
<go if="User.IsAdmin">
	Admin panel
</go>
<go else>
	Regular user
</go>

<!-- Multiple conditions -->
<go if="User.Age > 18">
	Adult content
</go>
<go else>
	Youth content
</go>

<!-- Negation -->
<go if="!User.IsActive">
	Inactive account
</go>
```

### 3. Loops

```html
<!-- For loop with range -->
<go for="item" in="Items">
	<div><go>item.Name</go> - $<go>item.Price</go></div>
</go>

<!-- With dot notation -->
<go for="user" in="Users">
	<li>
		<go>user.Name</go> (<go>user.Email</go>)
	</li>
</go>

<!-- Nested loops -->
<go for="category" in="Categories">
	<h2><go>category.Name</go></h2>
	<go for="product" in="category.Products">
		<span><go>product.Name</go></span>
	</go>
</go>
```

### 4. Functions

```html
<!-- String functions -->
<go upper="Name"></go>           <!-- UPPERCASE -->
<go lower="Name"></go>           <!-- lowercase -->
<go capitalize="Name"></go>      <!-- Capitalize -->
<go truncate="Text" 50></go>     <!-- Truncate -->
<go slug="Title"></go>           <!-- URL slug -->

<!-- Math functions -->
<go add "10" "5"></go>           <!-- 15 -->
<go multiply "10" "5"></go>      <!-- 50 -->
<go max "10" "20"></go>          <!-- 20 -->
<go min "10" "20"></go>          <!-- 10 -->

<!-- Date functions -->
<go date="CreatedDate"></go>     <!-- Format as date -->
<go datetime="CreatedDate"></go> <!-- Format as datetime -->
<go time="CreatedDate"></go>     <!-- Format as time -->

<!-- Array functions -->
<go first="Items"></go>          <!-- First element -->
<go last="Items"></go>           <!-- Last element -->
<go length="Items"></go>         <!-- Array length -->

<!-- Type functions -->
<go int="Price"></go>            <!-- Convert to int -->
<go float="Price"></go>          <!-- Convert to float -->
<go string="Age"></go>           <!-- Convert to string -->
```

### 5. Variable Assignment

```html
<!-- Assign variable (using raw Go template syntax inside tags) -->
<go>
	{{ $sum := add .Price .Tax }}
	Total: ${{ $sum }}
</go>

<!-- Using helper syntax -->
<go var="greeting" value="Hello"></go>
<p><go>greeting</go> World</p>
```

### 6. Raw Go Template Code

For complex logic, use raw Go template syntax:

```html
<go>
	{{ if eq .Status "active" }}
		<span class="active">Active</span>
	{{ else }}
		<span class="inactive">Inactive</span>
	{{ end }}
</go>

<go>
	{{ range .Products }}
		<div>{{ .Name }}: ${{ .Price }}</div>
	{{ end }}
</go>
```

---

## Built-in Functions

### String Functions

| Function | Usage | Example |
|----------|-------|---------|
| `upper` | Convert to uppercase | `<go upper="name"></go>` |
| `lower` | Convert to lowercase | `<go lower="name"></go>` |
| `title` | Title case | `<go title="name"></go>` |
| `capitalize` | Capitalize first letter | `<go capitalize="name"></go>` |
| `reverse` | Reverse string | `<go reverse="text"></go>` |
| `truncate` | Truncate string | `<go truncate="text" 50></go>` |
| `slug` | URL-friendly slug | `<go slug="title"></go>` |
| `strlen` | String length | `<go strlen="text"></go>` |
| `substr` | Substring | `<go substr="text" 0 10></go>` |
| `contains` | Check if contains | `<go contains="text" "hello"></go>` |
| `replace` | Replace text | `<go replace="text" "old" "new"></go>` |
| `split` | Split string | `<go split="csv" ","></go>` |
| `join` | Join array | `<go join="items" ", "></go>` |

### Math Functions

| Function | Usage | Example |
|----------|-------|---------|
| `add` | Addition | `<go add="a" "b"></go>` |
| `subtract` | Subtraction | `<go subtract="a" "b"></go>` |
| `multiply` | Multiplication | `<go multiply="a" "b"></go>` |
| `divide` | Division | `<go divide="a" "b"></go>` |
| `mod` | Modulo | `<go mod="a" "b"></go>` |
| `min` | Minimum | `<go min="a" "b"></go>` |
| `max` | Maximum | `<go max="a" "b"></go>` |
| `abs` | Absolute value | `<go abs="number"></go>` |

### Date/Time Functions

| Function | Usage | Example |
|----------|-------|---------|
| `date` | Format as date (YYYY-MM-DD) | `<go date="timestamp"></go>` |
| `time` | Format as time (HH:MM:SS) | `<go time="timestamp"></go>` |
| `datetime` | Format as datetime | `<go datetime="timestamp"></go>` |
| `unix` | Convert to Unix timestamp | `<go unix="timestamp"></go>` |
| `now` | Current time | `<go>now</go>` |

### Array/Collection Functions

| Function | Usage | Example |
|----------|-------|---------|
| `first` | First element | `<go first="items"></go>` |
| `last` | Last element | `<go last="items"></go>` |
| `length` | Array length | `<go length="items"></go>` |
| `range` | Create range | `<go range="1" "10"></go>` |

### Conditional Functions

| Function | Usage | Example |
|----------|-------|---------|
| `default` | Default value | `<go default="name" "Unknown"></go>` |
| `empty` | Check if empty | `<go if="empty name">...</go>` |
| `eq` | Equals | `<go if="eq status 'active'">...</go>` |
| `ne` | Not equals | `<go if="ne status 'active'">...</go>` |
| `lt` | Less than | `<go if="lt age 18">...</go>` |
| `gt` | Greater than | `<go if="gt age 18">...</go>` |
| `lte` | Less than or equal | `<go if="lte age 18">...</go>` |
| `gte` | Greater than or equal | `<go if="gte age 18">...</go>` |

### Type Conversion

| Function | Usage | Example |
|----------|-------|---------|
| `int` | Convert to int | `<go int="price"></go>` |
| `float` | Convert to float | `<go float="price"></go>` |
| `string` | Convert to string | `<go string="age"></go>` |
| `bool` | Convert to bool | `<go bool="active"></go>` |
| `type` | Get type | `<go type="value"></go>` |

### Security/Encoding

| Function | Usage | Example |
|----------|-------|---------|
| `escape` | HTML escape | `<go escape="text"></go>` |
| `unescape` | HTML unescape | `<go unescape="html"></go>` |
| `safe` | Mark as safe | `<go safe="html"></go>` |
| `md5` / `hash` | MD5 hash | `<go md5="text"></go>` |
| `urlencode` | URL encode | `<go urlencode="text"></go>` |

---

## Examples

### Example 1: Blog Post

**blog.mango**
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
			By <go>post.Author</go> 
			on <go date="post.PublishedDate"></go>
		</p>
		<p class="content"><go>post.Content</go></p>
		<footer>
			Tags: 
			<go for="tag" in="post.Tags">
				<a href="/tag/<go slug=tag></go>"><go>tag</go></a>
			</go>
		</footer>
	</article>
</body>
</html>
```

**Go Code**
```go
data := map[string]interface{}{
	"post": map[string]interface{}{
		"Title":        "My First Post",
		"Author":       "John Doe",
		"Content":      "Lorem ipsum dolor sit amet...",
		"PublishedDate": time.Now(),
		"Tags":         []string{"golang", "web", "templates"},
	},
}

result, _ := engine.Render("blog.mango", data)
```

### Example 2: Product Listing

**products.mango**
```html
<div class="products">
	<go if="empty products">
		<p>No products available.</p>
	</go>
	<go else>
		<go for="product" in="products">
			<div class="product-card">
				<h3><go>product.Name</go></h3>
				<p><go truncate="product.Description" 100></go></p>
				<div class="price">
					<go if="product.OnSale">
						<span class="original">$<go>product.OriginalPrice</go></span>
						<span class="sale">$<go>product.Price</go></span>
					</go>
					<go else>
						<span class="price">$<go>product.Price</go></span>
					</go>
				</div>
				<button onclick="addToCart(<go>product.ID</go>)">
					Add to Cart
				</button>
			</div>
		</go>
	</go>
</div>
```

### Example 3: User Dashboard

**dashboard.mango**
```html
<!DOCTYPE html>
<html>
<head>
	<title><go>user.Name</go>'s Dashboard</title>
</head>
<body>
	<header>
		<h1>Welcome, <go>user.Name</go>!</h1>
	</header>

	<main>
		<section class="stats">
			<div class="stat-card">
				<h3>Total Orders</h3>
				<p class="number"><go>user.TotalOrders</go></p>
			</div>
			<div class="stat-card">
				<h3>Total Spent</h3>
				<p class="number">$<go>user.TotalSpent</go></p>
			</div>
		</section>

		<section class="recent-orders">
			<h2>Recent Orders</h2>
			<go if="empty user.Orders">
				<p>You haven't placed any orders yet.</p>
			</go>
			<go else>
				<table>
					<thead>
						<tr>
							<th>Order ID</th>
							<th>Date</th>
							<th>Total</th>
							<th>Status</th>
						</tr>
					</thead>
					<tbody>
						<go for="order" in="user.Orders">
							<tr>
								<td><go>order.ID</go></td>
								<td><go date="order.CreatedAt"></go></td>
								<td>$<go>order.Total</go></td>
								<td>
									<span class="status status-<go lower=order.Status></go>">
										<go>order.Status</go>
									</span>
								</td>
							</tr>
						</go>
					</tbody>
				</table>
			</go>
		</section>
	</main>
</body>
</html>
```

### Example 4: Email Template

**email.mango**
```html
<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif;">
	<div style="max-width: 600px; margin: 0 auto;">
		<h2>Hello <go>recipient.FirstName</go>,</h2>
		
		<p><go>emailBody</go></p>

		<go if="items">
			<table style="width: 100%; border-collapse: collapse;">
				<tr style="border-bottom: 2px solid #333;">
					<th style="text-align: left;">Item</th>
					<th>Qty</th>
					<th>Price</th>
				</tr>
				<go for="item" in="items">
					<tr style="border-bottom: 1px solid #eee;">
						<td><go>item.Name</go></td>
						<td style="text-align: center;"><go>item.Qty</go></td>
						<td style="text-align: right;">$<go>item.Price</go></td>
					</tr>
				</go>
			</table>
			<p style="font-weight: bold;">Total: $<go>totalAmount</go></p>
		</go>

		<p>
			<a href="<go>actionLink</go>" 
			   style="background: #007bff; color: white; padding: 10px 20px; 
			           text-decoration: none; border-radius: 5px;">
				<go>actionText</go>
			</a>
		</p>

		<hr>
		<p style="color: #666; font-size: 12px;">
			&copy; <go>year></go> <go>companyName</go>. All rights reserved.
		</p>
	</div>
</body>
</html>
```

---

## API Reference

### Engine Methods

```go
// Create new engine
engine := mangotemplate.NewEngine(templateDir)

// Register functions
engine.RegisterFunc(name string, fn interface{})
engine.RegisterBuiltinFunctions()

// Load templates
template, err := engine.LoadTemplate(filename)
template, err := engine.GetTemplate(filename)
exists := engine.TemplateExists(filename)

// Render templates
result, err := engine.Render(filename, data)
result, err := engine.RenderString(content, data)
result, err := engine.RenderLayout(layoutFile, contentFile, data)

// Cache management
engine.ClearCache()
engine.ClearCacheFile(filename)
stats := engine.GetCacheStats()

// Configuration
engine.SetDelimiters(open, close)
engine.SetTemplateDir(dir)
engine.EnableDebug()
engine.DisableDebug()

// Include templates
result, err := engine.Include(filename, data)
```

### Template Methods

```go
// Parse template
err := template.Parse()

// Render
result, err := template.Render(data)
err := template.RenderToWriter(&buf, data)
```

---

## Best Practices

1. **Use Template Caching**: Templates are cached automatically
2. **Enable Debug Mode**: During development, enable debug for better error messages
3. **Organize Templates**: Use subdirectories for better organization
4. **Data Structure**: Pass well-structured data to templates
5. **Security**: Always escape user input - use the `escape` function
6. **Performance**: Pre-compile templates in production

### Example: Web Server

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/AlexanderXinaxrZenDev/mangotemplate"
)

func main() {
	engine := mangotemplate.NewEngine("./templates")
	engine.RegisterBuiltinFunctions()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title":   "Home",
			"Message": "Welcome to MangoTemplate!",
		}

		result, err := engine.Render("index.mango", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, result)
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## File Format Support

- `.mango` - MangoTemplate files (recommended)
- `.html` - Standard HTML files (Go template syntax only)

---

## License

MIT License

---

**Happy templating with MangoTemplate! 🥭**