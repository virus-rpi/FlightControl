package Graph

import "fyne.io/fyne/v2/widget"

type ResetAxisTool struct {
	ToolBase
}

func NewResetAxisTool() *ResetAxisTool {
	t := &ResetAxisTool{}
	t.intents.axis = true
	t.intents.button = true
	return t
}

func (t *ResetAxisTool) registerButtons() {
	if t.widget == nil {
		return
	}

	t.widget.buttons = append(t.widget.buttons, widget.NewButton("Reset axis", t.resetAxis))
}

func (t *ResetAxisTool) resetAxis() {
	t.widget.resetAxis()
}
