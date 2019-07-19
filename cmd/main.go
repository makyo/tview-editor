package main

import (
	"github.com/rivo/tview"

	editor "github.com/makyo/tview-editor"
)

func main() {
	e := editor.NewEditor()
	e.SetBorder(true).
		SetTitle("Test editor")
	app := tview.NewApplication().
		SetRoot(e, true).
		SetFocus(e)
	e.SetChangedFunc(func() {
		app.Draw()
	})
	e.SetText("ab")
	if err := app.Run(); err != nil {
		panic(err)
	}
}
