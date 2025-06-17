package handlecsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func Read(filePath string, linesCh chan<- []string) (headers []string, err error) {
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
	reader.Read()

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
	}

	return headers, nil
}
