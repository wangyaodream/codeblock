package editors

import (
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var extensions = map[string]func(fyne.URI) fyne.CanvasObject{
	".go":  makeGo,
	".txt": makeTxt,
	".png": makeImg,
	".md":  makeTxt,
}

var mimes = map[string]func(fyne.URI) fyne.CanvasObject{
	"text/plain": makeTxt,
}

func ForURI(u fyne.URI) fyne.CanvasObject {
	ext := strings.ToLower(u.Extension())
	// 这里的editor是一个函数名
	editor, ok := extensions[ext]
	if !ok {
		// 利用mimes就可以打开任何扩展类型的文本文件，只要符合mimes指定的标准
		editor, ok = mimes[u.MimeType()]
		if !ok {
			return widget.NewLabel("Unable to find editor for file:" + u.Name() + ", mime: " + u.MimeType())
		}
		return editor(u)
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

func makeImg(u fyne.URI) fyne.CanvasObject {
	img := canvas.NewImageFromURI(u)

	// 保持图片比例
	img.FillMode = canvas.ImageFillContain

	return img
}
