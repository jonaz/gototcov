package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"golang.org/x/tools/cover"
)

var filename = flag.String("f", "", "Filename of the cover profile")
var errLevel = flag.Float64("limit", 0, "% threshold to throw error")
var ignoreZero = flag.Bool("ignore-zero", false, "Ignore files with 0%. Example main.go")

func main() {
	flag.Parse()
	profiles, err := cover.ParseProfiles(*filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)

	files := 0.0
	covsum := 0.0
	for _, v := range profiles {
		coverage := percentCovered(v)
		if *ignoreZero && coverage == 0.0 {
			continue
		}

		files++
		covsum = covsum + coverage
		fmt.Fprintf(w, v.FileName+":\t%.2f\n", coverage)
	}

	totalCoverage := covsum / files

	fmt.Fprintf(w, "Total coverage:\t%.2f\n", totalCoverage)
	w.Flush()

	if *errLevel != 0 && totalCoverage < *errLevel {
		fmt.Println("Error: Expected coverage to be over " + strconv.FormatFloat(*errLevel, 'f', 2, 64))
		os.Exit(1)
	}
}

func percentCovered(p *cover.Profile) float64 {
	var total, covered int64
	for _, b := range p.Blocks {
		total += int64(b.NumStmt)
		if b.Count > 0 {
			covered += int64(b.NumStmt)
		}
	}
	if total == 0 {
		return 0
	}
	return float64(covered) / float64(total) * 100
}
