package main

import (
	"FlightControl/Graph"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"gonum.org/v1/plot/plotter"
	"image/color"
)

func analysisTab() fyne.CanvasObject {
	graph1 := Graph.NewGraphWidget().
		AddTool(Graph.NewResetAxisTool()).
		AddTool(Graph.NewZoomTool()).
		AddTool(Graph.NewDragTool()).
		SetMaxBounds(0, 2, 0, 2)

	l, _ := plotter.NewLine(plotter.XYs{{0, 0}, {1, 1}, {2, 2}})
	l.Color = color.RGBA{G: 255}
	graph1.Plot.Add(l)

	grid := plotter.NewGrid()
	graph1.Plot.Add(grid)

	graph2 := Graph.NewGraphWidget().
		AddTool(Graph.NewResetAxisTool()).
		AddTool(Graph.NewZoomTool()).
		AddTool(Graph.NewDragTool()).
		SetMaxBounds(0, 2, 0, 2)

	l2, _ := plotter.NewLine(plotter.XYs{{0, 0}, {1, 1}, {2, 2}})
	l2.Color = color.RGBA{R: 255}
	graph2.Plot.Add(l2)

	r1, _ := plotter.NewScatter(plotter.XYs{{0, 0}, {1, 2}, {2, 2}})
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
