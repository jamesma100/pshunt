package ui

import (
	"fmt"
	"github.com/jamesma100/pshunt/pkg/parser"
	"github.com/jamesma100/pshunt/pkg/runner"
	"github.com/jroimartin/gocui"
	"strings"
)

func (vi *view_info) killPs(g *gocui.Gui, v *gocui.View) error {
	ps := strings.TrimSpace(vi.lines[vi.cursor])
	pid := strings.Split(ps, " ")[0]
	outputMsg := runner.KillPs(pid)
	vi.output = outputMsg
	vi.refreshConsole()
	return nil
}

func (vi *view_info) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if vi.cursor == 1 && vi.ptr-vi.esp+2 == 1 {
		return nil
	}
	if vi.cursor > 1 {
		vi.cursor--
		start_idx := vi.ptr - vi.esp + 1
		vi.writeDown(true, start_idx)
	} else {
		start_idx := max(1, vi.ptr-vi.esp)
		vi.writeDown(true, start_idx)
	}
	return nil
}

func (vi *view_info) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if vi.cursor < vi.maxLines-1 {
		if vi.cursor < vi.esp {
			vi.cursor++
			vi.writeDown(true, vi.ptr-vi.esp+1)
		}
	} else {
		if vi.ptr < len(vi.haystack)-1 {
			vi.writeUp(true, vi.ptr+1)
		}
	}
	return nil
}

func (vi *view_info) cursorHigh(g *gocui.Gui, v *gocui.View) error {
	vi.cursor = 1
	vi.refreshView()
	return nil
}

func (vi *view_info) cursorLow(g *gocui.Gui, v *gocui.View) error {
	vi.cursor = vi.esp
	vi.refreshView()
	return nil
}

func (vi *view_info) pageUp(g *gocui.Gui, v *gocui.View) error {
	prev := vi.ptr - vi.esp
	if prev < 1 {
		return nil
	}
	if prev-vi.maxLines+1 >= 1 {
		vi.writeUp(true, prev)
	} else {
		vi.writeDown(true, 1)
	}
	return nil
}

func (vi *view_info) pageDown(g *gocui.Gui, v *gocui.View) error {
	next := vi.ptr + 1
	if next > len(vi.haystack)-1 {
		return nil
	}
	if next+vi.maxLines-1 < len(vi.haystack) {
		vi.writeDown(true, next)
	} else {
		vi.writeUp(true, len(vi.haystack)-1)
	}
	return nil
}

func (vi *view_info) gotoFirst(g *gocui.Gui, v *gocui.View) error {
	vi.cursor = 1
	vi.writeDown(true, 1)
	return nil
}

func (vi *view_info) gotoLast(g *gocui.Gui, v *gocui.View) error {
	start := max(1, len(vi.haystack)-vi.maxLines+1)
	vi.cursor = vi.esp
	vi.writeDown(true, start)
	return nil
}

func (vi *view_info) writeDown(refreshView bool, start int) {
	vi.lines = make([]string, vi.maxLines)
	vi.lines[0] = vi.haystack[0] // title row is sticky
	i := 1
	j := start
	for ; i < vi.maxLines; i++ {
		if j > len(vi.haystack)-1 {
			break
		}
		vi.lines[i] = vi.haystack[j]
		j++
	}
	vi.esp = i - 1
	vi.ptr = j - 1

	if refreshView {
		vi.refreshView()
	}
}

func (vi *view_info) writeUp(refreshView bool, end int) {
	start := 0
	vi.lines = make([]string, vi.maxLines)
	vi.lines[0] = vi.haystack[0] // title row is sticky
	j := end
	for i := vi.maxLines - 1; i >= start+1; i-- {
		vi.lines[i] = vi.haystack[j]
		j--
	}
	vi.ptr = end

	if refreshView {
		vi.refreshView()
	}
}

func (vi *view_info) refreshPsList(g *gocui.Gui, v *gocui.View) error {
	v.Clear()

	contents := runner.GetPsList()
	psList := parser.ParseList(contents)

	vi.haystack = psList
	vi.cursor = 1
	vi.writeDown(true, 1)
	vi.output = "Process list refreshed"
	vi.refreshConsole()
	return nil
}

func (vi *view_info) refreshView() {
	v, _ := vi.gui.View("main")
	v.Clear()
	for idx, line := range vi.lines {
		if idx == vi.cursor {
			fmt.Fprintf(v, "\x1b[0;30;47m")
			fmt.Fprintf(v, line)
			fmt.Fprintf(v, "\x1b[m")
		} else {
			fmt.Fprintf(v, line)
		}
		fmt.Fprintf(v, "\n")
	}
}

func (vi *view_info) filterPsList(v *gocui.View) {
	str, _ := v.Line(0)
	vi.needle = str[1:]
	vi.gui.Cursor = false
	vi.gui.SetCurrentView("main")
	vi.refreshPsList(vi.gui, v)
	filtered := runner.Grep(vi.needle, vi.haystack[1:])
	vi.haystack = append([]string{vi.haystack[0]}, filtered...)
	vi.writeDown(true, 1)
	vi.refreshView()
}

func (vi *view_info) simpleEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		x, _ := v.Cursor()
		if x > 1 {
			v.EditDelete(true)
		}
	case key == gocui.KeyArrowLeft:
		x, _ := v.Cursor()
		if x > 1 {
			v.MoveCursor(-1, 0, false)
		}
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	case key == gocui.KeyEsc:
		vi.gui.SetCurrentView("main")
		vi.refreshPsList(vi.gui, v)
	case key == gocui.KeyEnter:
		vi.filterPsList(v)
	}
}

func (vi *view_info) enterEditMode(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView("console"); err != nil {
		return err
	}
	cv, _ := vi.gui.View("console")
	cv.Editable = true
	cv.Clear()
	cv.SetCursor(0, 0)
	cv.EditWrite('/')
	cv.Editor = gocui.EditorFunc(vi.simpleEditor)
	g.Cursor = true

	return nil
}

func (vi *view_info) refreshConsole() {
	v, _ := vi.gui.View("console")
	v.Clear()
	fmt.Fprintf(v, "%s\n", vi.output)
}

func (vi *view_info) displayHelp() {
	v, _ := vi.gui.View("help")
	v.Clear()
	fmt.Fprintf(v, "k/j: up/down | /: search | enter: sigkill | r: refresh | ctrl-c: quit")
}
