package helper

import (
	"encoding/csv"
	"log"
	"os"
)

var AvailCmds [][]string

func ParseCommands() [][]string {
	f, err := os.Open("./config/commands.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	AvailCmds, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return AvailCmds[1:]
}
