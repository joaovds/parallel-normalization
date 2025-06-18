package handlecsv

import (
	"encoding/csv"
	"os"
	"strconv"
)

func Write(filePath string, values map[string]int) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write([]string{"id", "value"})
	for value, id := range values {
		writer.Write([]string{
			strconv.Itoa(id),
			value,
		})
	}

	writer.Flush()

	return nil
}

func ResetOutputDir(outputDir string) error {
	if err := os.RemoveAll(outputDir); err != nil {
		return err
	}
	return os.MkdirAll(outputDir, os.ModePerm)
}
