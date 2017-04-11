package words

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"path"
	"time"
)

type playerID string

type player struct {
}

// Game stores the data for a game ready to be played
type Game struct {
	ID       string
	conf     wordsConf
	board    board
	time     time.Time
	duration time.Duration
	stopped  bool
	players  map[playerID]player
}

// Die is responsible for possible values of a grid
type Die struct {
	Values []string
}

var templateName = "_wordsTable.html"
var templatePath = "../words/views/" + templateName
var gameTemplate *template.Template

var baseBoard [][]string

func init() {
	rand.Seed(time.Now().Unix())
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path := path.Join(wd, "words", templatePath)
	gameTemplate, err = template.New("words").ParseFiles(path)
	if err != nil {
		fmt.Printf("Unable to load template for words game, attempted: %v\n", path)
		panic(err)
	}
}

// GetDefaultGame returns the default game to be displaed when then view is loaded
func GetDefaultGame() (*Game, error) {
	g, err := NewGame()
	if err != nil {
		return g, err
	}
	var board [][]string
	board = append(board, []string{"W", "O", "R", "D"})
	board = append(board, []string{"S", "#", "#", "#"})
	board = append(board, []string{"#", "B", "Y", "#"})
	board = append(board, []string{"D", "P", "N", "Y"})
	for i, v := range board {
		for j, w := range v {
			g.board.Set(i, j, newStringValue(w))
		}
	}
	return g, nil
}

// NewGame returns a new game ready to be played with the default config
func NewGame() (*Game, error) {
	conf, err := newDefaultConf()
	if err != nil {
		return nil, err
	}
	return NewGameWithConf(*conf), err
}

// NewGameWithConf returns a new game ready to be played
func NewGameWithConf(conf wordsConf) *Game {
	g := &Game{conf: conf}
	g.init()
	return g
}

func (g *Game) init() {
	g.stopped = true
	g.duration = g.conf.UnpackDuration()
	i := rand.Intn(len(g.conf.DiceConf))
	dc := g.conf.DiceConf[i]
	var values []string
	for _, d := range dc {
		dValue := d.Values[rand.Intn(len(d.Values))]
		values = append(values, dValue)
	}
	l := g.conf.Size
	g.board = newArrayBoard(l, l)
	for _, i := range rand.Perm(len(values)) {
		s := newStringValue(values[i])
		if i != 0 {
			g.board.Set(i/l, i%l, s)
		} else {
			g.board.Set(0, 0, s)
		}

	}
}

//SetTime sets the time on the game to the provided time
func (g *Game) SetTime(t time.Time) {
	g.time = t
}

// JSONTime returns an ISO8601 formatted time string
func (g *Game) JSONTime() string {
	return g.time.Format("2006-01-02T15:04:05-07:00")
}

//Duration returns the duration
func (g *Game) Duration() time.Duration {
	return g.duration
}

//DurationMs returns the game's duration as an int64 in milliseconds
func (g *Game) DurationMs() int64 {
	return int64(g.duration.Seconds() * 1000)
}

//SetStopped sets the game as stopped, and stores the remaining time
func (g *Game) SetStopped(s bool) {
	if g.stopped == s {
		return
	}

	if g.stopped != s {
		g.stopped = s
		if !s {
			g.time = time.Now().Add(g.duration)
		} else {
			g.duration = time.Now().Sub(g.time)
		}
	}
}

//Stopped returns whether or not the game is stopped
func (g *Game) Stopped() bool {
	return g.stopped
}

//JSON returns a map representing this game in json
func (g *Game) JSON() (map[string]interface{}, error) {
	data := make(map[string]interface{})
	data["id"] = g.ID
	data["duration"] = g.DurationMs()
	data["time"] = g.JSONTime()
	data["stopped"] = g.stopped

	w := new(bytes.Buffer)
	err := gameTemplate.ExecuteTemplate(w, templateName, g.board.ToStringArray())
	if err != nil {
		return nil, err
	}
	data["html"] = template.HTML(w.String())
	return data, nil
}

func (g *Game) print() {
	c := 0
	for x := 0; x < g.conf.Size; x++ {
		for y := 0; y < g.conf.Size; y++ {
			s, ok := g.board.Get(x, y).(stringValue)
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
	for x := 0; x < g.conf.Size; x++ {
		for y := 0; y < g.conf.Size; y++ {
			s, ok := g.board.Get(x, y).(stringValue)
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
