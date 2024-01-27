package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/moskyb/upbank-fbar-calculator/fbar"
)

func main() {
	tok := os.Getenv("UP_TOKEN")
	if tok == "" {
		panic("UP_TOKEN environment variable not set")
	}

	year := os.Getenv("YEAR")
	if year == "" {
		panic("YEAR environment variable not set")
	}

	intYear, err := strconv.Atoi(year)
	if err != nil {
		panic(err)
	}

	r, err := fbar.GenerateReport(tok, intYear)
	if err != nil {
		panic(err)
	}

	fmt.Println(r.PrettyString())
}
