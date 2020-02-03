package commands

import (
	"encoding/json"
	"os"

	"github.com/covrom/workerpool/chart"
	"github.com/covrom/workerpool/measures"
)

func Chart(in string, out string, lx bool, ly bool) {
	f, err := os.Open(in)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var m map[string][]measures.Measures

	err = json.NewDecoder(f).Decode(&m)
	if err != nil {
		panic(err)
	}

	err = chart.Draw(m, out, lx, ly)
	if err != nil {
		panic(err)
	}
}
