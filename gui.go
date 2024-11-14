package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/wangyaodream/codeblock/internal/dialogs"
)

type gui struct {
	win   fyne.Window
	title binding.String
}

func makeBanner() fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), func() {}),
	)
	logo := canvas.NewImageFromResource(resourceIconPng)
	logo.FillMode = canvas.ImageFillContain
	return container.NewStack(toolbar, container.NewPadded(logo))
}

func (g *gui) makeGUI() fyne.CanvasObject {
	top := makeBanner()
	left := widget.NewLabel("Left")
	right := widget.NewLabel("Right")

	directory := widget.NewLabelWithData(g.title)
	content := container.NewStack(canvas.NewRectangle(color.Gray{Y: 0xee}), directory)
	// return container.NewBorder(makeBanner(), nil, left, right, content)
	dividers := [3]fyne.CanvasObject{
		widget.NewSeparator(),
		widget.NewSeparator(),
		widget.NewSeparator(),
	}
	objs := []fyne.CanvasObject{content, top, left, right, dividers[0], dividers[1], dividers[2]}
	return container.New(newFysionLayout(top, left, right, content, dividers), objs...)
}

func (g *gui) openProjectDialog() {
	dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, g.win)
			return
		}
		if dir == nil {
			return
		}
		g.openProject(dir)
	}, g.win)
}

func (g *gui) openProject(dir fyne.ListableURI) {
	name := dir.Name()
	g.title.Set(name)
}

func (g *gui) ShowCreate(win fyne.Window) {
	var wizard *dialogs.Wizard
	intro := widget.NewLabel(`Here you can create new project!
Or open an existing project.`)

	open := widget.NewButton("Open Project", func() {
		wizard.Hide()
		g.openProjectDialog()
	})
	create := widget.NewButton("Create Project", func() {
		step2 := widget.NewLabel("step 2 content")
		wizard.Push("Step 2", step2)
	})
	// 修改create按钮样式
	create.Importance = widget.HighImportance
	buttons := container.NewGridWithColumns(2, open, create)
	home := container.NewVBox(intro, buttons)

	wizard = dialogs.NewWizard("create project", home)
	wizard.Show(win)
	wizard.Resize(home.MinSize().AddWidthHeight(40, 80))
}
