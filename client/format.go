package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func formatBits(i uint64, precision uint8) string {
	const unit = 1024

	if i < unit {
		return fmt.Sprintf("%d b", i)
	}

	div, exp := int64(unit), 0
	for n := i / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	format := fmt.Sprintf("%%.%df %%cb", precision)
	return fmt.Sprintf(format, float64(i)/float64(div), "KMGTPE"[exp])
}

func formatBytes(i uint64, precision uint8) string {
	const unit = 1024

	if i < unit {
		return fmt.Sprintf("%d B", i)
	}

	div, exp := int64(unit), 0
	for n := i / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	format := fmt.Sprintf("%%.%df %%cB", precision)
	return fmt.Sprintf(format, float64(i)/float64(div), "KMGTPE"[exp])
}

func formatPercent(total, used float64, precision uint8) string {
	percent := (100 / total) * used

	format := fmt.Sprintf("%%.%df%%", precision)

	return fmt.Sprintf(format, percent)
}

func formatDuration(duration time.Duration) string {
	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		singularName string
		amount       int64
	}{
		{"day", days},
		{"hour", hours},
		{"minute", minutes},
		{"second", seconds},
	}

	parts := []string{}

	for _, chunk := range chunks {
		switch chunk.amount {
		case 0:
			continue
		case 1:
			parts = append(parts, fmt.Sprintf("%d %s", chunk.amount, chunk.singularName))
		default:
			parts = append(parts, fmt.Sprintf("%d %ss", chunk.amount, chunk.singularName))
		}
	}

	return strings.Join(parts, " ")
}
