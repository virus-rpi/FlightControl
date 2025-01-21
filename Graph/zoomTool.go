package Graph

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/vgimg"
	"image/color"
)

type ZoomTool struct {
	ToolBase
	currentSelection [4]float64
}

func NewZoomTool() *ZoomTool {
	z := &ZoomTool{}
	z.intents.drag = true
	z.intents.render = true
	z.intents.axis = true
	z.intents.button = true
	return z
}

func (z *ZoomTool) registerButtons() {
	if z.widget == nil {
		return
	}

	z.widget.buttons = append(z.widget.buttons, widget.NewButton("Zoom", z.Enable))
}

func (z *ZoomTool) Enable() {
	for _, tool := range z.widget.tools {
		if tool != z && tool.hasIntent("drag") {
			tool.Disable()
		}
	}
	z.ToolBase.Enable()
}

func (z *ZoomTool) onDrag(ev *fyne.DragEvent) {
	if !z.active {
		return
	}
	if z.currentSelection[0] == 0 && z.currentSelection[1] == 0 {
		z.currentSelection[0] = float64(ev.Position.X)
		z.currentSelection[1] = float64(ev.Position.Y)
	}
	z.currentSelection[2] = float64(ev.Position.X)
	z.currentSelection[3] = float64(ev.Position.Y)
}

func (z *ZoomTool) onDragEnd() {
	if !z.active {
		return
	}
	widthPlot := z.widget.Plot.X.Max - z.widget.Plot.X.Min
	xScale := widthPlot / float64(z.widget.Size().Width)
	newXMin := z.currentSelection[0] * xScale
	newXMax := z.currentSelection[2] * xScale
	heightPlot := z.widget.Plot.Y.Max - z.widget.Plot.Y.Min
	yScale := heightPlot / float64(z.widget.Size().Height-z.widget.buttonContainer.Size().Height)
	newYMin := heightPlot - z.currentSelection[1]*yScale
	newYMax := heightPlot - z.currentSelection[3]*yScale
	z.widget.Plot.X.Min = newXMin
	z.widget.Plot.X.Max = newXMax
	z.widget.Plot.Y.Min = newYMin
	z.widget.Plot.Y.Max = newYMax
	z.currentSelection[0] = 0
	z.currentSelection[1] = 0
	z.currentSelection[2] = 0
	z.currentSelection[3] = 0
}

func (z *ZoomTool) onRender(c *vgimg.Canvas) {
	if !z.active {
		return
	}
	if z.currentSelection[0] != 0 || z.currentSelection[1] != 0 || z.currentSelection[2] != 0 || z.currentSelection[3] != 0 {
		rect := vg.Rectangle{
			Min: vg.Point{X: vg.Length(z.widget.pixelToDotsX(z.currentSelection[0])), Y: vg.Length(z.widget.pixelToDotsY(z.currentSelection[1]))},
			Max: vg.Point{X: vg.Length(z.widget.pixelToDotsX(z.currentSelection[2])), Y: vg.Length(z.widget.pixelToDotsY(z.currentSelection[3]))},
		}

		c.SetColor(color.NRGBA{R: 173, G: 216, B: 230, A: 150})
		c.Fill(rect.Path())
	}
}
