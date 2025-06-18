package normalization

import "sync"

type CategoricalColumn struct {
	mapping map[string]int
	counter int
	mu      sync.Mutex
}

type CategoricalEncoder struct {
	columns map[int]*CategoricalColumn
}

func NewCategoricalEncoder(categoricalCols []int) *CategoricalEncoder {
	columns := make(map[int]*CategoricalColumn)

	for _, col := range categoricalCols {
		columns[col] = &CategoricalColumn{
			mapping: make(map[string]int),
			counter: 1,
			mu:      sync.Mutex{},
		}
	}

	return &CategoricalEncoder{
		columns: columns,
	}
}

func (ce *CategoricalEncoder) Encode(column int, value string) int {
	col, ok := ce.columns[column]
	if !ok {
		panic("column not initialized")
	}
	col.mu.Lock()
	defer col.mu.Unlock()

	if id, ok := col.mapping[value]; ok {
		return id
	}

	id := col.counter
	col.mapping[value] = id
	col.counter++

	return id
}
