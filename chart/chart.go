package chart

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/covrom/workerpool/measures"
	"github.com/wcharczuk/go-chart"
)

func Draw(m map[string][]measures.Measures, outName string) error {
	for _, n := range measures.Values() {
		err := drawOne(n, m, outName)
		if err != nil {
			return fmt.Errorf("%q: %w", n, err)
		}
	}

	return nil
}

func drawOne(field string, m map[string][]measures.Measures, outName string) error {
	log.Printf("creating %q chart", field)
	graph := chart.Chart{
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
		},
		YAxis: chart.YAxis{
			Name:  field,
			Style: chart.Style{Show: true},
		},
		XAxis: chart.XAxis{
			Name:  "Amount",
			Style: chart.Style{Show: true, TextRotationDegrees: 45.0},
		},
		Series: make([]chart.Series, 0, len(m)),
	}

	cases := getCasesSorted(m)

	for _, c := range cases {
		s := chart.ContinuousSeries{
			Name:    c,
			XValues: make([]float64, 0, len(m[c])),
			YValues: make([]float64, 0, len(m[c])),
			Style:   chart.Style{Show: true, DotWidth: 3},
		}

		for _, p := range m[c] {
			s.XValues = append(s.XValues, float64(p.Amount))
			s.YValues = append(s.YValues, p.Value(field))
		}

		graph.Series = append(graph.Series, s)
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
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
