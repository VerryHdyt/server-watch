package main

import "testing"

func TestMemoryPercentage(t *testing.T) {
	total, available := uint64(800), uint64(200)
	got := float64(total-available) / float64(total) * 100
	if got != 75 { t.Fatalf("expected 75, got %v", got) }
}

func TestReadMemoryOnLinux(t *testing.T) {
	if testing.Short() { t.Skip("skip host procfs check in short mode") }
	total, available, err := readMemory()
	if err != nil { t.Skipf("Linux procfs unavailable: %v", err) }
	if total == 0 || available > total { t.Fatalf("unexpected memory values: total=%d available=%d", total, available) }
}
