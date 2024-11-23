package editors

import (
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var extensions = map[string]func(fyne.URI) fyne.CanvasObject{
	".go":  makeGo,
	".txt": makeTxt,
	".md":  makeTxt,
}

func ForURI(u fyne.URI) fyne.CanvasObject {
	// 这里的editor是一个函数名
	editor, ok := extensions[u.Extension()]
	if !ok {
		return widget.NewLabel("Unable to find editor for file" + u.Name())
	}
	// 需要调用它并返回
	return editor(u)
}

func makeGo(u fyne.URI) fyne.CanvasObject {
	code := makeTxt(u)
	code.(*widget.Entry).TextStyle = fyne.TextStyle{Monospace: true}

	return code
}

func makeTxt(u fyne.URI) fyne.CanvasObject {
	code := widget.NewEntry()

	r, err := storage.Reader(u)
	if err != nil {
		code.SetText("Unable to read" + u.Name())
		return code
	}

	defer r.Close()

	data, _ := io.ReadAll(r)
	code.SetText(string(data))
	return code
}
