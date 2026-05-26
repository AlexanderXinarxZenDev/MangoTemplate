# MangoTemplate - Quick Start Guide

Get started with MangoTemplate in minutes!

## Installation

### Step 1: Add to Your Go Project

```bash
go get github.com/AlexanderXinaxrZenDev/mangotemplate
```

### Step 2: Create Template Directory

```bash
mkdir templates
```

### Step 3: Create Your First Template

Create `templates/index.mango`:

```html
<!DOCTYPE html>
<html>
<head>
	<title><go>PageTitle</go></title>
</head>
<body>
	<h1>Hello, <go upper="UserName"></go>!</h1>
	<p>Welcome to MangoTemplate</p>
</body>
</html>
```

### Step 4: Use in Your Go Code

```go
package main

import (
	"fmt"
	"log"
	"github.com/AlexanderXinaxrZenDev/mangotemplate"
)

func main() {
	// Create engine
	engine := mangotemplate.NewEngine("./templates")
	
	// Register built-in functions
	engine.RegisterBuiltinFunctions()

	// Prepare data
	data := map[string]interface{}{
		"PageTitle": "Home",
		"UserName":  "alice",
	}

	// Render template
	result, err := engine.Render("index.mango", data)
	if err != nil {
		log.Fatal(err)
	}

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
	<h1>Hello, ALICE!</h1>
	<p>Welcome to MangoTemplate</p>
</body>
</html>
```

---

## Basic Syntax (30 seconds)

### Variables

```html
<go>VariableName</go>
<go>User.Name</go>
<go>User.Address.City</go>
```

### Functions

```html
<go upper="name"></go>      <!-- UPPERCASE -->
<go lower="name"></go>      <!-- lowercase -->
<go length="items"></go>    <!-- Array length -->
<go date="timestamp"></go>  <!-- Format date -->
```

### Conditionals

```html
<go if="User.IsAdmin">
	Admin panel
</go>
<go else>
	Regular user
</go>
```

### Loops

```html
<go for="item" in="Items">
	<div><go>item.Name</go></div>
</go>
```

---

## Common Patterns

### Pattern 1: List Items

```html
<ul>
	<go for="product" in="Products">
		<li><go>product.Name</go> - $<go>product.Price</go></li>
	</go>
</ul>
```

### Pattern 2: Conditional Classes

```html
<div class="user <go if="User.Premium">premium</go>">
	<go>User.Name</go>
</div>
```

### Pattern 3: Format Numbers

```html
<p>You have <go>CartCount</go> items</p>
<p>Total: $<go>Total</go></p>
```

### Pattern 4: Default Values

```html
<p>Welcome, <go default="UserName" "Guest"></go>!</p>
```

### Pattern 5: Nested Loops

```html
<go for="category" in="Categories">
	<h2><go>category.Name</go></h2>
	<ul>
		<go for="item" in="category.Items">
			<li><go>item.Name</go></li>
		</go>
	</ul>
</go>
```

---

## Real World Example: Product Page

### Step 1: Create Template

Create `templates/product.mango`:

```html
<!DOCTYPE html>
<html>
<head>
	<title><go>product.Name</go> - Shop</title>
	<style>
		.sale { color: red; font-weight: bold; }
		.out-of-stock { color: gray; }
	</style>
</head>
<body>
	<div class="product">
		<h1><go>product.Name</go></h1>
		
		<p><go>product.Description</go></p>

		<div class="price">
			<go if="product.OnSale">
				<span class="sale">$<go>product.SalePrice</go></span>
				<span style="text-decoration: line-through;">$<go>product.Price</go></span>
			</go>
			<go else>
				<span>$<go>product.Price</go></span>
			</go>
		</div>

		<p>
			<go if="product.InStock">
				<span class="in-stock">✓ In Stock</span>
				<button>Add to Cart</button>
			</go>
			<go else>
				<span class="out-of-stock">Out of Stock</span>
			</go>
		</p>

		<div class="reviews">
			<h3>Customer Reviews (<go>product.ReviewCount</go>)</h3>
			<p>Rating: <go>product.Rating</go>/5.0</p>
			
			<go if="product.Reviews">
				<go for="review" in="product.Reviews">
					<div class="review" style="border: 1px solid #ccc; padding: 10px; margin: 10px 0;">
						<strong><go>review.Author</go></strong> (<go>review.Rating</go>★)
						<p><go>review.Text</go></p>
						<small><go date="review.Date"></go></small>
					</div>
				</go>
			</go>
		</div>
	</div>
</body>
</html>
```

### Step 2: Create Go Code

