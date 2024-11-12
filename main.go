package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(newTheme())
	w := a.NewWindow("Hello")
	w.Resize(fyne.Size{Width: 800, Height: 600})

	w.SetContent(makeGUI())
	w.ShowAndRun()
}
