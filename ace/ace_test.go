package ace

import (
	"bytes"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
)

func trim(str string) string {
	trimmed := strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(str, " "))
	trimmed = strings.Replace(trimmed, " <", "<", -1)
	trimmed = strings.Replace(trimmed, "> ", ">", -1)
	return trimmed
}

func Test_Render(t *testing.T) {
	engine := New("./views", ".ace")
	engine.AddFunc("isAdmin", func(user string) bool {
		return user == "admin"
	})
	if err := engine.Load(); err != nil {
		t.Fatalf("load: %v\n", err)
	}
	// Partials
	var buf bytes.Buffer
	engine.Render(&buf, "index", map[string]interface{}{
		"Title": "Hello, World!",
	})
	expect := `<h2>Header</h2><h1>Hello, World!</h1><h2>Footer</h2>`
	result := trim(buf.String())
	if expect != result {
		t.Fatalf("Expected:\n%s\nResult:\n%s\n", expect, result)
	}
	// Single
	buf.Reset()
	engine.Render(&buf, "errors/404", map[string]interface{}{
		"Title": "Hello, World!",
	})
	expect = `<h1>Hello, World!</h1>`
	result = trim(buf.String())
	if expect != result {
		t.Fatalf("Expected:\n%s\nResult:\n%s\n", expect, result)
	}
}

func Test_Layout(t *testing.T) {
	engine := New("./views", ".ace")
	engine.AddFunc("isAdmin", func(user string) bool {
		return user == "admin"
	})
	engine.Debug(true)
	if err := engine.Load(); err != nil {
		t.Fatalf("load: %v\n", err)
	}

	var buf bytes.Buffer
	err := engine.Render(&buf, "index", map[string]interface{}{
		"Title": "Hello, World!",
	}, "layouts/main")
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	expect := `<!DOCTYPE html><html><head><title>Main</title></head><body><h2>Header</h2><h1>Hello, World!</h1><h2>Footer</h2></body></html>`
	result := trim(buf.String())
	if expect != result {
		t.Fatalf("Expected:\n%s\nResult:\n%s\n", expect, result)
	}
}

func Test_Empty_Layout(t *testing.T) {
	engine := New("./views", ".ace")
	engine.AddFunc("isAdmin", func(user string) bool {
		return user == "admin"
	})
	engine.Debug(true)
	if err := engine.Load(); err != nil {
		t.Fatalf("load: %v\n", err)
	}

	var buf bytes.Buffer
	err := engine.Render(&buf, "index", map[string]interface{}{
		"Title": "Hello, World!",
	}, "")
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	expect := `<h2>Header</h2><h1>Hello, World!</h1><h2>Footer</h2>`
	result := trim(buf.String())
	if expect != result {
		t.Fatalf("Expected:\n%s\nResult:\n%s\n", expect, result)
	}
}

func Test_FileSystem(t *testing.T) {
	engine := NewFileSystem(http.Dir("./views"), ".ace")
	engine.AddFunc("isAdmin", func(user string) bool {
		return user == "admin"
	})
	engine.Debug(true)
	if err := engine.Load(); err != nil {
		t.Fatalf("load: %v\n", err)
	}

	var buf bytes.Buffer
	err := engine.Render(&buf, "index", map[string]interface{}{
		"Title": "Hello, World!",
	}, "layouts/main")
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	expect := `<!DOCTYPE html><html><head><title>Main</title></head><body><h2>Header</h2><h1>Hello, World!</h1><h2>Footer</h2></body></html>`
	result := trim(buf.String())
	if expect != result {
		t.Fatalf("Expected:\n%s\nResult:\n%s\n", expect, result)
	}
}

//goland:noinspection GoDeprecation
func Test_Reload(t *testing.T) {
	engine := NewFileSystem(http.Dir("./views"), ".ace")
	engine.AddFunc("isAdmin", func(user string) bool {
		return user == "admin"
	})
	engine.Reload(true) // Optional. Default: false

	engine.AddFunc("isAdmin", func(user string) bool {
		return user == "admin"
	})
	if err := engine.Load(); err != nil {
		t.Fatalf("load: %v\n", err)
	}

	if err := os.WriteFile("./views/reload.ace", []byte("after reload\n"), 0644); err != nil {
		t.Fatalf("write file: %v\n", err)
	}
	defer func() {
		if err := os.WriteFile("./views/reload.ace", []byte("before reload\n"), 0644); err != nil {
			t.Fatalf("write file: %v\n", err)
		}
	}()

	engine.Load()

	var buf bytes.Buffer
	engine.Render(&buf, "reload", nil)
	expect := "<after>reload</after>"
	result := trim(buf.String())
	if expect != result {
		t.Fatalf("Expected:\n%s\nResult:\n%s\n", expect, result)
	}
}

func Test_AddFuncMap(t *testing.T) {
	// Create a temporary directory
	dir, _ := os.MkdirTemp(".", "")
	defer os.RemoveAll(dir)

	// Create a temporary template file.
	_ = os.WriteFile(dir+"/func_map.ace", []byte(`
h2 {{lower .Var1}}
p {{upper .Var2}}`), 0700)

	engine := New(dir, ".ace")

	fm := map[string]interface{}{
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
	}

	engine.AddFuncMap(fm)

	if err := engine.Load(); err != nil {
		t.Fatalf("load: %v\n", err)
	}

	var buf bytes.Buffer
	engine.Render(&buf, "func_map", map[string]interface{}{
		"Var1": "LOwEr",
		"Var2": "upPEr",
	})
	expect := `<h2>lower</h2><p>UPPER</p>`
	result := trim(buf.String())
	if expect != result {
		t.Fatalf("Expected:\n%s\nResult:\n%s\n", expect, result)
	}

	// FuncMap
	fm2 := engine.FuncMap()
	if _, ok := fm2["lower"]; !ok {
		t.Fatalf("Function lower does not exist in FuncMap().\n")
	}
	if _, ok := fm2["upper"]; !ok {
		t.Fatalf("Function upper does not exist in FuncMap().\n")
	}
}

func Benchmark_Ace(b *testing.B) {
	expectSimple := `<h1>Hello, World!</h1>`
	expectExtended := `<!DOCTYPE html><html><head><title>Main</title></head><body><h2>Header</h2><h1>Hello, Admin!</h1><h2>Footer</h2></body></html>`
	engine := New("./views", ".ace")
	engine.AddFunc("isAdmin", func(user string) bool {
		return user == "admin"
	})
	var buf bytes.Buffer
	var err error

	b.Run("simple", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		for i := 0; i < bb.N; i++ {
			buf.Reset()
			err = engine.Render(&buf, "simple", map[string]interface{}{
				"Title": "Hello, World!",
			})
		}

		if err != nil {
			bb.Fatalf("Failed to render: %v", err)
		}
		result := trim(buf.String())
		if expectSimple != result {
			bb.Fatalf("Expected:\n%s\nResult:\n%s\n", expectSimple, result)
		}
	})

	b.Run("extended", func(bb *testing.B) {
		bb.ReportAllocs()
		bb.ResetTimer()
		for i := 0; i < bb.N; i++ {
			buf.Reset()
			err = engine.Render(&buf, "extended", map[string]interface{}{
				"User": "admin",
			}, "layouts/main")
		}

		if err != nil {
			bb.Fatalf("Failed to render: %v", err)
		}
		result := trim(buf.String())
		if expectExtended != result {
			bb.Fatalf("Expected:\n%s\nResult:\n%s\n", expectExtended, result)
		}
	})
}
