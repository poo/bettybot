package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "bettybot",
	Short: "Build Emails",
	Long:  `Build betty emails`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("foo")
	},
}
