package smtp

import (
	"bytes"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/projTemplate/goauth/src/providers/email"
	"github.com/projTemplate/goauth/src/providers/email/templates"
)

// VerificationCodeData holds the data for the email template
type VerificationCodeData struct {
	Name string
	Code string
}

type SmtpSender struct {
	Host     string
	Port     string
	Password string
	From     string
}

func NewSmtpSender(host, port, pwd, from string) email.EmailSender {
	return &SmtpSender{
		Host:     host,
		Port:     port,
		Password: pwd,
		From:     from,
	}
}

func NewVerificationCodeSender(host, port, pwd, from string) email.VerificationSender {
	return &SmtpSender{
		Host:     host,
		Port:     port,
		Password: pwd,
		From:     from,
	}
}

// SendEmail sends an email using Gmail SMTP
func (h *SmtpSender) SendEmail(fields email.EmailFields, templateStruct any) error {
	// Parse the email template
	body, err := ParseEmailTemplate(fields.HtmlPath, templateStruct)
	if err != nil {
		return err
	}
	return h.SendEmailT(fields, body)
}

// SendEmail sends an email using Gmail SMTP
func (h *SmtpSender) SendEmailT(fields email.EmailFields, body string) error {

	// Set up SMTP server configuration
	// smtpHost := "smtp.gmail.com"
	// smtpPort := "587"
	auth := smtp.PlainAuth("", fields.From, h.Password, h.Host)

	// Create the email message
	headers := make(map[string]string)
	headers["From"] = fields.From
	headers["To"] = fields.To
	headers["Subject"] = fields.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	// Build the message
	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(k + ": " + v + "\r\n")
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	// Send the email
	err := smtp.SendMail(
		h.Host+":"+h.Port,
		auth,
		fields.From,
		[]string{fields.To},
		msg.Bytes(),
	)
	return err
}

// SendEmail sends an email using Gmail SMTP
func (h *SmtpSender) SendVerificationCodeEmbeded(to string, code string) error {

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(wd, "src/providers/email/templates/verification_code.html")
	emailFields := email.EmailFields{
		To:       to,
		From:     h.From,
		Subject:  "Verify Your Email",
		HtmlPath: fullPath,
	}
	verificationData := VerificationCodeData{
		Code: code,
		Name: "Mr",
	}
	return h.SendEmail(emailFields, verificationData)
}

// SendEmail sends an email using Gmail SMTP
func (h *SmtpSender) SendVerificationCode(to string, code string) error {

	f, err := templates.Embedded.Open("verification_code.html")
	if err != nil {
		panic(err)
	}
	verificationData := VerificationCodeData{
		Code: code,
		Name: "Mr",
	}
	body, err := ParseEmbededTemplate(f, verificationData)
	if err != nil {
		return err
	}

	emailFields := email.EmailFields{
		To:      to,
		From:    h.From,
		Subject: "Verify Your Email",
	}

	return h.SendEmailT(emailFields, body)
}
