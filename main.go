package main

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var richtext = _richtext{}

type _richtext struct{}

func (me _richtext) ui(w fyne.Window) {
	lblMsg := widget.NewLabel("")
	lblMsg.Alignment = fyne.TextAlignCenter

	u, _ := url.Parse("https://github.com/fyne-io/fyne")
	richtext1 := widget.NewRichText(
		&widget.TextSegment{Text: "Text\t", Style: widget.RichTextStyleInline},
		&widget.HyperlinkSegment{Text: "Link", URL: u},
		&widget.SeparatorSegment{},
		&widget.ImageSegment{Source: storage.NewFileURI("dragon.jpeg"), Title: "dragon"},
		&widget.ListSegment{Items: []widget.RichTextSegment{
			&widget.TextSegment{Text: "Text"},
			&widget.HyperlinkSegment{Text: "Link"},
		}, Ordered: true},
		&widget.ParagraphSegment{Texts: []widget.RichTextSegment{
			&widget.TextSegment{Text: "Text"},
			&widget.HyperlinkSegment{Text: "Link"},
		}},
	)
	richtext2 := widget.NewRichTextFromMarkdown("# title1 \n- item*1* \n- item**2** \n## title 3")
	richtext3 := widget.NewRichTextWithText("test1\ntest2")

	btn1 := widget.NewButton("btn1", func() {
	})

	c := container.NewVBox(lblMsg, richtext1, richtext2,
		richtext3, btn1)

	w.SetContent(c)
}

func main() {
	a := app.New()
	w := a.NewWindow("demo")
	richtext.ui(w)

	a.Run()
}
