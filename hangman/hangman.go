package hangman

import (
	"time"

	"github.com/elafont/CbreChallenge/words"
)

type Hangman struct {
	id         int
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

// NewHangman returns a play with a new word to guess
func New(id int, dict string) (*Hangman, error) {
	wd, err := words.NewDict(dict)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	hm := Hangman{
		id:      id,
		word:    wd.RandomWord(),
		start:   t,
		end:     t,
		letters: make([]byte, 255),
	}

	return &hm, nil
}

// ID getter for id
func (hm *Hangman) ID() int {
	return hm.id
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
	var wdisplay = make([]byte, 0)
	var failed = make([]byte, 0)

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
		Id:           hm.id,
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

/* ToDo:
 */
