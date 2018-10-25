package hangman

import (
	"time"

	"github.com/elafont/CbreChallenge/words"
)

type Hangman struct {
	Id         int
	word       string
	start, end time.Time
	done       bool
	tries      int

	letters []byte
}

type Hstatus struct {
	Id           int    `json:"id"`
	GuessedSoFar string `json:"guessed_so_far"`
	Tries        int    `json:"tries"`
	Start        string `json:"start_time"`
	Failed       string `json:"not_guessed"`
	Done         bool   `json:"done"`
}

// Database keeping track of games asked
var Games []Hangman = make([]Hangman, 0, 10)

// NewHangman returns a play with a new word to guess
func New(dict string) (*Hangman, error) {
	wd, err := words.NewDict(dict)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	hm := Hangman{
		Id:      len(Games),
		word:    wd.RandomWord(),
		start:   t,
		end:     t,
		letters: make([]byte, 255),
	}

	Games = append(Games, hm)

	return &hm, nil
}

// Guess evaluates if a letter is part of the hidden word
func (hm *Hangman) Guess(letter byte) Hstatus {
	hm.end = time.Now()
	hm.tries++
	// letters will keep 0 when unguessed, 1 when guessed right and 255 when guessed wrong
	if hm.letters[letter] == 0 {
		hm.letters[letter] = 255
		// check if letter is available on hm.word
		for _, l := range hm.word {
			if l == rune(letter) {
				hm.letters[letter] = 1
			}
		}
	}
	return hm.Status()
}

// Status, returns the HStatus type with info on the game
func (hm *Hangman) Status() Hstatus {
	var wdisplay = make([]byte, 0, len(hm.word))
	var failed = make([]byte, 0, hm.tries)

	for _, l := range hm.word {
		if hm.letters[l] == 1 {
			wdisplay = append(wdisplay, byte(l))
		} else {
			wdisplay = append(wdisplay, '*')
		}
	}

	for k, l := range hm.letters {
		if l == 255 {
			failed = append(failed, byte(k))
		}
	}

	if hm.word == string(wdisplay) {
		hm.done = true
	}

	st := Hstatus{
		Id:           hm.Id,
		GuessedSoFar: string(wdisplay),
		Start:        hm.start.String()[:19],
		Tries:        hm.tries,
		Failed:       string(failed),
		Done:         hm.done,
	}

	return st
}

// func (hm *Hangman) Show() string {
// 	return hm.word
// }
