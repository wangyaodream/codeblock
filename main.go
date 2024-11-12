package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func (g *gui) makeMenu() *fyne.MainMenu {
	file := fyne.NewMenu("File",
		fyne.NewMenuItem("Open Project", g.openProject),
	)
	return fyne.NewMainMenu(file)
}

func main() {
	a := app.New()
	a.Settings().SetTheme(newTheme())
	w := a.NewWindow("Hello")
	w.Resize(fyne.Size{Width: 800, Height: 600})

	ui := &gui{win: w}

	w.SetContent(ui.makeGUI())
	w.SetMainMenu(ui.makeMenu())
	ui.openProject()
	w.ShowAndRun()
}
