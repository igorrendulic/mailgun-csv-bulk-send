package main

import "github.com/mailgun-cmd-bulk-sender/cmd"

var (
	version = "0.0.1"
)

func main() {
	cmd.Execute(version)
}
