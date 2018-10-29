package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/elafont/CbreChallenge/hangman"
)

func NotFoundHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		NewResponse(http.StatusNotFound, http.StatusText(http.StatusNotFound), nil).WriteTo(w)
	}
}

func Root(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		NewResponse(
			http.StatusOK,
			fmt.Sprintf("Hangman V%s", version),
			&Data{Type: "Root Version", Content: version}).WriteTo(w)
	}
}

func NewGame(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := len(s.Games)
		game, err := hangman.New(id, "")
		if err != nil {
			NewResponse(http.StatusServiceUnavailable, err.Error(), nil).WriteTo(w)
			return
		}

		s.Games = append(s.Games, game)
		NewResponse(
			http.StatusOK,
			fmt.Sprintf("New Game: %d", id),
			&Data{Type: "Game Status", Content: game.Status()}).WriteTo(w)
		return
	}
}

func ListGames(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var hs = make([]hangman.Hstatus, 0, len(s.Games))

		for _, g := range s.Games {
			hs = append(hs, g.Status())
		}

		NewResponse(
			http.StatusOK,
			fmt.Sprintf("Game List"),
			&Data{Type: "Games", Content: hs}).WriteTo(w)
		return
	}
}

func GetGame(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getInt64Param(r, "game_id")
		if err != nil {
			s.Respond(w, paramError("game_id"))
			return
		}

		if int(id) > len(s.Games) {
			// error, id not available  http.StatusBadRequest
			NewResponse(
				http.StatusBadRequest,
				fmt.Sprintf("ID: %d not available de %d", id, len(s.Games)),
				nil).WriteTo(w)
			return
		}

		game := s.Games[id]
		if game.ID() != int(id) {
			// game id is not in the right place
			// http.StatusInternalServerError
			NewResponse(
				http.StatusInternalServerError,
				fmt.Sprintf("Internal DB Corrupted for ID:%d - %d", id, game.ID()),
				nil).WriteTo(w)
			return
		}

		NewResponse(
			http.StatusOK,
			fmt.Sprintf("Game ID: %d", id),
			&Data{Type: "Game", Content: game.Status()}).WriteTo(w)
		return

	}
}

func Guess(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getInt64Param(r, "game_id")
		if err != nil {
			s.Respond(w, paramError("game_id"))
			return
		}

		if int(id) > len(s.Games) {
			// error, id not available  http.StatusBadRequest
			NewResponse(
				http.StatusBadRequest,
				fmt.Sprintf("ID: %d not available", id),
				nil).WriteTo(w)
			return
		}

		game := s.Games[id]
		if game.ID() != int(id) {
			// game id is not in the right place
			// http.StatusInternalServerError
			NewResponse(http.StatusInternalServerError, "Internal DB Corrupted", nil).WriteTo(w)
			return
		}

		letter, err := getParam(r, "letter")
		if err != nil {
			s.Respond(w, paramError("letter"))
			return
		}

		letter = strings.ToLower(letter)
		hm := game.Guess(letter[0])

		NewResponse(
			http.StatusOK,
			fmt.Sprintf("Guessing letter: %s", letter[:1]),
			&Data{Type: "Game Guess", Content: hm}).WriteTo(w)
	}
}

func WrongGuess(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Respond(w, NewResponse(
			http.StatusNotAcceptable,
			fmt.Sprint("invalid or missing parameter: letter must be only 1 char long"),
			nil))
		return
	}
}
