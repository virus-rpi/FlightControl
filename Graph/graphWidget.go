package Graph

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	"image"
	"image/color"
)

type Widget struct {
	widget.BaseWidget
	MinWidgetSize    fyne.Size
	image            *canvas.Image
	Plot             *plot.Plot
	PlotXMin         float64
	PlotXMax         float64
	PlotYMin         float64
	PlotYMax         float64
	scale            float64
	currentSelection [4]float64
	resetButton      *widget.Button
}

func NewGraphWidget() *Widget {
	w := &Widget{}
	w.ExtendBaseWidget(w)
	w.Plot = plot.New()

	w.Plot.BackgroundColor = color.RGBA{R: 20, G: 20, B: 20}
	w.Plot.X.Color = color.White
	w.Plot.Y.Color = color.White
	w.Plot.X.Label.TextStyle.Color = color.White
	w.Plot.Y.Label.TextStyle.Color = color.White
	w.Plot.X.Tick.LineStyle.Color = color.White
	w.Plot.Y.Tick.LineStyle.Color = color.White
	w.Plot.Y.Tick.Label.Color = color.White
	w.Plot.X.Tick.Label.Color = color.White

	w.resetButton = widget.NewButton("Reset axis", w.resetAxis)

	w.image = canvas.NewImageFromImage(w.render())
	w.resetAxis()
	return w
}

func (w *Widget) SetMaxBounds(xMin, xMax, yMin, yMax float64) {
	w.PlotXMin = xMin
	w.PlotXMax = xMax
	w.PlotYMin = yMin
	w.PlotYMax = yMax
	w.resetAxis()
}

func (w *Widget) SetMinWidgetSize(size fyne.Size) {
	w.MinWidgetSize = size
}

func (w *Widget) render() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(w.Size().Width), int(w.Size().Height-w.resetButton.Size().Height)))
	c := vgimg.NewWith(vgimg.UseImage(img))
	dc := draw.New(c)
	w.Plot.Draw(dc)

	if w.currentSelection[0] != 0 || w.currentSelection[1] != 0 || w.currentSelection[2] != 0 || w.currentSelection[3] != 0 {
		rect := vg.Rectangle{
			Min: vg.Point{X: vg.Length(w.pixelToDotsX(w.currentSelection[0])), Y: vg.Length(w.pixelToDotsY(w.currentSelection[1]))},
			Max: vg.Point{X: vg.Length(w.pixelToDotsX(w.currentSelection[2])), Y: vg.Length(w.pixelToDotsY(w.currentSelection[3]))},
		}

		c.SetColor(color.NRGBA{R: 173, G: 216, B: 230, A: 150})
		c.Fill(rect.Path())
	}

	return c.Image()
}

func (w *Widget) pixelToDotsX(x float64) float64 {
	scale := float64(vg.Inch / vgimg.DefaultDPI)
	return x * scale
}
func (w *Widget) pixelToDotsY(y float64) float64 {
	scale := float64(vg.Inch / vgimg.DefaultDPI)
	return (float64(w.Size().Height) - y - float64(w.resetButton.Size().Height)) * scale
}

func (w *Widget) resetAxis() {
	w.Plot.X.Min = w.PlotXMin
	w.Plot.X.Max = w.PlotXMax
	w.Plot.Y.Min = w.PlotYMin
	w.Plot.Y.Max = w.PlotYMax
	w.Refresh()
}

func (w *Widget) Dragged(ev *fyne.DragEvent) {
	if w.currentSelection[0] == 0 && w.currentSelection[1] == 0 {
		w.currentSelection[0] = float64(ev.Position.X)
		w.currentSelection[1] = float64(ev.Position.Y)
	}
	w.currentSelection[2] = float64(ev.Position.X)
	w.currentSelection[3] = float64(ev.Position.Y)
	w.Refresh()
}

func (w *Widget) DragEnd() {
	widthPlot := w.Plot.X.Max - w.Plot.X.Min
	xScale := widthPlot / float64(w.Size().Width)
	newXMin := w.currentSelection[0] * xScale
	newXMax := w.currentSelection[2] * xScale
	heightPlot := w.Plot.Y.Max - w.Plot.Y.Min
	yScale := heightPlot / float64(w.Size().Height-w.resetButton.Size().Height)
	newYMin := heightPlot - w.currentSelection[1]*yScale
	newYMax := heightPlot - w.currentSelection[3]*yScale
	w.Plot.X.Min = newXMin
	w.Plot.X.Max = newXMax
	w.Plot.Y.Min = newYMin
	w.Plot.Y.Max = newYMax
	w.currentSelection[0] = 0
	w.currentSelection[1] = 0
	w.currentSelection[2] = 0
	w.currentSelection[3] = 0
	w.Refresh()
}

func (w *Widget) Refresh() {
	w.image.Image = w.render()
	w.image.Refresh()
}

func (w *Widget) CreateRenderer() fyne.WidgetRenderer {
	return &graphRenderer{image: w.image, graphWidget: w}
}

type graphRenderer struct {
	image       *canvas.Image
	graphWidget *Widget
}

func (gr *graphRenderer) Destroy() {}

func (gr *graphRenderer) Layout(size fyne.Size) {
	gr.graphWidget.Resize(size)
	if fyne.CurrentDevice().IsMobile() {
		gr.graphWidget.resetButton.Resize(fyne.NewSize(130, 40))
	} else {
		gr.graphWidget.resetButton.Resize(fyne.NewSize(70, 30))
	}
	gr.image.Resize(fyne.NewSize(size.Width, size.Height-gr.graphWidget.resetButton.Size().Height))
	gr.graphWidget.resetButton.Move(fyne.NewPos(size.Width-gr.graphWidget.resetButton.Size().Width-5, size.Height-gr.graphWidget.resetButton.Size().Height))
	gr.Refresh()
}

func (gr *graphRenderer) MinSize() fyne.Size {
	if gr.graphWidget.MinWidgetSize != (fyne.Size{}) {
		return gr.graphWidget.MinWidgetSize
	}
	return fyne.NewSize(100, 150)
}

func (gr *graphRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{gr.image, gr.graphWidget.resetButton}
}

func (gr *graphRenderer) Refresh() {
	gr.graphWidget.Refresh()
	gr.image.Refresh()
}
