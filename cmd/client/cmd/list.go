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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available games",
	Long:  `List all available games.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("List of Games\n\n")
		list()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func list() {
	hs, err := listgame(host)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, game := range *hs {
		fmt.Println(game)
	}

}

func listgame(srv string) (*[]hangman.Hstatus, error) {
	resp, err := http.Get("http://" + srv + "/games")
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var answer responseHsArr

	if err := bindJSON(bytes.NewReader(body), &answer); err != nil {
		return nil, fmt.Errorf("error reading response %v", err)
	}

	if answer.Status == server.StatusFail {
		return nil, fmt.Errorf("Error: Can not list available games, code:%d, %s", answer.Code, answer.Message)
	}

	return answer.Data.Content, nil
}
