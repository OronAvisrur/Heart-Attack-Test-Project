package main

import (
	"bytes"
	"html/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

// This function send mail using SMTP
func (mail_obj *Mail) SendSMTPMessage(msg Message) error {

	// Fix empty important empty info if needed
	if msg.From == "" {
		msg.From = mail_obj.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = mail_obj.FromName
	}

	// Create map of the data
	data := map[string]any{
		"message": msg.Data,
	}

	// Add the data to the message DataMap
	msg.DataMap = data

	// Build HTML message format
	formattedMessage, possible_error := mail_obj.buildHTMLMessage(msg)
	if possible_error != nil {
		return possible_error
	}

	// Build plain text message format
	plainMessage, possible_error := mail_obj.buildPlainTextMessage(msg)
	if possible_error != nil {
		return possible_error
	}

	// Set up the SMTP client
	server := mail.NewSMTPClient()
	server.Host = mail_obj.Host
	server.Port = mail_obj.Port
	server.Username = mail_obj.Username
	server.Password = mail_obj.Password
	server.Encryption = mail.EncryptionNone
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// Connect to the client
	smtpClient, possible_error := server.Connect()
	if possible_error != nil {
		return possible_error
	}

	// Send the message
	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextAMP, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	// Check for error
	possible_error = email.Send(smtpClient)
	if possible_error != nil {
		return possible_error
	}

	return nil

}

// This function convert Message struct object to HTML
func (mail_obj *Mail) buildHTMLMessage(msg Message) (string, error) {
	// Get the premade templates
	templateToRender := "./templates/mail.html.gohtml"

	t, possible_error := template.New("email-html").ParseFiles(templateToRender)
	if possible_error != nil {
		return "", possible_error
	}

	var tpl bytes.Buffer

	if possible_error = t.ExecuteTemplate(&tpl, "body", msg.DataMap); possible_error != nil {
		return "", possible_error
	}

	formattedMessage := tpl.String()
	formattedMessage, possible_error = mail_obj.inlineCSS(formattedMessage)
	if possible_error != nil {
		return "", possible_error
	}

	return formattedMessage, nil
}

// This function convert Message struct object to plain text
func (mail_obj *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.plain.gohtml"

	t, possible_error := template.New("email-plain").ParseFiles(templateToRender)
	if possible_error != nil {
		return "", possible_error
	}

	var tpl bytes.Buffer

	if possible_error = t.ExecuteTemplate(&tpl, "body", msg.DataMap); possible_error != nil {
		return "", possible_error
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

func (mail_obj *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, possible_error := premailer.NewPremailerFromString(s, &options)
	if possible_error != nil {
		return "", possible_error
	}

	html, possible_error := prem.Transform()
	if possible_error != nil {
		return "", possible_error
	}

	return html, nil

}
