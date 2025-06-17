package normalization

import (
	"fmt"
	"sync"
	"time"

	"github.com/joaovds/parallel-normalization/handlecsv"
)

type FileToNormalize struct {
	SampleSize int
	BatchSize  int
	FilePath   string
	linesCh    chan []string
	batchesCh  chan *Batch
	wg         *sync.WaitGroup
	headers    []string
}

func NewFileToNormalize(filePath string, sampleSize, batchSize int) *FileToNormalize {
	return &FileToNormalize{
		FilePath:   filePath,
		SampleSize: sampleSize,
		BatchSize:  batchSize,
		wg:         &sync.WaitGroup{},
		linesCh:    make(chan []string, batchSize),
		batchesCh:  make(chan *Batch),
	}
}

func (f *FileToNormalize) Normalize() error {
	defer f.wg.Wait()

	f.wg.Add(1)
	go func() {
		defer f.wg.Done()

		headers, err := handlecsv.Read(f.FilePath, f.linesCh)
		if err != nil {
			panic(err)
		}
		f.headers = headers
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

	// categoricalCols := IdentifyCategoricalColumns(sample)
	// fmt.Println("Colunas categóricas (índices):", categoricalCols)
	// categoricalColsNamesBuider := strings.Builder{}
	// for _, colIDX := range categoricalCols {
	// 	categoricalColsNamesBuider.WriteString(fmt.Sprintf(" %s |", headers[colIDX]))
	// }
	// fmt.Println("Colunas categóricas (nomes):", categoricalColsNamesBuider.String())

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
			fmt.Println("Criando batch")
			time.Sleep(time.Second * 2)
		}
	}

	if len(batch.Data) > 0 {
		fmt.Println("Criando batch restante")
		f.batchesCh <- batch
	}
}

// func (f *FileToNormalize) HandleBatch(batch [][]string) {
// 	f.wg.Add(1)
// 	fmt.Println("------")
// 	// fmt.Println(batchCopy)
// 	f.wg.Done()
// }
