package accounting

import "testing"

func TestBuildItems(t *testing.T) {
	items, total, err := buildItems([]ItemRequest{
		{Name: "Kain", Quantity: 2.5, Unit: "meter", UnitPrice: 50000},
		{Name: "Benang", Quantity: 3, Unit: "pcs", UnitPrice: 10000},
	}, 0)
	if err != nil {
		t.Fatalf("buildItems: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("items = %d, want 2", len(items))
	}
	if total != 155000 {
		t.Fatalf("total = %d, want 155000", total)
	}
	if items[0].Subtotal != 125000 {
		t.Fatalf("first subtotal = %d, want 125000", items[0].Subtotal)
	}
}

func TestBuildItemsManualTotal(t *testing.T) {
	items, total, err := buildItems(nil, 75000)
	if err != nil {
		t.Fatalf("buildItems: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("items = %d, want 0", len(items))
	}
	if total != 75000 {
		t.Fatalf("total = %d, want 75000", total)
	}
}

func TestBuildItemsRejectsZeroTotal(t *testing.T) {
	if _, _, err := buildItems(nil, 0); err == nil {
		t.Fatal("expected zero total error")
	}
}
