package editor

import (
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const (
	cursor = "\u200b"
	selTag = `["selection"]`
	curTag = `["cursor"]`
	endTag = `[""]`
)

type handlerFunc func(*tcell.EventKey) *tcell.EventKey

type Editor struct {
	*tview.TextView
	cursorX     int
	cursorY     int
	cursorIndex int
	handlers    []handlerFunc
	selectable  bool
	lines       []string
}

func NewEditor() *Editor {
	e := &Editor{
		TextView: tview.NewTextView(),
	}
	e.SetWrap(true)
	e.SetRegions(true)
	e.SetDynamicColors(true)
	e.SetInputCapture(e.handleInput)
	e.Highlight("cursor")
	e.cursorX = -1
	e.cursorY = -1
	e.cursorIndex = 0
	return e
}

func (e *Editor) Draw(screen tcell.Screen) {
	e.TextView.Draw(screen)
	ix, iy, _, _ := e.GetInnerRect()
	if e.cursorX == -1 && e.cursorY == -1 {
		e.cursorX, e.cursorY = ix, iy
	}
	screen.ShowCursor(e.cursorX, e.cursorY)
	e.SetTitle(e.GetHighlights()[0])
}

func (e *Editor) SetText(text string) {
	if strings.Index(text, cursor) == -1 {
		text += selTag + curTag + cursor + endTag + endTag
	}
	e.TextView.SetText(text)
}

func (e *Editor) AddHandler(handler handlerFunc) {
	e.handlers = append(e.handlers, handler)
}

func (e *Editor) GetHandlers() []handlerFunc {
	return e.handlers
}

func (e *Editor) SetHandlers(handlers []handlerFunc) {
	e.handlers = handlers
}

func (e *Editor) SetSelectable(selectable bool) {
	e.selectable = selectable
}

func (e *Editor) handleInput(event *tcell.EventKey) *tcell.EventKey {
	for _, handler := range e.handlers {
		event = handler(event)
	}
	if event == nil {
		return nil
	}

	switch event.Key() {
	case tcell.KeyRune:
		e.TextView.SetText(e.TextView.GetText(true) + string(event.Rune()))
		e.cursorX++
		x, _, width, _ := e.GetInnerRect()
		if e.cursorX >= x+width {
			e.cursorX = x
			e.cursorY++
		}
	case tcell.KeyUp, tcell.KeyDown, tcell.KeyRight, tcell.KeyHome, tcell.KeyEnd, tcell.KeyPgUp, tcell.KeyPgDn:
		e.handleMovement(event)
	case tcell.KeyBackspace:
		e.del(-1)
	case tcell.KeyDelete:
		e.del(1)
	}
	return event
}

func (e *Editor) handleMovement(event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyUp:
		e.moveCursorRelative(0, -1, event.Modifiers() == tcell.ModCtrl)
	case tcell.KeyDown:
		e.moveCursorRelative(0, 1, event.Modifiers() == tcell.ModCtrl)
	case tcell.KeyLeft:
		e.moveCursorRelative(-1, 0, event.Modifiers() == tcell.ModCtrl)
	case tcell.KeyRight:
		e.moveCursorRelative(1, 0, event.Modifiers() == tcell.ModCtrl)
	case tcell.KeyHome:
		if event.Modifiers() == tcell.ModCtrl {
			e.moveCursorAbsolute(-1, -1, true, true, event.Modifiers() == tcell.ModCtrl)
		}
		e.moveCursorAbsolute(0, -1, true, false, event.Modifiers() == tcell.ModCtrl)
	case tcell.KeyEnd:
		if event.Modifiers() == tcell.ModCtrl {
			e.moveCursorAbsolute(0, 0, true, true, event.Modifiers() == tcell.ModCtrl)
		}
		e.moveCursorAbsolute(-1, 0, true, false, event.Modifiers() == tcell.ModCtrl)
	case tcell.KeyPgUp:
		e.page(-1, event.Modifiers() == tcell.ModCtrl)
	case tcell.KeyPgDn:
		e.page(1, event.Modifiers() == tcell.ModCtrl)
	}
}

func (e *Editor) moveCursorRelative(dx, dy int, selecting bool) {
}

func (e *Editor) moveCursorAbsolute(x, y int, overflowX, overflowY, selecting bool) {
}

func (e *Editor) page(direction int, selecting bool) {
}

func (e *Editor) del(direction int) {
}
