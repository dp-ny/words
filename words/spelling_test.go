package words

import (
	"testing"

	"github.com/trustmaster/go-aspell"
)

func TestAspellSanity(t *testing.T) {
	speller := getTestSpellChecker()
	defer speller.Delete()
	word := "testing"
	if !speller.Check(word) {
		t.Error("Expected valid word,", word)
	}
	word = "tsting"
	if speller.Check(word) {
		t.Error("Expected invalid word,", word)
	}
	s := speller.Suggest(word)
	if len(s) == 0 {
		t.Error("Expected to find some suggestions")
	}
}

func getTestSpellChecker() aspell.Speller {
	speller, err := aspell.NewSpeller(map[string]string{
		"lang": "en_US",
	})
	if err != nil {
		panic(err.Error())
	}
	return speller
}
