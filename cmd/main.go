package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joaovds/parallel-normalization/normalization"
)

func main() {
	inputFile := flag.String("i", "", "Path to input CSV file")
	outputDir := flag.String("o", "./output", "Output directory")
	sampleSize := flag.Int("samples", 2, "Number of samples to process (> 0)")
	batchSize := flag.Int("batch", 4000, "Number of lines per batch (> 0")

	flag.Parse()
	validateParams(*inputFile, *outputDir, *sampleSize, *batchSize)

	start := time.Now()

	fileToNormalize := normalization.NewFileToNormalize("./small_dataset.csv", "./output", 8, 200000)
	err := fileToNormalize.Normalize()
	if err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(start)
	fmt.Println("Dataset normalized! :)")
	fmt.Printf("Time: %s\n", elapsed)
}

func validateParams(inputFile, outputDir string, sampleSize, batchSize int) {
	if inputFile == "" {
		fmt.Println("Input file is required. Use -i <file>")
		flag.Usage()
		os.Exit(1)
	}

	if sampleSize == 0 {
		fmt.Println("Invalid sample size. Use sample size > 0")
		flag.Usage()
		os.Exit(1)
	}

	if batchSize == 0 {
		fmt.Println("Invalid batch size. Use batch size > 0")
		flag.Usage()
		os.Exit(1)
	}
}
