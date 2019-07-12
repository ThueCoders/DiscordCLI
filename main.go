package main

import (
	"fmt"
	"os"
	"strings"

	"encoding/json"
	"io/ioutil"

	"github.com/ThueCoders/DiscordCLI/logger"
    "github.com/ThueCoders/DiscordCLI/panels"
	"github.com/ThueCoders/discordgo"
    "github.com/BurntSushi/toml"
	"github.com/jroimartin/gocui"
)

const delta          = 1

var (
	token string

	dg *discordgo.Session

	views   = []string{}
	curView = -1
	idxView = 0
)

type Config struct {

}

type keybind struct {

}

func init() {
    var conf Config
    if _, err := toml.Decode("config.toml", &conf); err != nil {
        logger.Log.Println(err)
    }
	token = conf.userToken
}

func initGocui() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		logger.Log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorRed

	g.SetManagerFunc(layout)

	if err := initKeybindings(g); err != nil {
		logger.Log.Panicln(err)
	}

	if err := newView(g); err != nil {
		logger.Log.Panicln(err)
	}

	if err = panels.MakePrompt(g); err != nil {
		logger.Log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		logger.Log.Panicln(err)
	}
}

func initBot() (*discordgo.Session, error) {

	dg, err := discordgo.New(token)
	if err != nil {
		logger.Log.Println("error creating Discord session,", err)
		return nil, err
	}

	dg.AddHandler(messageCreate)
	// dg.AddHandler(channelUpdate)

	err = dg.Open()
	if err != nil {
		logger.Log.Println("error opening connection,", err)
		return nil, err
	}
	return dg, nil
}

func main() {
	dg, err := initBot()
	if err != nil {
		logger.Log.Println(err)
	} else {
		defer dg.Close()
	}
	initGocui()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		if m.Content == "ping" {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		}
		return
	}
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

// func channelUpdate(s *discordgo.Session, c *discordgo.ChannelUpdate) {}

func layout(g *gocui.Gui) error {
	maxX, _ := g.Size()
	v, err := g.SetView("help", maxX-25, 0, maxX-1, 9)
	v0, err := g.SetView("test", 50, 0, maxX-50, 9)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "KEYBINDINGS")
		fmt.Fprintln(v, "Space: New View")
		fmt.Fprintln(v, "Tab: Next View")
		fmt.Fprintln(v, "← ↑ → ↓: Move View")
		fmt.Fprintln(v, "^D: Delete View")
		fmt.Fprintln(v, "t: Set view on top")
		fmt.Fprintln(v, "b: Set view on bottom")
		fmt.Fprintln(v, "^C: Exit")
		fmt.Fprintln(v0, "test")
	}
	return nil
}

func initKeybindings(g *gocui.Gui, keybindings []keybind) error {
    for _, key := range keybindings {
        switch keybind.name {
            case "viewUp":
                if err := g.SetKeybinding(keybind.view, keybind.key, keybind.mod, moveViewUp) {
                    return err
                }
        }
    }
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return newView(g)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return delView(g)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return nextView(g, true)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return moveView(g, v, -delta, 0)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return moveView(g, v, delta, 0)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return moveView(g, v, 0, delta)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return moveView(g, v, 0, -delta)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 't', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			_, err := g.SetViewOnTop(views[curView])
			return err
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'b', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			_, err := g.SetViewOnBottom(views[curView])
			return err
		}); err != nil {
		return err
	}
	return nil
}

func newView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	name := fmt.Sprintf("v%v", idxView)
	v, err := g.SetView(name, maxX/2-5, maxY/2-5, maxX/2+5, maxY/2+5)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		fmt.Fprintln(v, strings.Repeat(name+" ", 30))
	}
	if _, err := g.SetCurrentView(name); err != nil {
		return err
	}

	views = append(views, name)
	curView = len(views) - 1
	idxView += 1
	return nil
}

func delView(g *gocui.Gui) error {
	if len(views) <= 1 {
		return nil
	}

	if err := g.DeleteView(views[curView]); err != nil {
		return err
	}
	views = append(views[:curView], views[curView+1:]...)

	return nextView(g, false)
}

func nextView(g *gocui.Gui, disableCurrent bool) error {
	next := curView + 1
	if next > len(views)-1 {
		next = 0
	}

	if _, err := g.SetCurrentView(views[next]); err != nil {
		return err
	}

	curView = next
	return nil
}

func moveView(g *gocui.Gui, v *gocui.View, dx, dy int) error {
	name := v.Name()
	x0, y0, x1, y1, err := g.ViewPosition(name)
	if err != nil {
		return err
	}
	if _, err := g.SetView(name, x0+dx, y0+dy, x1+dx, y1+dy); err != nil {
		return err
	}
	return nil
}

