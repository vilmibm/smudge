package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/vilmibm/smudge/game"
)

const (
	minHeight    = 10
	minWidth     = 12
	animInterval = time.Millisecond * 300
)

type characterCell struct {
	game.GameObject
	HasSpread bool
	Ignited   bool
	HP        int
}

func (c *characterCell) Update() {
	if c.HP <= 0 {
		if !c.HasSpread {
			c.Spread()
		}
		c.Game.AddDrawable(newSmoke(c.Game, c.Point()))
		c.Game.Destroy(c)
		return
	}

	if c.Ignited {
		fireColor := tcell.NewRGBColor(240, int32(rand.Intn(110)+60), 20)
		so := c.Game.Style.Foreground(fireColor)
		c.StyleOverride = &so
		c.HP--
		if !c.HasSpread {
			if rand.Intn(10) < 8 {
				c.Spread()
			}
		}
	}
}

func (c *characterCell) Ignite() {
	c.Ignited = true
}

func (c *characterCell) Spread() {
	c.HasSpread = true
	cells := c.Game.FilterGameObjects(func(d game.Drawable) bool {
		var o *characterCell
		var ok bool
		if o, ok = d.(*characterCell); !ok {
			return false
		}

		if o.Ignited {
			return false
		}

		r := game.NewRay(c.Point(), o.Point())
		if r.Length() == 0 {
			return false
		}

		return r.Length() == 2
	})

	if len(cells) == 0 {
		return
	}

	if len(cells) == 1 {
		cell := cells[0].(*characterCell)
		cell.Ignite()
		return
	}

	for _, d := range cells {
		cell := d.(*characterCell)
		if rand.Intn(100) < 25 {
			cell.Ignite()
		}
	}
}

type smoke struct {
	game.GameObject
	HP int
}

func (s *smoke) Update() {
	if s.HP == 0 || s.Y < 0 {
		s.Game.Destroy(s)
		return
	}
	spriteSheet := "....++++####"
	s.Sprite = string(spriteSheet[s.HP-1])
	color := int32(rand.Intn(120) + 60)
	so := s.Game.Style.Foreground(tcell.NewRGBColor(color, color, color))
	s.StyleOverride = &so
	s.HP--
	s.Y--
	s.X += rand.Intn(3) - 1
}

func newSmoke(g *game.Game, p game.Point) *smoke {
	so := g.Style.Foreground(tcell.NewRGBColor(160, 160, 160))
	return &smoke{
		GameObject: game.GameObject{
			X: p.X, Y: p.Y,
			W: 1, H: 1,
			Sprite:        "*",
			Game:          g,
			StyleOverride: &so,
		},
		HP: 12,
	}
}

func _main(sourceFiles []string) (err error) {
	sources, err := parseFiles(sourceFiles)
	if err != nil {
		return err
	}

	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}

	if err = s.Init(); err != nil {
		return err
	}

	defer s.Fini()

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	quit := make(chan struct{})
	blow := make(chan struct{})
	go inputLoop(s, blow, quit)()

	w, h := s.Size()
	if w < minWidth || h < minHeight {
		return errors.New("terminal is too small i'm sorry")
	}

	gg := &game.Game{
		Screen:   s,
		Style:    defStyle,
		MaxWidth: w,
	}

	smudgeWidth := w / 3

	rand.Seed(time.Now().Unix())

	sourceIx := 0
	sourcePointers := map[int]int{}
	sourceColors := map[int]int{}
	for i := range sources {
		sourcePointers[i] = 0
		sourceColors[i] = rand.Intn(120) + 60
	}

	nextChar := func() (string, tcell.Color) {
		color := int32(sourceColors[sourceIx])
		char := "x"

		source := sources[sourceIx]
		nIx := sourcePointers[sourceIx]
		if nIx < len(source) {
			char = strings.TrimSpace(string(source[nIx]))
			if char == "" {
				char = "+"
			}
			sourcePointers[sourceIx]++
		}
		sourceIx++
		if sourceIx >= len(sources) {
			sourceIx = 0
		}

		return char, tcell.NewRGBColor(color, color, color)
	}

	for y := 0; y < h; y++ {
		for x := smudgeWidth; x < smudgeWidth*2; x++ {
			char, color := nextChar()
			so := defStyle.Foreground(color)
			c := &characterCell{
				GameObject: game.GameObject{
					X: x, Y: y,
					W: 1, H: 1,
					Sprite:        char,
					Game:          gg,
					StyleOverride: &so,
				},
				HP: 20,
			}
			if y == 0 {
				c.Ignited = true
			}
			gg.AddDrawable(c)
		}
	}

	reignite := func() {
		topCells := gg.FilterGameObjects(func(d game.Drawable) bool {
			var cell *characterCell
			var ok bool
			if cell, ok = d.(*characterCell); !ok {
				return false
			}
			return !cell.Ignited
		})

		for ix, d := range topCells {
			if ix == 10 {
				break
			}
			cell := d.(*characterCell)
			cell.Ignite()
		}
	}

	var quitting bool
	for {
		select {
		case <-quit:
			quitting = true
		case <-blow:
			reignite()
		case <-time.After(animInterval):
		}

		if quitting {
			break
		}

		s.Clear()
		gg.Update()
		gg.Draw()
		s.Show()
	}

	return nil
}

func inputLoop(s tcell.Screen, blow chan struct{}, quit chan struct{}) func() {
	return func() {
		for {
			s.Show()

			ev := s.PollEvent()

			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEnter {
					blow <- struct{}{}
				}
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					close(quit)
				}
			}
		}
	}
}

func parseFiles(filenames []string) ([]string, error) {
	out := []string{}
	errs := []error{}
	for _, f := range filenames {
		bs, err := os.ReadFile(f)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		out = append(out, string(bs))
	}

	if len(out) == 0 {
		errMsg := "failed to read any files; errors encountered: "
		for _, e := range errs {
			errMsg += e.Error() + " "
		}
		return nil, fmt.Errorf(errMsg)
	}

	return out, nil

}

func main() {
	err := _main(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
