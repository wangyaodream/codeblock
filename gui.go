package main

import (
	"errors"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/wangyaodream/codeblock/internal/dialogs"
)

type gui struct {
	win   fyne.Window
	title binding.String

	fileTree binding.URITree
}

func (g *gui) makeBanner() fyne.CanvasObject {
	title := canvas.NewText("app Creator", theme.Color(theme.ColorNameForeground))
	title.TextSize = 14
	title.TextStyle = fyne.TextStyle{Bold: true}

	g.title.AddListener(binding.NewDataListener(func() {
		name, _ := g.title.Get()
		if name == "" {
			name = "App Creator"
		}
		title.Text = name
		title.Refresh()
	}))

	home := widget.NewButtonWithIcon("", theme.HomeIcon(), func() {

	})
	left := container.NewHBox(home, title)

	logo := canvas.NewImageFromResource(resourceIconPng)
	logo.FillMode = canvas.ImageFillContain
	return container.NewStack(container.NewPadded(left), container.NewPadded(logo))
}

func (g *gui) makeGUI() fyne.CanvasObject {
	top := g.makeBanner()
	g.fileTree = binding.NewURITree()
	files := widget.NewTreeWithData(g.fileTree, func(branch bool) fyne.CanvasObject {
		return widget.NewLabel("filename.jpg")
	}, func(data binding.DataItem, branch bool, obj fyne.CanvasObject) {
		l := obj.(*widget.Label)
		u, _ := data.(binding.URI).Get()

		name := u.Name()
		l.SetText(name)
	})
	left := widget.NewAccordion(
		widget.NewAccordionItem("Files", files),
		widget.NewAccordionItem("Screens", widget.NewLabel("TODO screens")),
	)
	// 左侧区域默认打开
	left.Open(0)
	right := widget.NewLabel("Right")

	name, _ := g.title.Get()
	window := container.NewInnerWindow(name,
		widget.NewLabel("App Preview"),
	)
	picker := widget.NewSelect([]string{"Desktop", "iPhone 15"}, func(string) {})
	preview := container.NewBorder(container.NewHBox(picker), nil, nil, nil, container.NewCenter(window))

	// 中间区域用灰色背景来区分
	content := container.NewStack(canvas.NewRectangle(color.Gray{Y: 0xee}), preview)
	// 各区域的分隔线
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

// 创建项目
func (g *gui) makeCreateDetail(wizard *dialogs.Wizard) fyne.CanvasObject {
	homeDir, _ := os.UserHomeDir()
	parent := storage.NewFileURI(homeDir)
	chosen, _ := storage.ListerForURI(parent)
	name := widget.NewEntry()
	// 项目名称校验
	name.Validator = func(in string) error {
		if in == "" {
			return errors.New("project name is required")
		}
		return nil
	}
	var dir *widget.Button
	// 选择项目目录
	dir = widget.NewButton(chosen.Name(), func() {
		d := dialog.NewFolderOpen(func(l fyne.ListableURI, err error) {
			if err != nil || l == nil {
				return
			}

			chosen = l
			dir.SetText(l.Name())

		}, g.win)
		d.SetLocation(chosen)
		d.Show()
	})

	// 创建项目表单
	form := widget.NewForm(
		widget.NewFormItem("Name", name),
		widget.NewFormItem("Parent Directory", dir),
	)
	// 表单提交
	form.OnSubmit = func() {
		if name.Text == "" {
			return
		}

		// 创建项目
		project, err := createProject(name.Text, chosen)
		if err != nil {
			dialog.ShowError(err, g.win)
			// 这里return的目的是不让程序退出
			return
		}
		wizard.Hide()
		g.openProject(project)
	}
	return form
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
		wizard.Push("Project Details", g.makeCreateDetail(wizard))
	})
	// 修改create按钮样式
	create.Importance = widget.HighImportance
	buttons := container.NewGridWithColumns(2, open, create)
	home := container.NewVBox(intro, buttons)

	wizard = dialogs.NewWizard("create project", home)
	wizard.Show(win)
	wizard.Resize(home.MinSize().AddWidthHeight(80, 80))
}
