package ui

import (
    "fmt"
    "sync"
    "github.com/jroimartin/gocui"
)
var (

)

type DisplayManager struct {
    done chan bool
    gui gocui.Gui
    wg sync.WaitGroup
}
func Run() DisplayManager{
    d := DisplayManager{}
    d.done = make(chan bool)
/*
    g, err := gocui.NewGui(gocui.OutputNormal)
    if err != nil {
        log.Panicln(err)
    }
    defer g.Close()

    g.SetManagerFunc(d.Layout)

    if err := d.Keybindings(g); err != nil {
        log.Panicln(err)
    }
    d.gui = *g*/
    return d
}
func (d *DisplayManager) Wg() *sync.WaitGroup{
    return &d.wg
}
func (d *DisplayManager) Done() chan bool{
    return d.done
}
func (d *DisplayManager) SetDone() {
    d.done <- true
}
/*func (d *DisplayManager) Gui() *gocui.Gui{
    return &d.gui
}
func (d *DisplayManager) SetGui(g *gocui.Gui) {
    d.gui = *g
}*/
func (d *DisplayManager) Layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()
    if _, err := g.SetView("ctr", 3, 3, maxX/2+9, maxY/2+8); err != nil {
        fmt.Println(err)
        if err != gocui.ErrUnknownView {
            return err
        }
    }

    if v, err := g.SetView("output", maxX/3-7, 3, maxX/2, maxY-2); err != nil {
        fmt.Println(err)
        if err != gocui.ErrUnknownView {
            return err
        }
        fmt.Fprintln(v, "Welcome to the Game")
    }
    return nil
}

func (d *DisplayManager) Keybindings(g *gocui.Gui) error {
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, d.quit); err != nil {
        return err
    }
    return nil
}

func (d *DisplayManager) quit(g *gocui.Gui, v *gocui.View) error {
    close(d.done)
    g.Close()
    return gocui.ErrQuit
}
