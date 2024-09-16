package main

import (
	"FlightControl/Graph"
	"fyne.io/fyne/v2"
	"gonum.org/v1/plot/plotter"
)

func analysisTab() fyne.CanvasObject {
	graph := Graph.NewGraphWidget()
	graph.SetMaxBounds(0, 2, 0, 2)

	l, err := plotter.NewLine(plotter.XYs{{0, 0}, {1, 1}, {2, 2}})
	if err != nil {
		panic(err)
	}
	l.Color = plotter.DefaultLineStyle.Color
	graph.Plot.Add(l)

	return graph
}
