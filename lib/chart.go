package lib

import (
	"fmt"
	"os"
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// GenerateLOCChart creates a line chart showing lines of code evolution across commits
func GenerateLOCChart(commits []CommitData, languageSet map[string]bool, outputFile string) error {
	// Convert language set to sorted slice
	var languages []string
	for lang := range languageSet {
		languages = append(languages, lang)
	}
	sort.Strings(languages)

	// Prepare data for chart
	var xAxis []string
	seriesData := make(map[string][]opts.LineData)

	// Initialize series data
	for _, lang := range languages {
		seriesData[lang] = make([]opts.LineData, 0, len(commits))
	}

	// Populate data
	for _, commit := range commits {
		// Use short hash for x-axis
		shortHash := commit.Hash[:8]
		xAxis = append(xAxis, shortHash)

		// Add data for each language
		for _, lang := range languages {
			count := commit.Languages[lang]
			seriesData[lang] = append(seriesData[lang], opts.LineData{
				Value: count,
			})
		}
	}

	// Create line chart
	chart := charts.NewLine()
	chart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Lines of Code Evolution",
			Subtitle: "Language breakdown across commits",
			Top:      "5%",
			Left:     "5%",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Commits",
			Type: "category",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Lines of Code",
			Type: "value",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    true,
			Trigger: "axis",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:  true,
			Top:   "5%",
			Right: "5%",
			Type:  "scroll",
		}),
		charts.WithGridOpts(opts.Grid{
			Top:    "15%",
			Bottom: "15%",
			Left:   "10%",
			Right:  "20%",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "800px",
		}),
	)

	// Add series for each language
	for _, lang := range languages {
		chart.AddSeries(lang, seriesData[lang],
			charts.WithLineChartOpts(opts.LineChart{
				Smooth: true,
			}),
		)
	}

	// Set x-axis data
	chart.SetXAxis(xAxis)

	// Save chart to file
	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create chart file: %w", err)
	}
	defer f.Close()

	return chart.Render(f)
}
