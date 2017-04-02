package words

import aspell "github.com/trustmaster/go-aspell"

// Speller is a spell checker
type Speller struct {
	s aspell.Speller
}

// NewSpeller creates a new speller for checking words
func NewSpeller() (*Speller, error) {
	s, err := getSpellerForLang("en_US")
	if err != nil {
		return nil, err
	}
	return &Speller{s}, nil
}

// Check returns true if word is spelled correctly
func (s Speller) Check(word string) bool {
	return s.s.Check(word)
}

func getSpellerForLang(lang string) (aspell.Speller, error) {
	return aspell.NewSpeller(map[string]string{
		"lang": lang,
	})
}
