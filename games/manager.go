package games

import (
	"math/rand"
	"time"

	"../words"
)

type uuid string

const idLength int = 10

// Manager manages the currently running games
type Manager struct {
	games map[string]*words.Game
}

// NewManager returns a new game Manager, which is the suggested way
// of creating many games
func NewManager() *Manager {
	return &Manager{make(map[string]*words.Game)}
}

// NewGame creates a new game in the manager
func (m Manager) NewGame() (words.Game, error) {
	id := m.generateUnusedGameID()
	g, err := words.NewDefaultGame()
	if err != nil {
		return words.Game{}, err
	}
	m.games[id] = g
	g.ID = id
	g.Time = time.Now()
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
