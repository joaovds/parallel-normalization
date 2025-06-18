package handlecsv

import (
	"strconv"
	"strings"
)

func IdentifyCategoricalColumns(lines [][]string) []int {
	if len(lines) == 0 {
		return nil
	}

	totalCols := len(lines[0])
	var categoricalCols []int

	for i := range totalCols {
		isNumeric := true

		for _, line := range lines {
			if i >= len(line) {
				continue
			}

			value := strings.TrimSpace(line[i])
			if value == "" {
				break
			}

			_, err := strconv.ParseFloat(value, 64)
			if err != nil {
				isNumeric = false
				break
			}
		}

		if !isNumeric {
			categoricalCols = append(categoricalCols, i)
		}
	}

	return categoricalCols
}
