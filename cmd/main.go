package main

import (
	"log"

	"github.com/joaovds/parallel-normalization/normalization"
)

func main() {
	fileToNormalize := normalization.NewFileToNormalize("./small_dataset.csv", "./output", 10, 100)
	err := fileToNormalize.Normalize()
	if err != nil {
		log.Fatal(err)
	}
}
