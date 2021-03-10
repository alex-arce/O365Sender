package o365

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
)

type MailServerConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}

func NewMailServerConfig(username, password, host, port string) *MailServerConfig {
	return &MailServerConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

type Attachment struct {
	filename string
	data     []byte
}

type O365Client struct {
	mailServerConfig *MailServerConfig
	fromAddress      string
	toAddress        []string
	subject          string
	body             []byte
	attachments      []Attachment
}

func NewMailClient(cfgMail *MailServerConfig) *O365Client {
	return &O365Client{
		mailServerConfig: cfgMail,
		fromAddress:      "",
		toAddress:        nil,
		subject:          "",
		body:             nil,
		attachments:      nil,
	}
}

//-------------------------------------------------

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}

//-------------------------------------------------

func (oc *O365Client) AddFrom(fromAddress string) {
	oc.fromAddress = fromAddress
}

func (oc *O365Client) AddTo(toAdd string) {
	oc.toAddress = append(oc.toAddress, toAdd)
}

func (oc *O365Client) AddSubject(subject string) {
	oc.subject = subject
}

func (oc *O365Client) AddBody(content string) {
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n\n", oc.subject, mimeHeaders)))
	body.Write([]byte(content))
	oc.body = body.Bytes()
}

func (oc *O365Client) AddAttachment(filename string) {
	// TODO
}

func (oc *O365Client) Send() error {
	server := fmt.Sprintf("%s:%s", oc.mailServerConfig.Host, oc.mailServerConfig.Port)

	conn, err := net.Dial("tcp", server)
	if err != nil {
		println(err)
		return err
	}

	c, err := smtp.NewClient(conn, oc.mailServerConfig.Host)
	if err != nil {
		println(err)
		return err
	}

	tlsconfig := &tls.Config{
		ServerName: oc.mailServerConfig.Host,
	}

	if err = c.StartTLS(tlsconfig); err != nil {
		println(err)
		return err
	}

	auth := LoginAuth(oc.mailServerConfig.Username, oc.mailServerConfig.Password)
	if err = c.Auth(auth); err != nil {
		println(err)
		return err
	}

	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n\n", oc.subject, mimeHeaders)))
	err = smtp.SendMail(oc.mailServerConfig.Host+":"+oc.mailServerConfig.Port, auth, oc.fromAddress, oc.toAddress, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
