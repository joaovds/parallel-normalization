package handlecsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func Read(filePath string, sampleSize int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true

	var sample [][]string
	lineCount := 0

	headers, err := reader.Read()
	if err != nil {
		fmt.Println("Err read headers:", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Err read line:", err)
			continue
		}
		if lineCount < sampleSize {
			sample = append(sample, record)
			lineCount++
		}
		fmt.Println(record, len(record))
	}

	categoricalCols := IdentifyCategoricalColumns(sample)
	fmt.Println("Colunas categóricas (índices):", categoricalCols)
	categoricalColsNamesBuider := strings.Builder{}
	for _, colIDX := range categoricalCols {
		categoricalColsNamesBuider.WriteString(fmt.Sprintf(" %s |", headers[colIDX]))
	}
	fmt.Println("Colunas categóricas (nomes):", categoricalColsNamesBuider.String())

	return nil
}
