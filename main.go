package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/storage"
)

func (g *gui) makeMenu() *fyne.MainMenu {
	file := fyne.NewMenu("File",
		fyne.NewMenuItem("Open Project", g.openProjectDialog),
	)
	return fyne.NewMainMenu(file)
}

func main() {
	a := app.New()
	a.Settings().SetTheme(newTheme())
	w := a.NewWindow("Hello")
	w.Resize(fyne.Size{Width: 800, Height: 600})

	ui := &gui{win: w, title: binding.NewString()}

	w.SetContent(ui.makeGUI())
	w.SetMainMenu(ui.makeMenu())
	ui.title.AddListener(binding.NewDataListener(func() {
		name, _ := ui.title.Get()
		w.SetTitle("App:" + name)
	}))

	flag.Usage = func() {
		fmt.Println("Usage: codeblock [project directory]")
	}

	flag.Parse()
	if flag.NArg() > 0 {
		dirPath := flag.Args()[0]
		dirPath, err := filepath.Abs(dirPath)
		if err != nil {
			fmt.Println("Error path:", err)
			return
		}
		dirURI := storage.NewFileURI(dirPath)
		dir, err := storage.ListerForURI(dirURI)
		if err != nil {
			fmt.Println("Error opening project:", err)
			return
		}
		ui.openProject(dir)

	} else {
		ui.openProjectDialog()
	}
	w.ShowAndRun()
}
