package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/elafont/CbreChallenge/hangman"
	"github.com/gorilla/mux"
)

var s *Server

// CreateRequest creates the request and response for the given path, method and body
func CreateRequest(method string, target string, body io.Reader) (response *httptest.ResponseRecorder, request *http.Request) {
	response = httptest.NewRecorder()
	request = httptest.NewRequest(method, target, body)

	return
}

// Equals performs a deep equal comparison against two
// values and fails if they are not the same.
func Equals(tb testing.TB, expected, actual interface{}) {
	tb.Helper()

	//log.Printf("Equals %[1]v :: %[1]T\n\tgot: %[2]v :: %[2]T\n", expected, actual)
	if !reflect.DeepEqual(expected, actual) {
		tb.Fatalf(
			"\n\texp: %#[1]v (%[1]T)\n\tgot: %#[2]v (%[2]T)\n",
			expected,
			actual,
		)
	}
}

func BindJSON(r io.Reader, target interface{}) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

func TestMain(m *testing.M) {
	// Create Server
	s = &Server{
		Router: mux.NewRouter(),
		Games:  make([]*hangman.Hangman, 0, 8),
		Logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	s.RegisterRoutes()
	os.Exit(m.Run())
}

/*
	Test WrongGuess
	Test NewGame
	Test ListGames
	Test GetGame
	Test Guess
*/

func TestNotFound(t *testing.T) {
	w, r := CreateRequest(
		http.MethodGet,
		"/nonexistant",
		nil,
	)

	s.Router.ServeHTTP(w, r)
	Equals(t, http.StatusNotFound, w.Code)
}

func TestWrongGuess(t *testing.T) {
	w, r := CreateRequest(
		http.MethodGet,
		"/game/99/guess/aa",
		nil,
	)

	s.Router.ServeHTTP(w, r)
	Equals(t, http.StatusNotAcceptable, w.Code)

}

func TestNewGame(t *testing.T) {
	w, r := CreateRequest(
		http.MethodGet,
		"/newgame",
		nil,
	)

	s.Router.ServeHTTP(w, r)
	Equals(t, http.StatusOK, w.Code)

	var body Response

	if err := BindJSON(w.Body, &body); err != nil {
		t.Fatalf("error reading response %v", err)
	}

	Equals(t, "New Game: 0", body.Message)
	// Equals(t, "Game Status", body.Data.Type)

	// var hs = body.Data.Content.(map[string]interface{})

	// Equals(t, 0, int(hs["id"].(float64)))
	// Equals(t, 0, int(hs["tries"].(float64)))
	// Equals(t, "", hs["not_guessed"].(string))
	// Equals(t, false, hs["done"].(bool))
}
