package manager

import (
	"math/rand"
	"time"

	"../words"
)

type uuid string

const idLength int = 4

// Manager manages the currently running games
type Manager struct {
	def   words.Game
	games map[string]*words.Game
}

// NewManager returns a new game Manager, which is the suggested way
// of creating many games
func NewManager() (*Manager, error) {
	g, err := words.GetDefaultGame()
	if err != nil {
		return nil, err
	}
	return &Manager{*g, make(map[string]*words.Game)}, nil
}

// GetDefaultGame returns the default game for this manager
func (m Manager) GetDefaultGame() words.Game {
	return m.def
}

// NewGame creates a new game in the manager
func (m Manager) NewGame() (words.Game, error) {
	id := m.generateUnusedGameID()
	g, err := words.NewGame()
	if err != nil {
		return words.Game{}, err
	}
	m.games[id] = g
	g.ID = id
	g.SetTime(time.Now().Add(g.Duration()))
	return *g, nil
}

// Get returns the game for the provided id, and true. If the game doesn't exist,
// returns nil, and false
func (m Manager) Get(id string) (*words.Game, bool) {
	g, ok := m.games[id]
	return g, ok
}

func (m Manager) generateUnusedGameID() string {
	s := ""
	for i := 0; i < idLength; i++ {
		r := rand.Intn(26 + 10)
		var c string
		if r-26 >= 0 {
			c = string('0' + r - 26)
		} else {
			c = string('a' + r)
		}
		s = s + c
	}
	if _, ok := m.games[s]; ok {
		return m.generateUnusedGameID()
	}
	return s
}

// DefaultDuration returns the duration of the default game
func (m Manager) DefaultDuration() time.Duration {
	return m.def.Duration()
}
