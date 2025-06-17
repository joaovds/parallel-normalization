package normalization

type Batch struct {
	Data [][]string
}

func NewBatch(batchSize int) *Batch {
	return &Batch{
		Data: make([][]string, 0, batchSize),
	}
}
