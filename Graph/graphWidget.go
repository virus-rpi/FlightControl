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
	image            *canvas.Image
	Plot             *plot.Plot
	XMin             float64
	XMax             float64
	YMin             float64
	YMax             float64
	scale            float64
	currentSelection [4]float64
	resetButton      *widget.Button
}

func NewGraphWidget() *Widget {
	w := &Widget{}
	w.ExtendBaseWidget(w)
	w.Plot = plot.New()
	w.resetButton = widget.NewButton("Reset axis", w.resetAxis)
	w.image = canvas.NewImageFromImage(w.render())
	w.resetAxis()
	return w
}

func (w *Widget) SetMaxBounds(xMin, xMax, yMin, yMax float64) {
	w.XMin = xMin
	w.XMax = xMax
	w.YMin = yMin
	w.YMax = yMax
	w.resetAxis()
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

		c.SetColor(color.NRGBA{R: 173, G: 216, B: 230, A: 204})
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
	w.Plot.X.Min = w.XMin
	w.Plot.X.Max = w.XMax
	w.Plot.Y.Min = w.YMin
	w.Plot.Y.Max = w.YMax
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
	gr.graphWidget.resetButton.Resize(fyne.NewSize(100, 40))
	gr.image.Resize(fyne.NewSize(size.Width, size.Height-gr.graphWidget.resetButton.Size().Height))
	gr.graphWidget.resetButton.Move(fyne.NewPos(size.Width-gr.graphWidget.resetButton.Size().Width, size.Height-gr.graphWidget.resetButton.Size().Height))
	gr.Refresh()
}

func (gr *graphRenderer) MinSize() fyne.Size { return fyne.NewSize(100, 100) }

func (gr *graphRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{gr.image, gr.graphWidget.resetButton}
}

func (gr *graphRenderer) Refresh() {
	gr.graphWidget.Refresh()
	gr.image.Refresh()
}
