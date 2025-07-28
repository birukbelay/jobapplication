package email

type EmailFields struct {
	To       string
	HtmlPath string
	Subject  string
	From     string
}

type EmailSender interface {
	SendEmailTmpl(fields EmailFields, templateStruct any) error
	SendEmail(to, subject string, templatePath EmailTemplates, templateStruct any) error
	//TODO: SendVerificationCode
}

type VerificationSender interface {
	SendVerificationCode(to string, code string) error
}

type EmailTemplates string

const (
	VerificationCodeTemplate = EmailTemplates("verification_code.html")
	UserInvitationTemplate   = EmailTemplates("user_invitation.html")
)

func (t EmailTemplates) S() string {
	return string(t)
}
