package core_utils

import (
	"fmt"
)

func ParseFloat(s string) float64 {
	var f float64

	if _, err := fmt.Sscanf(s, "%f", &f); err != nil {
		return 0
	}

	return f
}
