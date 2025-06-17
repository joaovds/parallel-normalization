package normalization

import (
	"fmt"
	"strings"
	"sync"

	"github.com/joaovds/parallel-normalization/handlecsv"
)

type FileToNormalize struct {
	SampleSize      int
	BatchSize       int
	FilePath        string
	linesCh         chan []string
	batchesCh       chan *Batch
	wg              *sync.WaitGroup
	categoricalCols []int
	headers         []string
}

func NewFileToNormalize(filePath string, sampleSize, batchSize int) *FileToNormalize {
	return &FileToNormalize{
		FilePath:        filePath,
		SampleSize:      sampleSize,
		BatchSize:       batchSize,
		wg:              &sync.WaitGroup{},
		linesCh:         make(chan []string, batchSize),
		batchesCh:       make(chan *Batch),
		categoricalCols: make([]int, 0),
		headers:         make([]string, 0),
	}
}

func (f *FileToNormalize) Normalize() error {
	sampleLinesCh := make(chan [][]string, f.SampleSize)

	f.wg.Add(1)
	go func() {
		defer f.wg.Done()

		headers, err := handlecsv.Read(f.FilePath, f.linesCh, sampleLinesCh, f.SampleSize)
		if err != nil {
			panic(err)
		}
		f.headers = headers
	}()

	f.wg.Add(1)
	go func() {
		defer f.wg.Done()
		for sample := range sampleLinesCh {
			f.categoricalCols = handlecsv.IdentifyCategoricalColumns(sample)
		}
	}()

	f.wg.Add(1)
	go f.CreateBatches()

	f.wg.Add(1)
	go func() {
		defer f.wg.Done()

		for batch := range f.batchesCh {
			fmt.Println("Processando batch com", len(batch.Data), "linhas")
		}
	}()

	f.wg.Wait()

	fmt.Println("Colunas categóricas (índices):", f.categoricalCols)
	categoricalColsNamesBuider := strings.Builder{}
	for _, colIDX := range f.categoricalCols {
		categoricalColsNamesBuider.WriteString(fmt.Sprintf(" %s |", f.headers[colIDX]))
	}
	fmt.Println("Colunas categóricas (nomes):", categoricalColsNamesBuider.String())

	return nil
}

func (f *FileToNormalize) CreateBatches() {
	defer f.wg.Done()
	defer close(f.batchesCh)

	batch := NewBatch(f.BatchSize)

	for line := range f.linesCh {
		batch.Data = append(batch.Data, line)

		if len(batch.Data) == f.BatchSize {
			f.batchesCh <- batch
			batch = NewBatch(f.BatchSize)
		}
	}

	if len(batch.Data) > 0 {
		f.batchesCh <- batch
	}
}
