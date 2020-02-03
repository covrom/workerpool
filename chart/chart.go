package chart

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/covrom/workerpool/measures"
	"github.com/wcharczuk/go-chart"
)

func Draw(m map[string][]measures.Measures, outName string, lx bool, ly bool) error {
	for _, n := range measures.Values() {
		err := drawOne(n, m, outName, lx, ly)
		if err != nil {
			return fmt.Errorf("%q: %w", n, err)
		}
	}

	return nil
}

func drawOne(field string, m map[string][]measures.Measures, outName string, lx bool, ly bool) error {
	log.Printf("creating %q chart", field)
	graph := chart.Chart{
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
		},
		YAxis: chart.YAxis{
			Name: field,
		},
		XAxis: chart.XAxis{
			Name:  "Amount",
			Style: chart.Style{TextRotationDegrees: 45.0},
		},
		Series: make([]chart.Series, 0, len(m)),
	}

	cases := getCasesSorted(m)

	for _, c := range cases {
		s := chart.ContinuousSeries{
			Name:    c,
			XValues: make([]float64, 0, len(m[c])),
			YValues: make([]float64, 0, len(m[c])),
			Style:   chart.Style{DotWidth: 3},
		}

		for _, p := range m[c] {
			s.XValues = append(s.XValues, float64(p.Amount))
			s.YValues = append(s.YValues, p.Value(field))
			graph.XAxis.Ticks = append(graph.XAxis.Ticks, chart.Tick{Value: float64(p.Amount), Label: strconv.Itoa(p.Amount)})
			//graph.YAxis.Ticks = append(graph.YAxis.Ticks, chart.Tick{Value: float64(int(p.Value(field))), Label: fmt.Sprintf("%d", int(p.Value(field)))})
		}

		graph.Series = append(graph.Series, s)
	}

	graph.XAxis.Ticks = cleanTicks(graph.XAxis.Ticks)
	//graph.YAxis.Ticks = cleanTicks(graph.YAxis.Ticks)

	if lx {
		graph.XAxis.Range = &chart.LogarithmicRange{}
	}

	if ly {
		graph.YAxis.Range = &chart.LogarithmicRange{}
	}

	graph.Elements = []chart.Renderable{
		chart.LegendThin(&graph),
	}

	f, err := os.Create(field + "." + outName)
	if err != nil {
		return err
	}
	defer f.Close()

	err = graph.Render(chart.PNG, f)
	if err != nil {
		return err
	}

	log.Printf("created %q chart", field)
	return nil
}

func getCasesSorted(m map[string][]measures.Measures) []string {
	l := make([]string, 0, len(m))

	for c := range m {
		l = append(l, c)
	}

	sort.Strings(l)

	return l
}

func cleanTicks(in []chart.Tick) []chart.Tick {
	sort.Slice(in, func(i, j int) bool { return in[i].Value < in[j].Value })

	j := 0
	for i := 1; i < len(in); i++ {
		if in[j].Value == in[i].Value {
			continue
		}

		j++

		in[j] = in[i]
	}

	return in[:j+1]
}
