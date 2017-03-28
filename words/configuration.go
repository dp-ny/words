package words

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path"
	"unicode"
)

var defaultConf string
var defaultConfErr error

// A wordsConf game to be intiialized and played
type wordsConf struct {
	Size     int
	DiceConf [][]Die `json:"dice"`
}

func init() {
	loadConf()
}

// Validate returns errors if there are issues with the WordsConf
func (b *wordsConf) init() error {
	for i, dc := range b.DiceConf {
		c := len(dc)
		sqrt := int(math.Floor(math.Sqrt(float64(c))))

		if sqrt*sqrt != c {
			return fmt.Errorf(
				"Dice configuration %d did not have a square number of dice (%d)",
				i, c)
		}
		b.Size = sqrt
	}
	return nil
}

// UnmarshalJSON returns a Words configuration based on a json string
func (d *Die) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	for i := 0; i < len(s); i++ {
		c := rune(s[i])
		if !unicode.IsLetter(c) {
			return fmt.Errorf("dice unmarshall: unexpected non character: %s", string(c))
		}
		if !unicode.IsUpper(c) {
			return fmt.Errorf("dice unmarshall: unexpected lowercase character %s", string(c))
		}
		value := string(c)
		if i+1 < len(s) {
			next := rune(s[i+1])
			if unicode.IsLower(next) {
				value += string(next)
				i = i + 1
			}
		}
		d.Values = append(d.Values, value)
	}
	return nil
}

// MarshalJSON returns a json represenation of a die
func (d Die) MarshalJSON() ([]byte, error) {
	var s string
	for _, v := range d.Values {
		s += v
	}
	return []byte(fmt.Sprintf("\"%s\"", s)), nil
}

// newDefaultConf returns a standard conf loaded from json
func newDefaultConf() (*wordsConf, error) {
	if defaultConfErr != nil {
		return nil, defaultConfErr
	}
	return newConf(defaultConf)
}

// newConf attempts to read a provided json file
func newConf(conf string) (*wordsConf, error) {
	file, err := os.Open(conf)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	words := &wordsConf{}
	if err := decoder.Decode(words); err != nil {
		return nil, err
	}
	if err := words.init(); err != nil {
		return nil, err
	}
	return words, nil
}

func (b wordsConf) printJSON() error {
	encoder := json.NewEncoder(os.Stdout)
	if err := encoder.Encode(b); err != nil {
		return err
	}
	return nil
}

func loadConf() {
	wd, err := os.Getwd()
	if err != nil {
		defaultConfErr = err
	}
	dir := wd
	for _, err = os.Stat(getConfDir(dir)); os.IsNotExist(err); _, err = os.Stat(getConfDir(dir)) {
		newDir := path.Clean(path.Join(dir, ".."))
		if newDir == dir {
			defaultConfErr = fmt.Errorf("Unable to find dir, started at %s", wd)
		}
		dir = newDir
	}
	defaultConf = getConfDir(dir)
}

func getConfDir(base string) string {
	return path.Join(base, "conf", "words.json")
}
