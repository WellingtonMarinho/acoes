package pricefeed

import (
	"strconv"
	"strings"
)

func parseDecimal(raw string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(raw), 64)
}
