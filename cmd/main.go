package main

import (
	"flag"
	"fmt"

	o365 "github.com/alex-arce/O365Sender"
	"github.com/BurntSushi/toml"
)

var configFilePath string

var mailServerCfg *o365.MailServerConfig

var from, to, subject, text string

//var filenameToAttach string

func main() {
	fmt.Println("[O365 Mail sender]")

	flag.StringVar(&configFilePath, "config", "o365sender.conf", "config location")

	flag.StringVar(&from, "from", "", "FROM address")
	flag.StringVar(&to, "to", "", "TO address")
	flag.StringVar(&subject, "subject", "", "Subject")
	flag.StringVar(&text, "text", "", "Mail text")
	//flag.StringVar(&filenameToAttach, "filename", "", "filename to add")

	flag.Parse()

	if _, err := toml.DecodeFile(configFilePath, &mailServerCfg); err != nil {
		fmt.Println(err)
		return
	}

	o365Client := o365.NewMailClient(mailServerCfg)

	o365Client.AddFrom(from)
	o365Client.AddTo(to)
	o365Client.AddSubject(subject)
	o365Client.AddBody(text)

	sendError := o365Client.Send()
	if sendError != nil {
		fmt.Errorf("[-] ERROR sending mail :-(")
	}

	fmt.Println("[+] Mail Sent!!")
}
