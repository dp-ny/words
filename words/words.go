package words

import (
	"fmt"
	"math/rand"
	"time"
)

type playerID string

type player struct {
}

// Game stores the data for a game ready to be played
type Game struct {
	ID      string
	conf    wordsConf
	Board   board
	Time    time.Time
	players map[playerID]player
}

// Die is responsible for possible values of a grid
type Die struct {
	Values []string
}

func init() {
	rand.Seed(time.Now().Unix())
}

// NewDefaultGame returns a new game ready to be played with the default config
func NewDefaultGame() (*Game, error) {
	conf, err := newDefaultConf()
	if err != nil {
		return nil, err
	}
	return NewGame(*conf), err
}

// NewGame returns a new game ready to be played
func NewGame(conf wordsConf) *Game {
	b := &Game{conf: conf}
	b.init()
	return b
}

func (b *Game) init() {
	i := rand.Intn(len(b.conf.DiceConf))
	dc := b.conf.DiceConf[i]
	var values []string
	for _, d := range dc {
		dValue := d.Values[rand.Intn(len(d.Values))]
		values = append(values, dValue)
	}
	l := b.conf.Size
	b.Board = newArrayBoard(l, l)
	for _, i := range rand.Perm(len(values)) {
		s := newStringValue(values[i])
		if i != 0 {
			b.Board.Set(i/l, i%l, s)
		} else {
			b.Board.Set(0, 0, s)
		}

	}
}

// Time returns an ISO8601 formatted time string
func (b *Game) JsonTime() string {
	return b.Time.Format("2006-01-02T15:04:05-07:00")
}

func (b *Game) print() {
	c := 0
	for x := 0; x < b.conf.Size; x++ {
		for y := 0; y < b.conf.Size; y++ {
			s, ok := b.Board.Get(x, y).(stringValue)
			if !ok {
				panic("Unable to print non stringValue")
			}
			size := len(s.String())
			if size > c {
				c = size
			}
		}
	}
	c++ // spacer
	for x := 0; x < b.conf.Size; x++ {
		for y := 0; y < b.conf.Size; y++ {
			s, ok := b.Board.Get(x, y).(stringValue)
			if !ok {
				panic("Unable to print non stringValue")
			}
			v := s.String()
			for len(v) < c {
				v += "."
			}
			fmt.Print(v)
		}
		fmt.Println()
	}
}
