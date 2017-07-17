package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/poo/bettybot/pkg/module"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(BuildCmd)
}

var BuildCmd = &cobra.Command{
	Use: "build",
	Run: buildCmd,
}

func buildCmd(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Please include JSON config")
	}

	jsn := args[0]

	b, err := ioutil.ReadFile(jsn)
	if err != nil {
		log.Fatalf("Error reading JSON config: %v", err)
	}

	files := module.Files{}
	if err := json.Unmarshal(b, &files); err != nil {
		log.Fatalf("Error decoding JSON config: %v", err)
	}

	output, err := files.Build()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}
