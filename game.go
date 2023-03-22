package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type Game struct {
	Repo         string
	debug        bool
	drawables    []Drawable
	Screen       tcell.Screen
	DefaultStyle tcell.Style
	Style        tcell.Style
	MaxWidth     int
	// Logger       *log.Logger
}

func (g *Game) DrawStr(x, y int, str string, style *tcell.Style) {
	var st tcell.Style
	if style == nil {
		st = g.DefaultStyle
	} else {
		st = *style
	}
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		g.Screen.SetContent(x, y, c, comb, st)
		x += w
	}
}

func (g *Game) AddDrawable(d Drawable) {
	g.drawables = append(g.drawables, d)
}

func (g *Game) Destroy(d Drawable) {
	newDrawables := []Drawable{}
	for _, dd := range g.drawables {
		if dd == d {
			continue
		}
		newDrawables = append(newDrawables, dd)
	}
	g.drawables = newDrawables
}

func (g *Game) Update() {
	for _, gobj := range g.drawables {
		gobj.Update()
	}
}

func (g *Game) Draw() {
	for _, gobj := range g.drawables {
		gobj.Draw()
	}
}

func (g *Game) FindGameObject(fn func(Drawable) bool) Drawable {
	for _, gobj := range g.drawables {
		if fn(gobj) {
			return gobj
		}
	}
	return nil
}

func (g *Game) FilterGameObjects(fn func(Drawable) bool) []Drawable {
	out := []Drawable{}
	for _, gobj := range g.drawables {
		if fn(gobj) {
			out = append(out, gobj)
		}
	}
	return out
}

type Drawable interface {
	Draw()
	Update()
}

type GameObject struct {
	x             int
	y             int
	w             int
	h             int
	Sprite        string
	Game          *Game
	StyleOverride *tcell.Style
}

func (g *GameObject) Update() {}

func (g *GameObject) Transform(x, y int) {
	g.x += x
	g.y += y
}

func (g *GameObject) Draw() {
	var style *tcell.Style
	if g.StyleOverride != nil {
		style = g.StyleOverride
	}
	lines := strings.Split(g.Sprite, "\n")
	for i, line := range lines {
		l := line
		w := runewidth.StringWidth(line)
		if g.x+w > g.Game.MaxWidth {
			space := g.Game.MaxWidth - g.x
			comb := []rune{}
			for i, r := range line {
				if i > space {
					break
				}
				comb = append(comb, r)
			}
			l = string(comb)
		}
		g.Game.DrawStr(g.x, g.y+i, l, style)
	}
}
