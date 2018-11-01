package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/elafont/CbreChallenge/server"

	"github.com/elafont/CbreChallenge/hangman"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new game.",
	Long:  "Creates a new game.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("New Game\n\n")
		new()
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func new() {
	hs, err := newgame(host)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hs)
}

func newgame(srv string) (*hangman.Hstatus, error) {
	resp, err := http.Get("http://" + srv + "/newgame")
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
		return nil, fmt.Errorf("Error: Can not generate a new game, code:%d, %s", answer.Code, answer.Message)
	}

	return answer.Data.Content, nil
}
