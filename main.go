package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"go/build"

	"io/ioutil"

	"github.com/ThueCoders/DiscordCLI/panels"
	"github.com/ThueCoders/discordgo"
	"github.com/jroimartin/gocui"
	toml "github.com/pelletier/go-toml"
)

const delta = 1

var (
	conf Config

	token string

	logger *log.Logger

	dg *discordgo.Session

	views   = []string{}
	curView = -1
	idxView = 0
)

type Keybind struct {
	Name string         `toml:"name"`
	View string         `toml:"view"`
	Key  gocui.Key      `toml:"key"`
	Mod  gocui.Modifier `toml:"mod"`
}

type Config struct {
	BotToken  string    `toml:"botToken"`
	UserToken string    `toml:"userToken"`
	LogPath   string    `toml:"logPath"`
	Keybinds  []Keybind `toml:"keybind"`
}

func init() {
	var contents, err = ioutil.ReadFile("config.toml")
	if err != nil {
		log.Println(err)
	}
	toml.Unmarshal(contents, &conf)
	token = conf.UserToken
	var file, err1 = os.OpenFile(strings.Join([]string{build.Default.GOPATH, conf.LogPath}, ""), os.O_APPEND|os.O_WRONLY, 0664)
	if err1 != nil {
		log.Println(err1)
	}
	logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	logger.Println("Initialization complete")
    logger.Println(conf)
}

func initGocui() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		logger.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorRed

	g.SetManagerFunc(layout)

	if err := initKeybindings(g, conf.Keybinds); err != nil {
		logger.Panicln(err)
	}

	if err := newView(g); err != nil {
		logger.Panicln(err)
	}

	if err = panels.MakePrompt(g); err != nil {
		logger.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		logger.Panicln(err)
	}
}

func initBot() (*discordgo.Session, error) {

	dg, err := discordgo.New(token)
	if err != nil {
		logger.Println("error creating Discord session,", err)
		return nil, err
	}

	dg.AddHandler(messageCreate)
	// dg.AddHandler(channelUpdate)

	err = dg.Open()
	if err != nil {
		logger.Println("error opening connection,", err)
		return nil, err
	}
	return dg, nil
}

func main() {
	dg, err := initBot()
	if err != nil {
		logger.Println(err)
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

func initKeybindings(g *gocui.Gui, keybindings []Keybind) error {
	for _, key := range keybindings {
        // simulates ternary operator. 
        // if we are between 32 and 127 non inclusive gocui expects a rune
        // otherwise it takes a uint16
        tern := map[bool]interface{}{true: rune(key.Key), false: key.Key} [key.Key > 32 && key.Key < 127]
		switch key.Name {
		case "viewUp":
			if err := g.SetKeybinding(key.View, tern, key.Mod, moveViewUp); err != nil {
				return err
			}
		case "viewDown":
			if err := g.SetKeybinding(key.View, tern, key.Mod, moveViewDown); err != nil {
				return err
			}
		case "viewLeft":
			if err := g.SetKeybinding(key.View, tern, key.Mod, moveViewLeft); err != nil {
				return err
			}
		case "viewRight":
            if err := g.SetKeybinding(key.View, tern, key.Mod, moveViewRight); err != nil {
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

func moveViewUp(g *gocui.Gui, v *gocui.View) error {
	name := v.Name()
	x0, y0, x1, y1, err := g.ViewPosition(name)
	if err != nil {
		return err
	}
	if _, err := g.SetView(name, x0, y0-delta, x1, y1-delta); err != nil {
		return err
	}
	return nil
}

func moveViewDown(g *gocui.Gui, v *gocui.View) error {
	name := v.Name()
	x0, y0, x1, y1, err := g.ViewPosition(name)
	if err != nil {
		return err
	}
	if _, err := g.SetView(name, x0, y0+delta, x1, y1+delta); err != nil {
		return err
	}
	return nil
}

func moveViewLeft(g *gocui.Gui, v *gocui.View) error {
	name := v.Name()
	x0, y0, x1, y1, err := g.ViewPosition(name)
	if err != nil {
		return err
	}
	if _, err := g.SetView(name, x0-delta, y0, x1-delta, y1); err != nil {
		return err
	}
	return nil
}

func moveViewRight(g *gocui.Gui, v *gocui.View) error {
	name := v.Name()
	x0, y0, x1, y1, err := g.ViewPosition(name)
	if err != nil {
		return err
	}
	if _, err := g.SetView(name, x0+delta, y0, x1+delta, y1); err != nil {
		return err
	}
	return nil
}
