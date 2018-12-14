package mailgunsend

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/igorrendulic/mailgun-csv-bulk-send/mustache"

	mailgun "github.com/mailgun/mailgun-go"
)

// CSV Config with send methods
type CSV struct {
	Subject      string `json:"subject"`
	From         string `json:"from"`
	DomainName   string `json:"domain_name"`
	APIKey       string `json:"api_key"`
	Test         bool   `json:"test"`
	HTMLTemplate string `json:"html_template"`
	TEXTTemplate string `json:"text_template,omitempty"`
}

// ReadCSVAndSend reads csv file line by line, augments with mustache parameters and send an email
func (c *CSV) ReadCSVAndSend(filename string) error {

	err := c.validateConfig()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}
	defer f.Close()

	mg := mailgun.NewMailgun(c.DomainName, c.APIKey)

	readChannel := c.csvReader(f)

	for lineJSON := range readChannel {
		html, text, recipient, err := c.mustacheIt(lineJSON)
		if err != nil {
			log.Fatal(err)
			continue
		}
		message := mg.NewMessage(c.From, c.Subject, text, recipient)
		message.SetReplyTo(c.From)
		message.SetHtml(html)
		message.SetTracking(true)

		if !c.Test {
			resp, id, err := mg.Send(message)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("ID: %s Resp: %s\n", id, resp)
		} else {
			fmt.Printf("Test success for %s\n", recipient)
		}

	}
	return nil
}

// mustacheIt takes html/text file and replaces placeholders with text from csv
func (c *CSV) mustacheIt(lineJSON map[string]interface{}) (string, string, string, error) {
	html := mustache.RenderFile(c.HTMLTemplate, lineJSON)
	var text string
	if c.TEXTTemplate != "" {
		text = mustache.RenderFile(c.TEXTTemplate, lineJSON)
	}

	return html, text, lineJSON["email"].(string), nil
}

func (c *CSV) csvReader(rc io.Reader) (ch chan map[string]interface{}) {
	ch = make(chan map[string]interface{}, 10)

	go func() {
		r := csv.NewReader(rc)
		r.LazyQuotes = true
		r.Comma = ','
		r.Comment = '#'

		header, err := r.Read()
		if err != nil { //read header
			log.Fatal(err)
			return
		}

		headerMap := make(map[int]string)
		for i, h := range header {
			headerMap[i] = strings.Trim(h, " \t")
		}

		if err != nil {
			log.Fatal(err)
			return
		}

		defer close(ch)
		for {
			rec, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					log.Fatal(err)
				}
			}
			line := make(map[string]interface{})
			for index, record := range rec {
				name := headerMap[index]
				line[name] = strings.Trim(record, " \t")
			}
			ch <- line
		}
	}()
	return
}

func (c *CSV) validateConfig() error {
	if c.APIKey == "" {
		return errors.New("mailgun api key missing")
	}
	if c.DomainName == "" {
		return errors.New("mailgun domain name missing")
	}
	if c.From == "" {
		return errors.New("from is missing")
	}
	if c.Subject == "" {
		return errors.New("subject is missing")
	}
	if c.HTMLTemplate == "" {
		return errors.New("Email template required")
	}
	if _, err := os.Stat(c.HTMLTemplate); os.IsNotExist(err) {
		return fmt.Errorf("csv does not exist: %s", c.HTMLTemplate)
	}
	if c.TEXTTemplate != "" {
		if _, err := os.Stat(c.TEXTTemplate); os.IsNotExist(err) {
			return fmt.Errorf("csv does not exist: %s", c.TEXTTemplate)
		}
	}

	return nil
}
