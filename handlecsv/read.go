package handlecsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func Read(
	filePath string,
	linesCh chan<- []string,
	sampleLinesCh chan<- [][]string,
	sampleSize int,
) (headers []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	defer close(linesCh)

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true

	headers, err = reader.Read()
	if err != nil {
		fmt.Println("Err read headers:", err)
	}

	sampleCounter := 0
	sample := make([][]string, 0, sampleSize)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Err read line:", err)
			continue
		}

		linesCh <- line

		if sampleCounter < sampleSize {
			sample = append(sample, line)
			sampleCounter++

			if sampleCounter == sampleSize {
				sampleLinesCh <- sample
				close(sampleLinesCh)
			}
		}
	}

	return headers, nil
}
