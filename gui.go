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
	"github.com/wangyaodream/codeblock/internal/editors"
)

type gui struct {
	win   fyne.Window
	title binding.String

	fileTree binding.URITree
	// content  *fyne.Container
	content  *container.DocTabs // 文档标签
	openTabs map[fyne.URI]*container.TabItem
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
	files.OnSelected = func(id widget.TreeNodeID) {
		// 获取选中的文件
		u, err := g.fileTree.GetValue(id)
		if err != nil {
			// 通过dialog显示错误信息
			dialog.ShowError(err, g.win)
			return
		}
		// 打开文件
		g.openFile(u)
	}

	left := widget.NewAccordion(
		widget.NewAccordionItem("Files", files),
		widget.NewAccordionItem("Screens", widget.NewLabel("TODO screens")),
	)
	// 左侧区域默认打开
	left.Open(0)
	// 左侧区域可以同时打开多个
	left.MultiOpen = true

	right := widget.NewRichTextFromMarkdown("## Settings")

	name, _ := g.title.Get()
	window := container.NewInnerWindow(name,
		widget.NewLabel("App Preview"),
	)
	window.CloseIntercept = func() {
		// 关闭窗口时可以做点事
	}
	picker := widget.NewSelect([]string{"Desktop", "iPhone 15"}, func(string) {})
	picker.Selected = "Desktop"
	preview := container.NewBorder(container.NewHBox(picker), nil, nil, nil, container.NewCenter(window))

	// 中间区域用灰色背景来区分，并增加padding
	content := container.NewStack(canvas.NewRectangle(color.Gray{Y: 0xee}), container.NewPadded(preview))

	// 利用文档标签来切换预览和设置，而不是每次都重新创建
	g.content = container.NewDocTabs(
		container.NewTabItem("Preview", content),
	)

	// 关闭文档标签时，删除对应的打开文件
	g.content.CloseIntercept = func(item *container.TabItem) {
		var u fyne.URI
		for child, childItem := range g.openTabs {
			if childItem == item {
				u = child
			}
		}
		if u != nil {
			delete(g.openTabs, u)
		}

		// 最终关闭文档标签
		g.content.Remove(item)
	}

	// 各区域的分隔线
	dividers := [3]fyne.CanvasObject{
		widget.NewSeparator(),
		widget.NewSeparator(),
		widget.NewSeparator(),
	}
	objs := []fyne.CanvasObject{g.content, top, left, right, dividers[0], dividers[1], dividers[2]}
	return container.New(newFysionLayout(top, left, right, g.content, dividers), objs...)
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

func (g *gui) openFile(uri fyne.URI) {
	listable, err := storage.CanList(uri)
	if listable || err != nil {
		// TODO should we unselect this item
		return
	}
	// 判断是否已经打开

	if item, ok := g.openTabs[uri]; ok {
		g.content.Select(item)
		return
	}
	// 打开文件
	editor := editors.ForURI(uri)

	item := container.NewTabItem(uri.Name(), editor)
	// 如果没有打开文件列表，初始化
	if g.openTabs == nil {
		g.openTabs = make(map[fyne.URI]*container.TabItem)
	}
	// 新打开的文件加入到打开文件列表
	g.openTabs[uri] = item
	g.content.Append(item)
	// 当打开文件时，将当前选中的文档标签切换到新打开的文件
	g.content.Select(item)

}
