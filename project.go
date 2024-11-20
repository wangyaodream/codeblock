package main

import (
	"fmt"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/storage"
)

func createProject(name string, parent fyne.ListableURI) (fyne.ListableURI, error) {
	dir, err := storage.Child(parent, name)
	if err != nil {
		return nil, err
	}
	err = storage.CreateListable(dir)
	if err != nil {
		return nil, err
	}
	mod, err := storage.Child(dir, "go.mod")
	if err != nil {
		return nil, err
	}
	w, err := storage.Writer(mod)
	if err != nil {
		return nil, err
	}
	defer w.Close()
	_, err = io.WriteString(w, fmt.Sprintf(`module %s
go 1.23

require fyne.io/fyne/v2 v2.0.0
`, name))
	list, _ := storage.ListerForURI(dir)
	return list, err
}

// 打开项目
func (g *gui) openProject(dir fyne.ListableURI) {
	name := dir.Name()
	g.title.Set(name)

	// 当在已打开的项目状态下再次打开项目需要重置文件树
	g.fileTree.Set(map[string][]string{}, map[string]fyne.URI{})

	items, _ := dir.List()
	for _, uri := range items {
		// TODO handle directories
		g.fileTree.Append(binding.DataTreeRootID, uri.String(), uri)
	}
}
