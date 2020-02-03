package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"time"

	"github.com/covrom/workerpool/cases"
	"github.com/covrom/workerpool/chart"
	"github.com/covrom/workerpool/measures"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("workerspool-test", "run the workerspool test")

	cmdRunAll        = app.Command("run", "run all the defined test, write results to JSON file")
	cmdRunAllCases   = cmdRunAll.Arg("cases", "cases to run, omit to run all the defined").Strings()
	cmdRunAllOut     = cmdRunAll.Flag("output", "output file name").Default("results.json").String()
	cmdRunAllWorkers = cmdRunAll.Flag("workers", "number of workers for workerpool types").Default("20").Int()
	cmdRunAllChanLen = cmdRunAll.Flag("chan", "channel buffer size for workerpool type").Default("200").Int()
	cmdRunAllAmount  = cmdRunAll.Flag("amount", "objects count, separate test will be running for each value listed").Short('a').Int32List()
	cmdRunAllProfile = cmdRunAll.Flag("profile", "write cpu profile file name Sprint4f template, case name and amount will be passed to it as the params").String()
	cmdRunAllChart   = cmdRunAll.Flag("chart", "file name for the PNG format chart").String()
	cmdRunAllLX      = cmdRunAll.Flag("lx", "logarithmic X axis").Bool()
	cmdRunAllLY      = cmdRunAll.Flag("ly", "logarithmic Y axis").Bool()

	cmdRunOne        = app.Command("runone", "run the specified test, write results to STDOUT")
	cmdRunOneName    = cmdRunOne.Arg("case", "name of the case to run").Required().String()
	cmdRunOneWorkers = cmdRunOne.Flag("workers", "number of workers for workerpool types").Int()
	cmdRunOneChanLen = cmdRunOne.Flag("chan", "channel buffer size for workerpool type").Int()
	cmdRunOneAmount  = cmdRunOne.Flag("amount", "objects count").Int()
	cmdRunOneProfile = cmdRunOne.Flag("profile", "write cpu profile to `file`").Default("").String()

	cmdChart    = app.Command("chart", "generate chart from the previously collected data")
	cmdChartIn  = cmdChart.Flag("in", "input file").Default("results.json").String()
	cmdChartOut = cmdChart.Flag("out", "file name for the PNG format chart").Default("results.png").String()
	cmdChartLX  = cmdChart.Flag("lx", "logarithmic X axis").Bool()
	cmdChartLY  = cmdChart.Flag("ly", "logarithmic Y axis").Bool()

	cmdCSV    = app.Command("csv", "convert results to CSV file")
	cmdCSVIn  = cmdCSV.Flag("in", "input file").Default("results.json").String()
	cmdCSVOut = cmdCSV.Flag("out", "file name for the CSV file").Default("results.csv").String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case cmdRunAll.FullCommand():
		runAll(
			*cmdRunAllCases,
			*cmdRunAllOut,
			*cmdRunAllWorkers,
			*cmdRunAllChanLen,
			*cmdRunAllAmount,
			*cmdRunAllProfile,
			*cmdRunAllChart,
			*cmdRunAllLX,
			*cmdRunAllLY,
		)
	case cmdRunOne.FullCommand():
		runOne(*cmdRunOneName, *cmdRunOneWorkers, *cmdRunOneChanLen, *cmdRunOneAmount, *cmdRunOneProfile)
	case cmdChart.FullCommand():
		runChart(*cmdChartIn, *cmdChartOut, *cmdChartLX, *cmdChartLY)
	case cmdCSV.FullCommand():
		runCSV(*cmdCSVIn, *cmdCSVOut)
	}
}

func runAll(casesList []string, out string, workers int, chanLen int, amount []int32, profile string, chartOut string, lx bool, ly bool) {
	if len(casesList) == 0 {
		casesList = enumerateCases(cases.Cases)
	}

	log.Printf("runAll: Cases: %v, Out %v, Workers %v, ChanLen %v, Amount %v, Profile %v, Chart %v", casesList, out, workers, chanLen, amount, profile, chartOut)

	res := make(map[string][]measures.Measures, len(casesList))

	for _, a := range amount {
		for _, c := range casesList {
			p := profile
			if p != "" {
				p = fmt.Sprintf(p, c, a)
			}

			cmd := exec.Command(
				os.Args[0], "runone", c,
				"--chan", strconv.Itoa(chanLen),
				"--workers", strconv.Itoa(int(workers)),
				"--amount", strconv.Itoa(int(a)),
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

			err = json.NewEncoder(os.Stdout).Encode(m)
			if err != nil {
				panic(err)
			}

			res[c] = append(res[c], m)
		}
	}

	f, err := os.Create(*cmdRunAllOut)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", " ")

	err = enc.Encode(res)
	if err != nil {
		panic(err)
	}

	if chartOut != "" {
		err = chart.Draw(res, chartOut, lx, ly)
		if err != nil {
			panic(err)
		}
	}
}

var ErrCaseUndefined = errors.New("case undefined")

func runOne(name string, workers int, chanLen int, amount int, profile string) {
	log.Printf("runOne: Name %v, Workers %v, ChanLen %v, Amount %v, Profile %v", name, workers, chanLen, amount, profile)

	if profile != "" {
		f, err := os.Create(profile)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}

		defer pprof.StopCPUProfile()
	}

	c, ok := cases.Cases[name]
	if !ok {
		panic(fmt.Errorf("%q: %w", name, ErrCaseUndefined))
	}

	chin, chout := c.Prepare(workers, chanLen, amount)

	m1, m2 := &runtime.MemStats{}, &runtime.MemStats{}

	start := time.Now()
	runtime.ReadMemStats(m1)

	c.Run(workers, chanLen, amount, chin, chout)

	spent := time.Now().Sub(start)

	runtime.ReadMemStats(m2)

	err := json.NewEncoder(os.Stdout).Encode(
		measures.Measures{
			Case:            name,
			Workers:         workers,
			ChanLen:         chanLen,
			Amount:          amount,
			Spent:           spent,
			AllocBytes:      m2.Alloc - m1.Alloc,
			AllocObjects:    m2.Mallocs - m1.Mallocs,
			AllocBytesTotal: m2.TotalAlloc - m1.TotalAlloc,
			SystemMem:       m2.Sys - m1.Sys,
		},
	)
	if err != nil {
		panic(err)
	}
}

func runChart(in string, out string, lx bool, ly bool) {
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

func runCSV(inName string, outName string) {
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

func enumerateCases(c map[string]cases.Case) []string {
	l := make([]string, 0, len(c))

	for n := range c {
		l = append(l, n)
	}

	sort.Strings(l)

	return l
}
