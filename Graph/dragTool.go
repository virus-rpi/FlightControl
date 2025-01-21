package Graph

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type DragTool struct {
	ToolBase
}

func NewDragTool() *DragTool {
	d := &DragTool{}
	d.intents.drag = true
	d.intents.axis = true
	d.intents.button = true
	return d
}

func (d *DragTool) registerButtons() {
	if d.widget == nil {
		return
	}

	d.widget.buttons = append(d.widget.buttons, widget.NewButton("Drag", d.Enable))
}

func (d *DragTool) Enable() {
	for _, tool := range d.widget.tools {
		if tool != d && tool.hasIntent("drag") {
			tool.Disable()
		}
	}
	d.ToolBase.Enable()
}

func (d *DragTool) onDrag(ev *fyne.DragEvent) {
	if !d.active {
		return
	}
	d.widget.Plot.X.Min -= float64(ev.Dragged.DX) * (d.widget.Plot.X.Max - d.widget.Plot.X.Min) / float64(d.widget.Size().Width)
	d.widget.Plot.X.Max -= float64(ev.Dragged.DX) * (d.widget.Plot.X.Max - d.widget.Plot.X.Min) / float64(d.widget.Size().Width)
	d.widget.Plot.Y.Min += float64(ev.Dragged.DY) * (d.widget.Plot.Y.Max - d.widget.Plot.Y.Min) / float64(d.widget.Size().Height)
	d.widget.Plot.Y.Max += float64(ev.Dragged.DY) * (d.widget.Plot.Y.Max - d.widget.Plot.Y.Min) / float64(d.widget.Size().Height)
}
