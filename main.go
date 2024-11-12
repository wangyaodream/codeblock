package main

import "fyne.io/fyne/v2/app"

func main() {
	a := app.New()
	a.Settings().SetTheme(newTheme())
	w := a.NewWindow("Hello")

	w.SetContent(makeGUI())
	w.ShowAndRun()
}
