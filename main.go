package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	minHeight    = 10
	minWidth     = 12
	animInterval = time.Millisecond * 100
)

/*
	the overall goal of this program is to interleave a set of text files
	together into a smudge stick and then burn them away character by character.
	the purpose of this program is to aid a user in meditating over their
	computer as a place and a space as opposed to just a tool or content delivery
	mechanism.

	input files are listed ad naueseum as positional arguments.
*/

func _main(sources []string) (err error) {
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
	go inputLoop(s, quit)()

	w, h := s.Size()
	if w < minWidth || h < minHeight {
		return errors.New("terminal is too small i'm sorry")
	}

	game := &Game{
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
		sourceColors[i] = rand.Intn(120)
		sourceColors[i] += 40
	}

	nextChar := func() (string, tcell.Color) {
		color := int32(sourceColors[sourceIx])
		char := "x"

		source := sources[sourceIx]
		nIx := sourcePointers[sourceIx]
		if nIx < len(source) {
			char = string(source[nIx])
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
			game.AddDrawable(&GameObject{
				x: x, y: y,
				w: 1, h: 1,
				Sprite:        char,
				Game:          game,
				StyleOverride: &so,
			})
		}
	}

	var quitting bool
	for {
		select {
		case <-quit:
			quitting = true
		case <-time.After(animInterval):
		}

		if quitting {
			break
		}

		s.Clear()

		game.Update()
		game.Draw()

		// TODO burn animation

		s.Show()
	}

	return nil
}

func inputLoop(s tcell.Screen, quit chan struct{}) func() {
	return func() {
		for {
			s.Show()

			ev := s.PollEvent()

			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					close(quit)
				}
			}
		}
	}
}

func main() {
	// TODO parse each argument, read file, do the interleaving

	dummyTexts := []string{`Celia agreed with her sister, but she did not say so. The two little
girls had been sitting by the fireside, for the April evening
was chilly; but now the daylight had nearly faded, and Joy, rising,
went to the door and peeped into the passage to make certain that
Jane had lit the gas there. Satisfied on that point, she returned
to her former seat by the fire, and continued the conversation.`,

		`"I wonder if we ought to send Jane to the drawing-room to light
the gas?" Celia suggested presently. "But, no, mother would be sure
to ring if she wished it. Oh, the gentleman's going at last!"`,

		`There was a sound of footsteps in the passage. The front door opened
and shut, and the next minute Mrs. Wallis joined her little
daughters. She was a tall, stately woman with a pale, handsome face,
and hair which was prematurely grey.`,

		`"My visitor kept me some time," she remarked, as she seated herself
in an easy chair, and glanced from one to the other of the children.
"I suppose you have been cogitating about him, and wondering who he
could possibly be?"`}

	err := _main(dummyTexts)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
