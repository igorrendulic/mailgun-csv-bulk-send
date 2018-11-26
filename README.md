# Command line Personalized Email Sending with Mailgun from CSV file

![alt text](https://raw.githubusercontent.com/igorrendulic/mailgun-csv-bulk-send/master/docs/mailgun-bulk.png)

This simple program can handle CSV files for personalized email sending using [Mailgun API](https://www.mailgun.com). It uses [template language Mustache](http://mustache.github.io/mustache.5.html) for both HTML and Plain text templates. 

The reason for this command line interface is fast preparation of Email Template Contents with exported CSV list.

# Sending email via Command Line tool

Edit the config file under `data/config.json`, prepare `domain name and api key from Mailgun`, change `subject` and `from`. 
```
{
	"domain_name": "domainname",
	"api_key": "appkey",
	"subject": "this is my subject",
	"from": "mymail <mymail@mymail.com>",
	"html_template":"templates/simple-template.html",
	"text_template":"templates/simple-template.txt",
	"test":true
}
```

Modify Exampl Excel file under `data/example.csv`. 

- **CSV file must have headers as first line**
- **There must be at least 1 field named: `email` in the headers**
- **Other header names are used as Mustache placeholders in the template**

```
email,name,date
myemail@mymail.com,igor,3/12/18
myemail@mymail.com,Ashley,3/12/18
myemail@mymail.com,igor2,3/12/18
```

Run: 
```
go run main.go send -c data/config.json -s data/example.csv 
```

-c is a path to config.json file

-s is path to CSV file

# Customizing template

Templates are stored under folder `templates/...`

Template customization is based on Mustache templating language. [You can read more about how to customize templates here](http://mustache.github.io/mustache.5.html).

Plain/Text template example
```
Hi {{name}},

Sometimes you just want to send a simple HTML email with a simple design and clear call to action. This is it.

This is a really simple email template. Its sole purpose is to get the recipient to click the button with no distractions.

Good luck! The date is {{date}}.
```