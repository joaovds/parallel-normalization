package main

import (
	"log"

	"github.com/joaovds/parallel-normalization/handlecsv"
)

func main() {
	err := handlecsv.Read("./small_dataset.csv", 10)
	if err != nil {
		log.Fatal(err)
	}
}
