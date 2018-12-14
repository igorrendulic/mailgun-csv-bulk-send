package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/igorrendulic/mailgun-csv-bulk-send/mailgunsend"
	"github.com/spf13/cobra"
)

var (
	config  string
	csvFile string

	sendCmd = &cobra.Command{
		Use:   "send",
		Short: "Bulk Send Emails",
		Long:  ``,
		Run:   send,
	}
)

func init() {

	sendCmd.Flags().StringVarP(&config, "config", "c", "", "Config file. See example in github.com/igorrendulic/mailgun-csv-bulk-send")
	sendCmd.Flags().StringVarP(&csvFile, "csv", "s", "", "Csv File. Make sure first line of CSV has headers")

	rootCmd.AddCommand(sendCmd)
}

func send(ccmd *cobra.Command, args []string) {
	err := validateSend(csvFile, config)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(1)
	}

	confJSON, err := readConfig(config)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(1)
	}

	err = confJSON.ReadCSVAndSend(csvFile)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(1)
	}

}

func readConfig(configFilePath string) (*mailgunsend.CSV, error) {

	conf, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	var dat mailgunsend.CSV

	if err := json.Unmarshal(conf, &dat); err != nil {
		return nil, err
	}
	return &dat, nil
}

func validateSend(csvFile, config string) error {
	if csvFile == "" {
		return errors.New("CSV file required")
	}
	if config == "" {
		return errors.New("mailgun config file missing")
	}
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		return fmt.Errorf("csv does not exist: %s", csvFile)
	}
	if _, err := os.Stat(config); os.IsNotExist(err) {
		return fmt.Errorf("mailgun config does not exist: %s", config)
	}
	return nil
}
