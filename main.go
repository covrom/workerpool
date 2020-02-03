package main

import (
	"os"

	"github.com/covrom/workerpool/commands"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("workerspool-test", "run the workerspool test")

	cmdRunAll        = app.Command("run", "run all the defined test, write results to JSON file")
	cmdRunAllCases   = cmdRunAll.Arg("cases", "cases to run, omit to run all the defined").Strings()
	cmdRunAllOut     = cmdRunAll.Flag("res", "file the results are stored").Default("results.json").String()
	cmdRunAllWorkers = cmdRunAll.Flag("workers", "number of workers for workerpool types").Default("20").Int()
	cmdRunAllChanLen = cmdRunAll.Flag("chan", "channel buffer size for workerpool type").Default("10000").Int()
	cmdRunAllAmount  = cmdRunAll.Flag("amount", "objects count, separate test will be running for each value listed").Short('a').Int32List()
	cmdRunAllProfile = cmdRunAll.Flag("profile", "write cpu profile to `file`").String()
	cmdRunAllChart   = cmdRunAll.Flag("chart", "file name for the PNG format chart").String()
	cmdRunAllLX      = cmdRunAll.Flag("lx", "logarithmic X axis").Bool()
	cmdRunAllLY      = cmdRunAll.Flag("ly", "logarithmic Y axis").Bool()

	cmdRunOne        = app.Command("runone", "run the specified test, write results to STDOUT")
	cmdRunOneName    = cmdRunOne.Arg("case", "name of the case to run").Required().String()
	cmdRunOneWorkers = cmdRunOne.Flag("workers", "number of workers for workerpool types").Required().Int()
	cmdRunOneChanLen = cmdRunOne.Flag("chan", "channel buffer size for workerpool type").Required().Int()
	cmdRunOneAmount  = cmdRunOne.Flag("amount", "objects count").Required().Int()
	cmdRunOneProfile = cmdRunOne.Flag("profile", "write cpu profile to `file`").Default("").String()

	cmdChart    = app.Command("chart", "generate chart from the previously collected data")
	cmdChartIn  = cmdChart.Flag("res", "file the results are stored").Default("results.json").String()
	cmdChartOut = cmdChart.Flag("chart", "file name for the PNG format chart").Default("results.png").String()
	cmdChartLX  = cmdChart.Flag("lx", "logarithmic X axis").Bool()
	cmdChartLY  = cmdChart.Flag("ly", "logarithmic Y axis").Bool()

	cmdCSV    = app.Command("csv", "convert results to CSV file")
	cmdCSVIn  = cmdCSV.Flag("res", "file the results are stored").Default("results.json").String()
	cmdCSVOut = cmdCSV.Flag("out", "file name for the CSV file").Default("results.csv").String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case cmdRunAll.FullCommand():
		commands.RunAll(
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
		commands.RunOne(*cmdRunOneName, *cmdRunOneWorkers, *cmdRunOneChanLen, *cmdRunOneAmount, *cmdRunOneProfile)
	case cmdChart.FullCommand():
		commands.Chart(*cmdChartIn, *cmdChartOut, *cmdChartLX, *cmdChartLY)
	case cmdCSV.FullCommand():
		commands.CSV(*cmdCSVIn, *cmdCSVOut)
	}
}
