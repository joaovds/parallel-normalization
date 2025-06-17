package handlecsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func Read(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = 59
	reader.LazyQuotes = true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Err read line:", err)
			continue
		}
		fmt.Println(record)
	}

	return nil
}
