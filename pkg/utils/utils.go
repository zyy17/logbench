package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

var sizeUnitToBytes = map[string]int64{
	// IEC units.
	"bytes": 1,
	"kib":   1024,
	"mib":   1024 * 1024,
	"gib":   1024 * 1024 * 1024,
	"tib":   1024 * 1024 * 1024 * 1024,
	"pib":   1024 * 1024 * 1024 * 1024 * 1024,

	// SI units.
	"kb": 1000,
	"mb": 1000 * 1000,
	"gb": 1000 * 1000 * 1000,
	"tb": 1000 * 1000 * 1000 * 1000,
	"pb": 1000 * 1000 * 1000 * 1000 * 1000,
}

// ParseSize parses a size string and returns the size in bytes.
// The size string can be in format of "100bytes", "1kib", "1mib", etc.
func ParseSize(size string) (int64, error) {
	if size == "" {
		return 0, fmt.Errorf("size is empty")
	}

	// Normalize the size string.
	size = strings.ToLower(strings.TrimSpace(size))

	// unitStart is the index of the unit string.
	unitStart := 0
	for unitStart = 0; unitStart < len(size); unitStart++ {
		if size[unitStart] >= '0' && size[unitStart] <= '9' {
			continue
		}

		break
	}

	if unitStart == len(size) {
		return 0, fmt.Errorf("invalid size: %s", size)
	}

	numStr := size[:unitStart]
	unitStr := strings.TrimSpace(size[unitStart:])

	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size: %s", size)
	}

	unit, ok := sizeUnitToBytes[unitStr]
	if !ok {
		return 0, fmt.Errorf("invalid size: %s", size)
	}

	return num * unit, nil
}

func ParsePercentage(percentage string) (float64, error) {
	if percentage == "" {
		return 0, fmt.Errorf("percentage is empty")
	}

	parts := strings.Split(percentage, "%")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid percentage: %s", percentage)
	}

	num, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid percentage: %s", percentage)
	}

	return num / 100, nil
}

// Generate a random number in [min, max).
func RandomNumber(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}
