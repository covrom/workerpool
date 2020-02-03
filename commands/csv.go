package commands

import (
	"encoding/csv"
	"encoding/json"
	"os"

	"github.com/covrom/workerpool/measures"
)

func CSV(inName string, outName string) {
	in, err := os.Open(inName)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	var m map[string][]measures.Measures

	err = json.NewDecoder(in).Decode(&m)
	if err != nil {
		panic(err)
	}

	out, err := os.Create(outName)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	w := csv.NewWriter(out)

	err = w.Write(measures.Fields())
	if err != nil {
		panic(err)
	}

	for _, group := range m {
		for _, line := range group {
			err = w.Write(line.Fields())
			if err != nil {
				panic(err)
			}
		}
	}

	w.Flush()
}
