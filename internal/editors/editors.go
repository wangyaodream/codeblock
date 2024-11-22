package editors

import (
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var extensions = map[string]func(fyne.URI) fyne.CanvasObject{
	".go": makeGo,
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
	// TODO code editor
	code := widget.NewEntry()
	code.TextStyle = fyne.TextStyle{Monospace: true}

	r, err := storage.Reader(u)
	if err != nil {
		code.SetText("Unable to read" + u.Name())
		return code
	}

	defer r.Close()

	io.ReadAll(r)
	data, _ := io.ReadAll(r)
	code.SetText(string(data))
	return code
}
