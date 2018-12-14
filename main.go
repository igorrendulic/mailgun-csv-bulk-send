package main

import "github.com/igorrendulic/mailgun-csv-bulk-send/cmd"

var (
	version = "0.0.1"
)

func main() {
	cmd.Execute(version)
}
