package utils

import (
	"testing"
	"time"
)

func TestParseMonth(t *testing.T) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}

	start, end, err := ParseMonth("2026-11", location)
	if err != nil {
		t.Fatalf("parse month: %v", err)
	}

	wantStart := time.Date(2026, time.November, 1, 0, 0, 0, 0, location)
	wantEnd := time.Date(2026, time.December, 1, 0, 0, 0, 0, location)
	if !start.Equal(wantStart) {
		t.Fatalf("start = %s, want %s", start, wantStart)
	}
	if !end.Equal(wantEnd) {
		t.Fatalf("end = %s, want %s", end, wantEnd)
	}
}

func TestParseMonthDefaultsToCurrentMonth(t *testing.T) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}

	start, end, err := ParseMonth("", location)
	if err != nil {
		t.Fatalf("parse month: %v", err)
	}
	if start.Day() != 1 {
		t.Fatalf("start day = %d, want 1", start.Day())
	}
	if !end.Equal(start.AddDate(0, 1, 0)) {
		t.Fatalf("end = %s, want %s", end, start.AddDate(0, 1, 0))
	}
}

func TestParseMonthRejectsInvalidValue(t *testing.T) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}

	if _, _, err := ParseMonth("2026-13", location); err == nil {
		t.Fatal("expected invalid month error")
	}
}
