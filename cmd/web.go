package cmd

import (
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(WebCmd)
}

// WebCmd is used to run bettybot as a client/server
var WebCmd = &cobra.Command{
	Use: "web",
	Run: webCmd,
}

func webCmd(cmd *cobra.Command, args []string) {
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
