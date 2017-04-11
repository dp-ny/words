package words

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestNewDefaultGame(t *testing.T) {
	conf, err := newDefaultConf()
	if err != nil {
		t.Error(err.Error())
	}
	testConfValues(t, conf)
}

func TestNewInvaildGame(t *testing.T) {
	words, err := newDefaultConf()
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	var newDice []Die
	for i, v := range words.DiceConf {
		for j, d := range v {
			if j == 0 {
				continue
			}
			newDice = append(newDice, d)
		}
		words.DiceConf[i] = newDice
	}
}

func TestInvalidJsons(t *testing.T) {
	testDecode(t, oneDieJSON("abc", true), false)
	testDecode(t, oneDieJSON("$!", true), false)
	testDecode(t, oneDieJSON("1234", true), false)
	testDecode(t, oneDieJSON("1234", false), false)
	testDecode(t, oneDieJSON("1234", false), false)
}

func TestInvalidConf(t *testing.T) {
	conf := testDecode(t, "{\"dice\":[[\"AB\",\"AB\"]]}", true)
	err := conf.init()
	if err == nil {
		t.Errorf("Expected error")
	}
}

func testDecode(t *testing.T, input string, valid bool) wordsConf {
	decoder := json.NewDecoder(strings.NewReader(input))
	var b wordsConf
	if err := decoder.Decode(&b); (err == nil) != valid {
		if valid {
			t.Errorf("Expected to be valid, but got %s\n", err.Error())
		} else {
			t.Errorf("Expected to be invalid\n")
		}
	}
	return b
}

func testConfValues(t *testing.T, conf *wordsConf) {
	if conf == nil {
		t.Error("Expected conf to be non-nil")
		t.FailNow()
	}
	if conf.Size == 0 {
		t.Error("Expected conf to have Size")
	}
	if len(conf.DiceConf) == 0 {
		t.Error("Expected conf to have some dice")
	}
	if conf.UnpackDuration().Seconds() == 0 {
		t.Error("Expected conf to have a duration")
	}
}

func oneDieJSON(dieConf string, quote bool) string {
	if quote {
		dieConf = fmt.Sprintf("\"%s\"", dieConf)
	}
	return fmt.Sprintf("{\"dice\":[[%s]]}", dieConf)
}

func TestPrint(t *testing.T) {
	game, err := NewGame()
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	game.print()
}
