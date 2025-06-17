package main

import (
	"log"

	"github.com/joaovds/parallel-normalization/normalization"
)

func main() {
	fileToNormalize := normalization.NewFileToNormalize("./small_dataset.csv", 10, 100)
	err := fileToNormalize.Normalize()
	if err != nil {
		log.Fatal(err)
	}
}
