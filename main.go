package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func updateTime(clock *widget.Label) {
	formated := time.Now().Format("Time: 03:04:05")
	clock.SetText(formated)
}

func updateNum(clock *widget.Label) {
	randomInt := rand.Intn(100) + 1
	randomIntString := strconv.Itoa(randomInt)
	clock.SetText(randomIntString)

}

func main() {
	a := app.New()
	w := a.NewWindow("demo")
	w.SetContent(widget.NewLabel("Hello world!"))
	w.Show()

	w2 := a.NewWindow("Larger")
	w2.SetContent(widget.NewButton("open", func() {
		w3 := a.NewWindow("new widnow")
		w3.SetContent(widget.NewLabel("Third"))
		w3.Show()

	}))
	w2.Resize(fyne.NewSize(100, 100))
	w2.Show()

	a.Run()
}

func foo() {
	fmt.Println("Run in exited！")
}
