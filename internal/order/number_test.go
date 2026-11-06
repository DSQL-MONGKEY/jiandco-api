package order

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateOrderNumber(t *testing.T) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}

	number, err := generateOrderNumber(location)
	if err != nil {
		t.Fatalf("generate order number: %v", err)
	}
	if !strings.HasPrefix(number, "JDC-") {
		t.Fatalf("order number = %q, want JDC- prefix", number)
	}
	parts := strings.Split(number, "-")
	if len(parts) != 3 {
		t.Fatalf("order number = %q, want 3 parts", number)
	}
	if len(parts[1]) != 8 {
		t.Fatalf("date part = %q, want 8 digits", parts[1])
	}
	if len(parts[2]) != 6 {
		t.Fatalf("random part = %q, want 6 chars", parts[2])
	}
}
