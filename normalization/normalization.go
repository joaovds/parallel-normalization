package normalization

import (
	"fmt"
	"sync"

	"github.com/joaovds/parallel-normalization/handlecsv"
)

type FileToNormalize struct {
	SampleSize         int
	BatchSize          int
	FilePath           string
	outputDir          string
	linesCh            chan []string
	batchesCh          chan *Batch
	wg                 *sync.WaitGroup
	categoricalCols    []int
	headers            []string
	categoricalEncoder *CategoricalEncoder
}

func NewFileToNormalize(filePath, outputDir string, sampleSize, batchSize int) *FileToNormalize {
	return &FileToNormalize{
		FilePath:        filePath,
		outputDir:       outputDir,
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
	if err := handlecsv.ResetOutputDir(f.outputDir); err != nil {
		panic(fmt.Sprintf("Err reset output dir: %v", err))
	}
	sampleLinesCh := make(chan [][]string, f.SampleSize)

	f.wg.Add(1)
	go f.ReadCSV(sampleLinesCh)

	f.wg.Add(1)
	go f.FindCategoricalCols(sampleLinesCh)

	f.wg.Add(1)
	go f.CreateBatches()

	f.wg.Add(1)
	f.HandleBatches()

	f.wg.Wait()

	f.WriteCategoricalsCSV()

	return nil
}

func (f *FileToNormalize) ReadCSV(sampleLinesCh chan<- [][]string) {
	defer f.wg.Done()
	headers, err := handlecsv.Read(f.FilePath, f.linesCh, sampleLinesCh, f.SampleSize)
	if err != nil {
		panic(err)
	}
	f.headers = headers
}

func (f *FileToNormalize) FindCategoricalCols(sampleLinesCh <-chan [][]string) {
	defer f.wg.Done()
	for sample := range sampleLinesCh {
		f.categoricalCols = handlecsv.IdentifyCategoricalColumns(sample)
		f.categoricalEncoder = NewCategoricalEncoder(f.categoricalCols)
	}
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

func (f *FileToNormalize) HandleBatches() {
	defer f.wg.Done()
	batchesWg := sync.WaitGroup{}
	for batch := range f.batchesCh {
		batchesWg.Add(1)
		go func() {
			batch.Normalize(f.categoricalCols, f.categoricalEncoder)
			batchesWg.Done()
		}()
	}
	batchesWg.Wait()
}

func (f *FileToNormalize) WriteCategoricalsCSV() {
	for i, column := range f.categoricalEncoder.columns {
		columnName := f.headers[i]

		filePath := fmt.Sprintf("%s/%s_mapping.csv", f.outputDir, columnName)
		handlecsv.Write(filePath, column.mapping)
	}
}
