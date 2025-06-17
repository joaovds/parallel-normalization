package main

import (
	"log"

	"github.com/joaovds/parallel-normalization/handlecsv"
)

func main() {
	err := handlecsv.Read("./small_dataset.csv")
	if err != nil {
		log.Fatal(err)
	}
}
