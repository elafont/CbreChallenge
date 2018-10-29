package server

import (
	"encoding/json"
	"fmt"
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
	Test Guess
*/

func testNotFound(t *testing.T) {
	w, r := CreateRequest(
		http.MethodGet,
		"/nonexistant",
		nil,
	)

	s.Router.ServeHTTP(w, r)
	Equals(t, http.StatusNotFound, w.Code)
}

func testWrongGuess(t *testing.T) {
	w, r := CreateRequest(
		http.MethodGet,
		"/game/99/guess/aa",
		nil,
	)

	s.Router.ServeHTTP(w, r)
	Equals(t, http.StatusNotAcceptable, w.Code)

}

func testNewGame(t *testing.T) {
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

	Equals(t, fmt.Sprintf("New Game: %d", len(s.Games)-1), body.Message)
}

func testGetGame(t *testing.T) {
	w, r := CreateRequest(
		http.MethodGet,
		"/game/0",
		nil,
	)

	s.Router.ServeHTTP(w, r)
	Equals(t, http.StatusOK, w.Code)

	var body Response

	if err := BindJSON(w.Body, &body); err != nil {
		t.Fatalf("error reading response %v", err)
	}

	Equals(t, "Game ID: 0", body.Message)
	Equals(t, "Game Status", body.Data.Type)

	var hs = body.Data.Content.(map[string]interface{})

	Equals(t, 0, int(hs["id"].(float64)))
	Equals(t, 0, int(hs["tries"].(float64)))
	Equals(t, "", hs["not_guessed"].(string))
	Equals(t, false, hs["done"].(bool))
}

func testListGames(t *testing.T) {
	testNewGame(t)
	w, r := CreateRequest(
		http.MethodGet,
		"/games",
		nil,
	)

	s.Router.ServeHTTP(w, r)
	Equals(t, http.StatusOK, w.Code)

	var body Response

	if err := BindJSON(w.Body, &body); err != nil {
		t.Fatalf("error reading response %v", err)
	}

	Equals(t, "Game List", body.Message)
	Equals(t, "Games", body.Data.Type)
	if len(body.Data.Content.([]interface{})) != 2 {
		t.Fatal("error getting list, expecting 2 games")
	}
}

func testGuessGame(t *testing.T) {
	var lastDone bool

	for letter := 97; letter < 123; letter++ {
		w, r := CreateRequest(
			http.MethodGet,
			fmt.Sprintf("/game/0/guess/%s", string(letter)),
			nil,
		)

		s.Router.ServeHTTP(w, r)
		Equals(t, http.StatusOK, w.Code)

		var body Response

		if err := BindJSON(w.Body, &body); err != nil {
			t.Fatalf("error reading response %v", err)
		}
		var hs = body.Data.Content.(map[string]interface{})
		lastDone = hs["done"].(bool)
		if lastDone == true {
			break
		}
	}

	Equals(t, true, lastDone)
}

func TestGames(t *testing.T) {
	t.Run("Not Found", testNotFound)
	t.Run("New Game", testNewGame)
	t.Run("Wrong Guess", testWrongGuess)
	t.Run("Get Game", testGetGame)
	t.Run("List Games", testListGames)
	t.Run("Guess Game", testGuessGame)
}
