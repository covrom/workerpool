package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/covrom/workerpool/cases"
	"github.com/covrom/workerpool/chart"
	"github.com/covrom/workerpool/measures"
)

func RunAll(
	casesList []string,
	out string,
	workers int,
	chanLen int,
	amount []string,
	profile string,
	chartOut string,
	lx bool,
	ly bool,
) {
	if len(casesList) == 0 {
		casesList = enumerateCases(cases.Cases)
	}

	log.Printf("runAll: Cases: %v, Out %v, Workers %v, ChanLen %v, Amount %v, Profile %v, Chart %v", casesList, out, workers, chanLen, amount, profile, chartOut)

	res := make(map[string][]measures.Measures, len(casesList))

	for _, a0 := range amount {
		for _, a := range strings.Split(a0, ",") {
			if _, err := strconv.Atoi(a); err != nil {
				log.Fatalf("incorrect -a argument: %s", a)
			}
			for _, c := range casesList {
				p := profile
				if p != "" {
					p = fmt.Sprintf("%s.%d.%s", c, a, p)
				}

				cmd := exec.Command(
					os.Args[0], "runone", c,
					"--chan", strconv.Itoa(chanLen),
					"--workers", strconv.Itoa(int(workers)),
					"--amount", a,
					"--profile", p,
				)

				b, err := cmd.Output()
				if err != nil {
					log.Printf("Error running %q: %v", cmd.String(), err)
					continue
				}

				m := measures.Measures{}
				err = json.Unmarshal(b, &m)
				if err != nil {
					panic(err)
				}

				log.Printf("%#v", m)

				res[c] = append(res[c], m)

				err = saveResults(res, out)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	if chartOut != "" {
		err := chart.Draw(res, chartOut, lx, ly)
		if err != nil {
			panic(err)
		}
	}
}

func enumerateCases(c map[string]cases.Case) []string {
	l := make([]string, 0, len(c))

	for n := range c {
		l = append(l, n)
	}

	sort.Strings(l)

	return l
}

func saveResults(res map[string][]measures.Measures, out string) error {
	err := writeResults(res, out+".tmp")
	if err != nil {
		return err
	}

	err = os.Rename(out+".tmp", out)
	if err != nil {
		return err
	}

	return nil
}

func writeResults(res map[string][]measures.Measures, out string) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}

	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", " ")

	err = enc.Encode(res)
	if err != nil {
		return err
	}

	return nil
}
