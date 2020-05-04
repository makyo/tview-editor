package editor

import (
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const (
	cursor = "\u2002"
	lTag   = `["left"]`
	rTag   = `["right"]`
	curTag = `["cursor"]`
	selTag = `["selection"]`
	endTag = `[""]`
)

type handlerFunc func(*tcell.EventKey) *tcell.EventKey

type Editor struct {
	*tview.TextView
	handlers   []handlerFunc
	selectable bool
	lines      []string
}

func NewEditor() *Editor {
	e := &Editor{
		TextView: tview.NewTextView(),
	}
	e.SetWrap(true)
	e.SetRegions(true)
	e.SetDynamicColors(true)
	e.SetInputCapture(e.handleInput)
	e.Highlight("selection")
	return e
}

func (e *Editor) Draw(screen tcell.Screen) {
	e.TextView.Draw(screen)
	ix, iy, w, h := e.TextView.GetInnerRect()
	set := false
	for wx := 0; wx < w; wx++ {
		x := ix + wx
		for hy := 0; hy < h; hy++ {
			y := iy + hy
			mainc, combc, style, _ := screen.GetContent(x, y)
			if _, _, attr := style.Decompose(); tcell.AttrDim&attr == tcell.AttrDim {
				screen.ShowCursor(x, y)
				screen.SetContent(x, y, mainc, combc, style.Dim(false))
				set = true
				break
			}
		}
		if set {
			break
		}
	}
}

func (e *Editor) SetText(text string) {
	if len(e.TextView.GetRegionText("cursor")) == 0 {
		text = lTag + tview.Escape(text) + endTag + selTag + curTag + blink(cursor) + endTag + endTag + rTag + endTag
	}
	e.TextView.SetText(text)
}

func (e *Editor) SetTitle(text string) {
	e.TextView.SetTitle(text)
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

func (e *Editor) GetEditorText() string {
	return e.TextView.GetRegionText("left") +
		strings.Replace(e.TextView.GetRegionText("selection"), cursor, "", -1) +
		e.TextView.GetRegionText("right")
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
		e.insertString(string(event.Rune()))
	case tcell.KeyEnter:
		e.insertString("\n")
	case tcell.KeyUp, tcell.KeyDown, tcell.KeyRight, tcell.KeyHome, tcell.KeyEnd, tcell.KeyPgUp, tcell.KeyPgDn:
		e.handleMovement(event)
	case tcell.KeyBS:
		e.del(-1)
	case tcell.KeyDEL:
		e.del(1)
	}
	return event
}

// insertString inserts a string into the text buffer at the point of the
// cursor. If there is any text selected, the seelection is removed.
func (e *Editor) insertString(toAppend string) {
	e.SetText(lTag + tview.Escape(e.TextView.GetRegionText("left")+toAppend) + endTag +
		selTag + tview.Escape(e.TextView.GetRegionText("selection")) + curTag + blink(cursor) + endTag + endTag +
		rTag + tview.Escape(e.TextView.GetRegionText("right")) + endTag)
}

func (e *Editor) handleMovement(event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyUp:
		e.moveCursorRelative(0, -1, shift(event))
	case tcell.KeyDown:
		e.moveCursorRelative(0, 1, shift(event))
	case tcell.KeyLeft:
		e.moveCursorRelative(-1, 0, shift(event))
	case tcell.KeyRight:
		e.moveCursorRelative(1, 0, shift(event))
	case tcell.KeyHome:
		if ctrl(event) {
			e.moveCursorAbsolute(-1, -1, true, true, shift(event))
		}
		e.moveCursorAbsolute(0, -1, true, false, shift(event))
	case tcell.KeyEnd:
		if event.Modifiers() == tcell.ModCtrl {
			e.moveCursorAbsolute(0, 0, true, true, shift(event))
		}
		e.moveCursorAbsolute(-1, 0, true, false, shift(event))
	case tcell.KeyPgUp:
		e.page(-1, shift(event))
	case tcell.KeyPgDn:
		e.page(1, shift(event))
	}
}

func (e *Editor) moveCursorRelative(dx, dy int, selecting bool) {
}

func (e *Editor) moveCursorAbsolute(x, y int, overflowX, overflowY, selecting bool) {
}

func (e *Editor) page(direction int, selecting bool) {
}

func (e *Editor) del(direction int) {
	left := tview.Escape(e.TextView.GetRegionText("left"))
	selected := tview.Escape(e.TextView.GetRegionText("selection"))
	right := tview.Escape(e.TextView.GetRegionText("right"))
	if len(selected) <= 1 {
		if direction > 0 && len(right) > 0 {
			panic("DEL")
			right = right[1:]
		} else if direction < 0 && len(left) > 0 {
			panic("BS")
			left = left[0 : len(left)-1]
		}
	}
	e.SetText(lTag + left + endTag +
		selTag + curTag + blink(cursor) + endTag + endTag +
		rTag + right + endTag)
}

func blink(c string) string {
	return "[::d]" + c + "[::-]"
}

func shift(event *tcell.EventKey) bool {
	return tcell.ModShift&event.Modifiers() == tcell.ModShift
}

func ctrl(event *tcell.EventKey) bool {
	return tcell.ModCtrl&event.Modifiers() == tcell.ModCtrl
}
