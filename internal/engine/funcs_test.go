package engine

import "testing"

func TestCaseConversions(t *testing.T) {
	cases := []struct {
		in, pascal, camel, snake, kebab string
	}{
		{"user_profile", "UserProfile", "userProfile", "user_profile", "user-profile"},
		{"BlogPost", "BlogPost", "blogPost", "blog_post", "blog-post"},
		{"HTTPServer", "HTTPServer", "httpServer", "http_server", "http-server"},
		{"order-item", "OrderItem", "orderItem", "order_item", "order-item"},
		{"id", "ID", "id", "id", "id"},
	}
	for _, c := range cases {
		if got := Pascal(c.in); got != c.pascal {
			t.Errorf("Pascal(%q)=%q want %q", c.in, got, c.pascal)
		}
		if got := Camel(c.in); got != c.camel {
			t.Errorf("Camel(%q)=%q want %q", c.in, got, c.camel)
		}
		if got := Snake(c.in); got != c.snake {
			t.Errorf("Snake(%q)=%q want %q", c.in, got, c.snake)
		}
		if got := Kebab(c.in); got != c.kebab {
			t.Errorf("Kebab(%q)=%q want %q", c.in, got, c.kebab)
		}
	}
}

func TestPluralSingular(t *testing.T) {
	cases := []struct{ singular, plural string }{
		{"post", "posts"},
		{"category", "categories"},
		{"box", "boxes"},
		{"church", "churches"},
		{"person", "people"},
		{"User", "Users"},
	}
	for _, c := range cases {
		if got := Plural(c.singular); got != c.plural {
			t.Errorf("Plural(%q)=%q want %q", c.singular, got, c.plural)
		}
		if got := Singular(c.plural); got != c.singular {
			t.Errorf("Singular(%q)=%q want %q", c.plural, got, c.singular)
		}
	}
}
