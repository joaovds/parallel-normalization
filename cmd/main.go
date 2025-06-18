package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joaovds/parallel-normalization/normalization"
)

func main() {
	start := time.Now()

	fileToNormalize := normalization.NewFileToNormalize("./extra_large_dataset.csv", "./output", 8, 200000)
	err := fileToNormalize.Normalize()
	if err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(start)
	fmt.Println("Dataset normalized! :)")
	fmt.Printf("Time: %s\n", elapsed)
}
