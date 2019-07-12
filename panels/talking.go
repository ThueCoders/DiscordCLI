package panels

import (
    "fmt"

    "github.com/jroimartin/gocui"
    "github.com/ThueCoders/DiscordCLI/color"
)

const (
	inputHeight    = 2
	inputCursorPos = 17
	promptWidth    = 21
)

func MakePrompt(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	viewHeight := maxY - inputHeight

	// Prompt view
	if v, err := g.SetView("Prompt", 0, viewHeight, promptWidth, viewHeight+inputHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		printPrompt(g)
	}

	// User input view
	if v, err := g.SetView("Input", inputCursorPos, viewHeight, maxX, viewHeight+inputHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Editable = true
		v.Wrap = false
		v.Editor = gocui.EditorFunc(promptEditor)
		if _, err := g.SetCurrentView("Input"); err != nil {
			return err
		}
	}
	return nil
}

func printPrompt(g *gocui.Gui) {
	promptString := "test"

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("Prompt")
		if err != nil {
			return err
		}
		v.Clear()
		v.MoveCursor(0, 0, true)

		if promptString == "a" {
			fmt.Fprintf(v, color.Green(color.Regular, promptString))
		} else {
			fmt.Fprintf(v, color.Red(color.Regular, promptString))
		}
		return nil
	})
}

func promptEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if ch != 0 && mod == 0 {
		v.EditWrite(ch)
		return
	}

	switch key {
	case gocui.KeySpace:
		v.EditWrite(' ')
	case gocui.KeyBackspace, gocui.KeyBackspace2:
		v.EditDelete(true)
	case gocui.KeyDelete:
		v.EditDelete(false)
	case gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case gocui.KeyArrowDown:
		_ = v.SetCursor(len(v.Buffer())-1, 0)
	case gocui.KeyArrowUp:
		v.MoveCursor(0, -1, false)
	case gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
}
