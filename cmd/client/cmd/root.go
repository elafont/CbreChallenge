package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/elafont/CbreChallenge/hangman"
	"github.com/spf13/cobra"
)

// var cfgFile string
const DEFAULTHOST = "localhost"
const DEFAULTWEBPORT = "8080"

// Commodity structs to make simple to unmarshal json responses
type responseHs struct { // Adapted from server.Response
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *struct {
		Type    string
		Content *hangman.Hstatus
	} `json:"data"`
}

type responseHsArr struct { // Adapted from server.Response
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *struct {
		Type    string
		Content *[]hangman.Hstatus
	} `json:"data"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "A client for the hangman server.",
	Long: `
This client connects with the hangman server and lets you play the hangman game. 
* The tipical scenario is using the "New" command to instruct the server to pick 
	a random word and start a new game.
* "List" and "Show" commands will let you see the status of any or all games started
	on the server.
* The "Guess" command let you guess one letter of the hidden word.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

var host string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&host, "server", "s", DEFAULTHOST+":"+DEFAULTWEBPORT, "Host:port of the hangman server, ie: hangame.com:8080")
}

func bindJSON(r io.Reader, target interface{}) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}
