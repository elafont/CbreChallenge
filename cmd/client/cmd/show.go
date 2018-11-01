package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/elafont/CbreChallenge/hangman"
	"github.com/elafont/CbreChallenge/server"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of a given game.",
	Long:  "Show details of a given game.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Showing a game\n\n")
		show()
	},
}

var game *int

func init() {
	rootCmd.AddCommand(showCmd)

	game = showCmd.Flags().IntP("game", "g", 0, "Number of game to show.")

}

func show() {
	hs, err := getGame(host, *game)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(hs)
}

func getGame(srv string, game int) (*hangman.Hstatus, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/game/%d", srv, game))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var answer responseHs

	if err := bindJSON(bytes.NewReader(body), &answer); err != nil {
		return nil, fmt.Errorf("error reading response %v", err)
	}

	if answer.Status == server.StatusFail {
		return nil, fmt.Errorf("Error: Can not show given game, code:%d, %s", answer.Code, answer.Message)
	}

	return answer.Data.Content, nil
}
