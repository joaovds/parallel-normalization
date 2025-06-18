package normalization

import (
	"fmt"
	"slices"
	"strconv"
)

type Batch struct {
	Data [][]string
}

func NewBatch(batchSize int) *Batch {
	return &Batch{
		Data: make([][]string, 0, batchSize),
	}
}

func (b *Batch) Normalize(categoricalColumns []int, ce *CategoricalEncoder) {
	batchNormalized := make([][]string, 0, len(b.Data))

	for _, row := range b.Data {
		for colIdx, value := range row {
			if slices.Contains(categoricalColumns, colIdx) {
				id := ce.Encode(colIdx, value)
				row[colIdx] = strconv.Itoa(id)
			}
		}

		batchNormalized = append(batchNormalized, row)
	}

	fmt.Println("Normalizado:", len(batchNormalized))
}
