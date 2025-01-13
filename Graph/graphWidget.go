package Graph

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
	MinWidgetSize   fyne.Size
	image           *canvas.Image
	Plot            *plot.Plot
	PlotXMin        float64
	PlotXMax        float64
	PlotYMin        float64
	PlotYMax        float64
	scale           float64
	tools           []tool
	buttons         []fyne.CanvasObject
	buttonContainer *fyne.Container
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

	w.buttonContainer = container.NewHBox(w.buttons...)

	w.image = canvas.NewImageFromImage(w.render())
	w.resetAxis()
	return w
}

func (w *Widget) AddTool(tool tool) *Widget {
	tool.setWidget(w)
	if tool.hasIntent("button") {
		tool.registerButtons()
		w.buttonContainer.Objects = w.buttons
	}
	tool.Enable()
	w.tools = append(w.tools, tool)
	return w
}

func (w *Widget) SetMaxBounds(xMin, xMax, yMin, yMax float64) *Widget {
	w.PlotXMin = xMin
	w.PlotXMax = xMax
	w.PlotYMin = yMin
	w.PlotYMax = yMax
	w.resetAxis()

	return w
}

func (w *Widget) SetMinWidgetSize(size fyne.Size) {
	w.MinWidgetSize = size
}

func (w *Widget) render() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(w.Size().Width), int(w.Size().Height-w.buttonContainer.Size().Height)))
	c := vgimg.NewWith(vgimg.UseImage(img))
	dc := draw.New(c)
	w.Plot.Draw(dc)

	for _, tool := range w.tools {
		tool.onRender(c)
	}

	return c.Image()
}

func (w *Widget) pixelToDotsX(x float64) float64 {
	scale := float64(vg.Inch / vgimg.DefaultDPI)
	return x * scale
}
func (w *Widget) pixelToDotsY(y float64) float64 {
	scale := float64(vg.Inch / vgimg.DefaultDPI)
	return (float64(w.Size().Height) - y - float64(w.buttonContainer.Size().Height)) * scale
}

func (w *Widget) resetAxis() {
	w.Plot.X.Min = w.PlotXMin
	w.Plot.X.Max = w.PlotXMax
	w.Plot.Y.Min = w.PlotYMin
	w.Plot.Y.Max = w.PlotYMax
	w.Refresh()
}

func (w *Widget) Dragged(ev *fyne.DragEvent) {
	for _, tool := range w.tools {
		tool.onDrag(ev)
	}
	w.Refresh()
}

func (w *Widget) DragEnd() {
	for _, tool := range w.tools {
		tool.onDragEnd()
	}
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
		gr.graphWidget.buttonContainer.Resize(fyne.NewSize(size.Width, 40))
		for _, button := range gr.graphWidget.buttons {
			button.Resize(fyne.NewSize(130, 40))
		}
	} else {
		gr.graphWidget.buttonContainer.Resize(fyne.NewSize(size.Width, 30))
		for _, button := range gr.graphWidget.buttons {
			button.Resize(fyne.NewSize(70, 30))
		}
	}
	gr.image.Resize(fyne.NewSize(size.Width, size.Height-gr.graphWidget.buttonContainer.Size().Height))
	gr.graphWidget.buttonContainer.Move(fyne.NewPos(0, size.Height-gr.graphWidget.buttonContainer.Size().Height))
	for i, button := range gr.graphWidget.buttons {
		button.Move(fyne.NewPos(5+float32(i)*(button.Size().Width+10), 0))
	}
	gr.Refresh()
}

func (gr *graphRenderer) MinSize() fyne.Size {
	if gr.graphWidget.MinWidgetSize != (fyne.Size{}) {
		return gr.graphWidget.MinWidgetSize
	}
	return fyne.NewSize(100, 150)
}

func (gr *graphRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{gr.image, gr.graphWidget.buttonContainer}
}

func (gr *graphRenderer) Refresh() {
	gr.graphWidget.Refresh()
	gr.image.Refresh()
}
