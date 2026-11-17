package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type Metrics struct {
	Hostname       string  `json:"hostname"`
	MemoryTotalMB  uint64  `json:"memory_total_mb"`
	MemoryUsedMB   uint64  `json:"memory_used_mb"`
	MemoryUsedPct  float64 `json:"memory_used_percent"`
	DiskTotalMB    uint64  `json:"disk_total_mb"`
	DiskUsedMB     uint64  `json:"disk_used_mb"`
	DiskUsedPct    float64 `json:"disk_used_percent"`
}

type Collector interface {
	Collect() (Metrics, error)
}

type ProcCollector struct{}

func (ProcCollector) Collect() (Metrics, error) {
	host, err := os.Hostname()
	if err != nil { return Metrics{}, fmt.Errorf("hostname: %w", err) }
	total, available, err := readMemory()
	if err != nil { return Metrics{}, err }
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil { return Metrics{}, fmt.Errorf("disk: %w", err) }
	diskTotal := stat.Blocks * uint64(stat.Bsize)
	diskFree := stat.Bavail * uint64(stat.Bsize)
	diskUsed := diskTotal - diskFree
	m := Metrics{Hostname: host, MemoryTotalMB: total / 1024, MemoryUsedMB: (total-available) / 1024,
		DiskTotalMB: diskTotal / 1024 / 1024, DiskUsedMB: diskUsed / 1024 / 1024}
	if total > 0 { m.MemoryUsedPct = float64(total-available) / float64(total) * 100 }
	if diskTotal > 0 { m.DiskUsedPct = float64(diskUsed) / float64(diskTotal) * 100 }
	return m, nil
}

func readMemory() (total, available uint64, err error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil { return 0, 0, fmt.Errorf("memory: %w", err) }
	defer f.Close()
	values := map[string]uint64{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 { continue }
		if fields[0] == "MemTotal:" || fields[0] == "MemAvailable:" {
			v, parseErr := strconv.ParseUint(fields[1], 10, 64)
			if parseErr != nil { return 0, 0, fmt.Errorf("parse memory: %w", parseErr) }
			values[strings.TrimSuffix(fields[0], ":")] = v
		}
	}
	if err := scanner.Err(); err != nil { return 0, 0, err }
	total, okTotal := values["MemTotal"]
	available, okAvailable := values["MemAvailable"]
	if !okTotal || !okAvailable { return 0, 0, fmt.Errorf("required memory fields missing") }
	return total, available, nil
}
