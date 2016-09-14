package mail

import "log"

var m Mailer

type Mailer interface {
	SendMail(emails []string, content string) error
}

func InitFakeMailer() {
	m = &FakeMailer{}
}

type FakeMailer struct {
}

func (m *FakeMailer) SendMail(emails []string, content string) error {
	log.Println("Send mail to %v: content %s", emails, content)
	return nil
}

func SendMail(emails []string, content string) {
	m.SendMail(emails, content)
}
