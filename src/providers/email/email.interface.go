package email

type EmailFields struct {
	To       string
	HtmlPath string
	Subject  string
	From     string
}

type EmailSender interface {
	SendEmail(fields EmailFields, templateStruct any) error
	//TODO: SendVerificationCode
}

type VerificationSender interface {
	SendVerificationCode(to string, code string) error
}
