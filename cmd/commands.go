package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// VERSION of the client app
var VERSION string

var rootCmd = &cobra.Command{
	Use:   "MailGun Bulk Sender",
	Short: "Send Test/Html bulk emails from CSV",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

// Execute start running the Dtable node
func Execute(version string) {
	VERSION = version

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
