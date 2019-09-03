package panels

import (
    "fmt"
	"github.com/jroimartin/gocui"
    "github.com/ThueCoders/DiscordCLI/logger"
)

const (
    rows = 7
    cols = 7
)
func LoadDefaultPanels(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	serverV, errS := g.SetView("server", 0, 0, maxX/cols, maxY)
	guildV, errG := g.SetView("guild", maxX/cols, 0, 2*maxX/cols, maxY)
	readV, errR := g.SetView("reading", 2*maxX/cols, 0, 6*maxX/cols, 6*maxY/rows)
	talkV, errT := g.SetView("talking", 2*maxX/cols, 6*maxY/rows, 6*maxX/cols, maxY)
	memberV, errM := g.SetView("member", 6*maxX/cols, 0, maxX, maxY)
	if errS != nil && errS != gocui.ErrUnknownView {
        logger.Log.Println(errS)
		return errS
	}
	if errG != nil && errG != gocui.ErrUnknownView {
        logger.Log.Println(errG)
		return errG
	}
	if errT != nil && errT != gocui.ErrUnknownView {
        logger.Log.Println(errT)
		return errT
	}
	if errR != nil && errR != gocui.ErrUnknownView {
        logger.Log.Println(errR)
		return errR
	}
	if errM != nil && errM != gocui.ErrUnknownView {
        logger.Log.Println(errM)
		return errM
	}
	fmt.Fprintln(serverV, "server panel")
	fmt.Fprintln(guildV, "guild panel")
	fmt.Fprintln(talkV, "talking panel")
	fmt.Fprintln(readV, "reading panel")
	fmt.Fprintln(memberV, "member panel")
	return nil
}
