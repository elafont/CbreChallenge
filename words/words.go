package words

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// This package reads the a file containing a dictionary and offers
// a function to retrieve random words from that dictionary

type wordDict struct {
	words []string
}

const UKdicFile = "./dictionary.UK"

// init Initialize the global random generator
func init() {
	rand.Seed(int64(time.Now().UnixNano()))
}

func (wd *wordDict) readDictionary(dict string) error {
	if wd.words == nil {
		fh, err := os.Open(dict)
		if os.IsNotExist(err) {
			file := filepath.Base(dict)
			dir, err := os.Getwd()
			for {
				testDict := filepath.Join(dir, file)
				fh, err = os.Open(testDict)
				if err != nil {
					if dir != "/" {
						dir = filepath.Dir(dir)
						continue
					} else {
						return err
					}
				}
				break
			}
		} else if err != nil {
			return err
		}

		defer fh.Close()

		data, err := ioutil.ReadAll(fh)
		if err != nil {
			return err
		}

		wd.words = strings.Split(string(data), string('\n'))
	}
	return nil
}

// NewDict returns the dictionary defined by dict
func NewDict(dict string) (*wordDict, error) {
	var wd *wordDict = &wordDict{}
	if dict == "" { // Default dictionary
		dict = UKdicFile
	}

	return wd, wd.readDictionary(dict)
}

// Reads the
func (wd *wordDict) RandomWord() string {
	return wd.words[rand.Intn(len(wd.words))]
}
