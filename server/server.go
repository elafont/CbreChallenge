package server

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"time"

	"github.com/elafont/CbreChallenge/hangman"
	"github.com/gorilla/mux"
)

const version = "1.0"

// Server implements the http server for the application
type Server struct {
	Games  []*hangman.Hangman
	Router *mux.Router
	Logger *log.Logger
}

func (s *Server) Start(ipWeb string) error {
	srv := http.Server{
		Handler:           s.Router,
		Addr:              ipWeb,
		WriteTimeout:      5 * time.Second,
		ReadTimeout:       5 * time.Second,
		IdleTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
	}

	// Serve
	fmt.Printf("Serving on %s. Press CTRL+C to cancel.\n", ipWeb)
	return srv.ListenAndServe()
}

func (s *Server) RegisterRoutes() {
	s.Router.NotFoundHandler = NotFoundHandler(s)

	// default endpoints
	s.Router.HandleFunc("/", logHandler(s, Root(s))).Methods(http.MethodGet)

	// Ganme endpoints
	s.Router.HandleFunc("/games", logHandler(s, ListGames(s))).Methods(http.MethodGet)
	s.Router.HandleFunc("/newgame", logHandler(s, NewGame(s))).Methods(http.MethodGet)
	s.Router.HandleFunc("/game/{game_id:[0-9]+}", logHandler(s, GetGame(s))).Methods(http.MethodGet)
	s.Router.HandleFunc("/game/{game_id:[0-9]+}/guess/{letter:[a-z]}", logHandler(s, Guess(s))).Methods(http.MethodGet)
	s.Router.HandleFunc("/game/{game_id:[0-9]+}/guess/{letter:[a-z]+}", logHandler(s, WrongGuess(s))).Methods(http.MethodGet)
}

// Respond write handler for responses
func (s *Server) Respond(w http.ResponseWriter, r *Response) {

	if err := r.WriteTo(w); err != nil {
		format := fmt.Sprint("error while writing response", err)

		// Try to get caller information to append to this log
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			s.Logger.Printf(format, err)
			return
		}

		s.Logger.Printf("%s:%d: %s", filepath.Base(file), line, format)
	}
}
