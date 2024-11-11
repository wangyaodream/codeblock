package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeGUI() fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), func() {}),
	)
	left := widget.NewLabel("Left")
	right := widget.NewLabel("Right")
	content := widget.NewLabel("Content")
	return container.NewBorder(toolbar, nil, left, right, content)
}
