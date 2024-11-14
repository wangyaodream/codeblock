package dialogs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

type Wizard struct {
	title   string
	stack   []fyne.CanvasObject
	content *fyne.Container

	d dialog.Dialog
}

func NewWizard(title string, content fyne.CanvasObject) *Wizard {
	w := &Wizard{title: title, stack: []fyne.CanvasObject{content}}
	w.content = container.NewStack(content)
	return w
}

func (w *Wizard) Hide() {
	w.d.Hide()
}

func (w *Wizard) Show(win fyne.Window) {
	w.d = dialog.NewCustomWithoutButtons(w.title, w.content, win)

	w.d.Show()
}
