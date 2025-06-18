package ui

import (
	"github.com/jroimartin/gocui"
	"log"
)

type view_info struct {
	haystack []string   // global process list
	maxLines int        // maximum number of lines in a view
	lines    []string   // list of processes shown in current view
	cursor   int        // location of cursor, relative to current lines
	esp      int        // "local stack pointer": index of last process in view, relative to current lines
	ptr      int        // global pointer: index of last process in view, relative to global process list
	output   string     // console output
	needle   string     // search string
	gui      *gocui.Gui // gui object  reference
}

// generate view
func (vi *view_info) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", maxX/8, 3, maxX-maxX/8, maxY-7); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "pshunt"

		_, vi.maxLines = v.Size()
		if vi.lines == nil {
			vi.lines = make([]string, vi.maxLines)
		}
		vi.writeDown(true, 1)

		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	if _, err := g.SetView("help", maxX/8, maxY-6, maxX-maxX/8, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		vi.displayHelp()
	}
	if _, err := g.SetView("console", maxX/8, maxY-3, maxX-maxX/8, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		vi.refreshConsole()
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func initKeybindings(g *gocui.Gui, vi *view_info) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	keyMap := map[interface{}]func(g *gocui.Gui, v *gocui.View) error{
		'k':                vi.cursorUp,
		gocui.KeyArrowUp:   vi.cursorUp,
		'j':                vi.cursorDown,
		gocui.KeyArrowDown: vi.cursorDown,
		gocui.KeyCtrlB:     vi.pageUp,
		gocui.KeyCtrlF:     vi.pageDown,
		'H':                vi.cursorHigh,
		'L':                vi.cursorLow,
		'g':                vi.gotoFirst,
		'G':                vi.gotoLast,
		'r':                vi.refreshPsList,
		'/':                vi.enterEditMode,
		gocui.KeyEnter:     vi.killPs,
		gocui.KeyEsc:       vi.refreshPsList,
	}

	for key, f := range keyMap {
		if err := g.SetKeybinding("main", key, gocui.ModNone, f); err != nil {
			return err
		}
	}

	return nil
}

func StartUI(psList []string) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	vi := view_info{
		haystack: psList,
		maxLines: 0,
		lines:    nil,
		cursor:   1,
		esp:      0,
		ptr:      0,
		output:   "",
		needle:   "",
		gui:      g,
	}

	g.InputEsc = true
	g.ASCII = true

	g.SetManagerFunc(vi.layout)

	if err := initKeybindings(g, &vi); err != nil {
		log.Fatal(err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatal(err)
	}
}
