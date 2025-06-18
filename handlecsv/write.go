package handlecsv

import (
	"encoding/csv"
	"os"
	"strconv"
	"sync"
)

type NormalizedWriter struct {
	file   *os.File
	writer *csv.Writer
	mu     sync.Mutex
}

func NewNormalizedWriter(filePath string) (*NormalizedWriter, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	return &NormalizedWriter{
		file:   file,
		writer: csv.NewWriter(file),
		mu:     sync.Mutex{},
	}, nil
}

func (n *NormalizedWriter) Write(rows [][]string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	for _, row := range rows {
		if err := n.writer.Write(row); err != nil {
			return err
		}
	}
	n.writer.Flush()
	return n.writer.Error()
}

func (n *NormalizedWriter) Close() error {
	n.writer.Flush()
	if err := n.writer.Error(); err != nil {
		return err
	}
	return n.file.Close()
}

func WriteCategoricalFiles(filePath string, values map[string]int) error {
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
