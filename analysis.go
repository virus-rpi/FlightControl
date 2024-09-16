package main

import (
	"FlightControl/Graph"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"gonum.org/v1/plot/plotter"
	"image/color"
)

func analysisTab() fyne.CanvasObject {
	graph1 := Graph.NewGraphWidget()
	graph1.SetMaxBounds(0, 2, 0, 2)

	l, err := plotter.NewLine(plotter.XYs{{0, 0}, {1, 1}, {2, 2}})
	if err != nil {
		panic(err)
	}
	l.Color = color.RGBA{G: 255}
	graph1.Plot.Add(l)

	grid := plotter.NewGrid()
	graph1.Plot.Add(grid)

	graph2 := Graph.NewGraphWidget()
	graph2.SetMaxBounds(0, 2, 0, 2)

	l2, err := plotter.NewLine(plotter.XYs{{0, 0}, {1, 1}, {2, 2}})
	if err != nil {
		panic(err)
	}
	l2.Color = color.RGBA{R: 255}
	graph2.Plot.Add(l2)

	r1, err := plotter.NewScatter(plotter.XYs{{0, 0}, {1, 2}, {2, 2}})
	if err != nil {
		panic(err)
	}
	r1.GlyphStyle.Color = color.RGBA{G: 255}
	graph2.Plot.Add(r1)

	grid2 := plotter.NewGrid()
	graph2.Plot.Add(grid2)

	graph1.SetMinWidgetSize(fyne.NewSize(10, 300))
	graph2.SetMinWidgetSize(fyne.NewSize(10, 300))

	return container.NewVScroll(
		container.NewVBox(
			graph1,
			graph2,
		),
	)
}
