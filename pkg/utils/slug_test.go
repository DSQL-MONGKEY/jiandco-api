package utils

import "testing"

func TestSlugify(t *testing.T) {
	tests := map[string]string{
		"Kaos Oversize Black":      "kaos-oversize-black",
		"  Hello, World!  ":        "hello-world",
		"Sepatu---Running":         "sepatu-running",
		"":                         "item",
		"T-Shirt / Crew Neck 2026": "t-shirt-crew-neck-2026",
	}

	for input, want := range tests {
		if got := Slugify(input); got != want {
			t.Fatalf("Slugify(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestUniqueSlug(t *testing.T) {
	existing := map[string]bool{"kaos": true, "kaos-2": true}
	slug, err := UniqueSlug("Kaos", func(candidate string) (bool, error) {
		return existing[candidate], nil
	})
	if err != nil {
		t.Fatalf("UniqueSlug: %v", err)
	}
	if slug != "kaos-3" {
		t.Fatalf("slug = %q, want kaos-3", slug)
	}
}
