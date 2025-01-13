package Graph

import (
	"fyne.io/fyne/v2"
	"gonum.org/v1/plot/vg/vgimg"
)

type intents struct {
	drag,
	render,
	axis,
	button bool
}

type tool interface {
	onDrag(ev *fyne.DragEvent)
	onDragEnd()
	onRender(c *vgimg.Canvas)
	setWidget(w *Widget)
	getIntents() intents
	hasIntent(intent string) bool
	registerButtons()

	Disable()
	Enable()
	IsEnabled() bool
}

type ToolBase struct {
	widget  *Widget
	active  bool
	intents intents
}

func (t *ToolBase) getIntents() intents {
	return t.intents
}

func (t *ToolBase) hasIntent(intent string) bool {
	switch intent {
	case "drag":
		return t.intents.drag
	case "render":
		return t.intents.render
	case "axis":
		return t.intents.axis
	case "button":
		return t.intents.button
	}
	return false
}

func (t *ToolBase) setWidget(w *Widget) {
	t.widget = w
}

func (t *ToolBase) Disable() {
	t.active = false
}

func (t *ToolBase) Enable() {
	t.active = true
}

func (t *ToolBase) IsEnabled() bool {
	return t.active
}

func (t *ToolBase) onDrag(ev *fyne.DragEvent) {}

func (t *ToolBase) onDragEnd() {}

func (t *ToolBase) onRender(c *vgimg.Canvas) {}

func (t *ToolBase) registerButtons() {}
