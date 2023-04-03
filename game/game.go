package game

import (
	"fmt"
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
	X             int
	Y             int
	W             int
	H             int
	Sprite        string
	Game          *Game
	StyleOverride *tcell.Style
}

func (g *GameObject) Update() {}

func (g *GameObject) Transform(x, y int) {
	g.X += x
	g.Y += y
}

func (g *GameObject) Point() Point {
	return Point{g.X, g.Y}
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
		if g.X+w > g.Game.MaxWidth {
			space := g.Game.MaxWidth - g.X
			comb := []rune{}
			for i, r := range line {
				if i > space {
					break
				}
				comb = append(comb, r)
			}
			l = string(comb)
		}
		g.Game.DrawStr(g.X, g.Y+i, l, style)
	}
}

type Point struct {
	X int
	Y int
}

func (p Point) String() string {
	return fmt.Sprintf("<%d, %d>", p.X, p.Y)
}

func (p Point) Equals(o Point) bool {
	return p.X == o.X && p.Y == o.Y
}

type Ray struct {
	Points []Point
}

func NewRay(a Point, b Point) *Ray {
	r := &Ray{
		Points: []Point{},
	}

	if a.Equals(b) {
		return r
	}

	xDir := 1
	if a.X > b.X {
		xDir = -1
	}
	yDir := 1
	if a.Y > b.Y {
		yDir = -1
	}

	x := a.X
	y := a.Y

	for x != b.X || y != b.Y {
		r.AddPoint(x, y)
		if x != b.X {
			x += xDir * 1
		}
		if y != b.Y {
			y += yDir * 1
		}
	}

	return r
}

func (r *Ray) AddPoint(x, y int) {
	r.Points = append(r.Points, Point{X: x, Y: y})
}

func (r *Ray) Length() int {
	return len(r.Points)
}