```go
package main

import (
	"fmt"
	"log"
	"time"
	"github.com/AlexanderXinaxrZenDev/mangotemplate"
)

type Review struct {
	Author string
	Rating int
	Text   string
	Date   time.Time
}

type Product struct {
	Name        string
	Description string
	Price       float64
	SalePrice   float64
	OnSale      bool
	InStock     bool
	Rating      float64
	ReviewCount int
	Reviews     []Review
}

func main() {
	engine := mangotemplate.NewEngine("./templates")
	engine.RegisterBuiltinFunctions()

	product := Product{
		Name:        "Premium Headphones",
		Description: "High-quality audio headphones with noise cancellation",
		Price:       199.99,
		SalePrice:   149.99,
		OnSale:      true,
		InStock:     true,
		Rating:      4.5,
		ReviewCount: 42,
		Reviews: []Review{
			{
				Author: "John",
				Rating: 5,
				Text:   "Excellent quality!",
				Date:   time.Now().Add(-7 * 24 * time.Hour),
			},
			{
				Author: "Sarah",
				Rating: 4,
				Text:   "Good sound, comfortable fit",
				Date:   time.Now().Add(-3 * 24 * time.Hour),
			},
		},
	}

	data := map[string]interface{}{
		"product": product,
	}

	result, err := engine.Render("product.mango", data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
```

### Output

```html
<!DOCTYPE html>
<html>
<head>
	<title>Premium Headphones - Shop</title>
	<style>
		.sale { color: red; font-weight: bold; }
		.out-of-stock { color: gray; }
	</style>
</head>
<body>
	<div class="product">
		<h1>Premium Headphones</h1>
		
		<p>High-quality audio headphones with noise cancellation</p>

		<div class="price">
			<span class="sale">$149.99</span>
			<span style="text-decoration: line-through;">$199.99</span>
		</div>

		<p>
			<span class="in-stock">✓ In Stock</span>
			<button>Add to Cart</button>
		</p>

		<div class="reviews">
			<h3>Customer Reviews (42)</h3>
			<p>Rating: 4.5/5.0</p>
			
			<div class="review" style="border: 1px solid #ccc; padding: 10px; margin: 10px 0;">
				<strong>John</strong> (5★)
				<p>Excellent quality!</p>
				<small>2024-05-20</small>
			</div>

			<div class="review" style="border: 1px solid #ccc; padding: 10px; margin: 10px 0;">
				<strong>Sarah</strong> (4★)
				<p>Good sound, comfortable fit</p>
				<small>2024-05-24</small>
			</div>
		</div>
	</div>
</body>
</html>
```

---

## Web Server Example

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
			"Message": "Welcome!",
		}

		result, err := engine.Render("index.mango", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, result)
	})

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## Tips & Tricks

### Tip 1: Use Meaningful Variable Names

```html
<!-- Good -->
<go>User.FirstName</go>

<!-- Bad -->
<go>u</go>
```

### Tip 2: Organize Templates in Subdirectories

```
templates/
├── layouts/
│   └── base.mango
├── pages/
│   ├── home.mango
│   └── about.mango
└── components/
    ├── header.mango
    └── footer.mango
```

### Tip 3: Use Debug Mode During Development

```go
engine.EnableDebug()
```

### Tip 4: Cache Templates in Production

Templates are cached automatically. Monitor with:

```go
stats := engine.GetCacheStats()
fmt.Printf("Cached: %v\n", stats["cached_templates"])
```

### Tip 5: Register Custom Functions

```go
engine.RegisterFunc("discount", func(price, percent float64) float64 {
	return price * (1 - percent/100)
})
```

---

## Troubleshooting

### Issue: Template not found

```
Error: failed to load template: open ./templates/xyz.mango: no such file or directory
```

**Solution**: Check file path and extension

```bash
ls -la templates/  # Verify template exists
```

### Issue: Syntax error in template

```
Error: template parsing error: ...
```

**Solution**: Verify <go> tag syntax

```html
<!-- Correct -->
<go>Name</go>
<go if="Condition">...</go>

<!-- Wrong -->
<go>Name</go >
<go if='Condition'>...</go>
```

### Issue: Variable not showing

**Solution**: Ensure data is passed with correct key

```go
// Template expects: {{ .Title }}
data := map[string]interface{}{
	"Title": "My Title",  // Key must match
}
```

---

## Next Steps

1. **Read Full Documentation**: Check `MANGOTEMPLATE_DOCUMENTATION.md`
2. **Explore Examples**: Check `examples.go` for more patterns
3. **Build Your App**: Start with `example6WebServer()`
4. **Register Custom Functions**: Add your own template functions

---

## File Structure Best Practice

```
project/
├── main.go
├── templates/
│   ├── layouts/
│   │   └── base.mango
│   ├── pages/
│   │   ├── home.mango
│   │   ├── products.mango
│   │   └── contact.mango
│   ├── components/
│   │   ├── header.mango
│   │   ├── footer.mango
│   │   └── nav.mango
│   └── email/
│       ├── welcome.mango
│       └── confirmation.mango
└── public/
    ├── css/
    ├── js/
    └── images/
```

---

## Performance Tips

1. **Reuse Engine**: Create one engine instance for your app
2. **Enable Caching**: It's on by default
3. **Pre-compile**: Templates are compiled on first use
4. **Monitor Cache**: Use `GetCacheStats()` to monitor

---

**You're ready to go! Start building with MangoTemplate! 🥭**