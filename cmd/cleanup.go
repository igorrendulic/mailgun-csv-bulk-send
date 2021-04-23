package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/igorrendulic/mailgun-csv-bulk-send/mailgunsend"
	"github.com/spf13/cobra"
)

var (
	cleanedCsv  string
	originalCsv string

	cleanupCmd = &cobra.Command{
		Use:   "cleanup",
		Short: "cleanup validasted emails",
		Long:  `Cleanup validated email addresses from bulk csv mailgun validation`,
		Run:   cleanup,
	}
)

func init() {

	cleanupCmd.Flags().StringVarP(&originalCsv, "original", "o", "", "Csv File. Make sure first line of CSV has headers")
	cleanupCmd.Flags().StringVarP(&cleanedCsv, "cleaned", "c", "", "Mailgun bulk CSV cleaned up/validate email addresses")

	rootCmd.AddCommand(cleanupCmd)
}

func cleanup(ccmd *cobra.Command, args []string) {
	err := validateCleanup(originalCsv, cleanedCsv)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		os.Exit(1)
	}

	f, err := os.Open(originalCsv)
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}
	defer f.Close()

	csvReader := mailgunsend.NewCSVReader()
	originalCsvMap := csvReader.LoadCSV(originalCsv)
	cleanedCsvMap := csvReader.LoadCSV(cleanedCsv)

	cleanedUpResultCsv := [][]string{{"email", "alias", "referralCode"}}

	for email, lineJSON := range originalCsvMap {
		if val, ok := cleanedCsvMap[email]; ok {
			risk := val.(map[string]interface{})["risk"]
			if risk == "low" {
				// add to new list lineJSON
				line := lineJSON.(map[string]interface{})
				e := line["email"].(string)
				alias := line["alias"].(string)
				refCode := line["referralCode"].(string)
				cleanupLine := []string{e, alias, refCode}
				cleanedUpResultCsv = append(cleanedUpResultCsv, cleanupLine)
				fmt.Printf("adding low risk email: %v\n", email)
			} else {
				fmt.Printf("skipping non low risk email: %v\n", email)
			}
		} else {
			fmt.Printf("email %v\n not found in validation list. Skipping", email)
		}
	}

	outputFile := strings.Split(originalCsv, ".")[0] + "_cleanedup.csv"
	fmt.Printf("creating a cleaned up CSV file with %d emails. Original size: %d. Saving to: %v\n", len(cleanedUpResultCsv), len(originalCsvMap), outputFile)

	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	w := csv.NewWriter(file)
	w.WriteAll(cleanedUpResultCsv)
	w.Flush()
	file.Close()
}

func validateCleanup(originalCsv, cleanedCsv string) error {
	if originalCsv == "" {
		return errors.New("Original CSV file required")
	}
	if cleanedCsv == "" {
		return errors.New("mailgun validated CSV file required")
	}
	if _, err := os.Stat(originalCsv); os.IsNotExist(err) {
		return fmt.Errorf("csv does not exist: %s", csvFile)
	}
	if _, err := os.Stat(cleanedCsv); os.IsNotExist(err) {
		return fmt.Errorf("mailgun config does not exist: %s", config)
	}
	return nil
}
