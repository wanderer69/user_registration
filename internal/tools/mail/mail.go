package mail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
)

type ConfigMailer struct {
	User       string `envconfig:"MAILER_USER_NAME" required:"true"`
	Password   string `envconfig:"MAILER_USER_PASSWORD" required:"true"`
	MailBox    string `envconfig:"MAILER_MAIL_BOX" required:"true"`
	ConnectURL string `envconfig:"MAILER_CONNECT_URL" required:"true"`
	IsUseTLS   bool   `envconfig:"MAILER_IS_USE_TLS" default:"false"`
}

type Mailer struct {
	config ConfigMailer
}

func NewMailer(cfg ConfigMailer) *Mailer {
	return &Mailer{
		config: cfg,
	}
}

func (m *Mailer) Send(email string, subject string, messageBody string, fromName string) error {
	from := mail.Address{
		Name:    fromName,
		Address: m.config.MailBox,
	}
	to := mail.Address{
		Name:    "",
		Address: email,
	}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + messageBody

	host, _, _ := net.SplitHostPort(m.config.ConnectURL)

	auth := smtp.PlainAuth("", m.config.User, m.config.Password, host)

	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	var c *smtp.Client

	if m.config.IsUseTLS {
		conn, err := tls.Dial("tcp", m.config.ConnectURL, tlsConfig)
		if err != nil {
			//Handle error
		}

		c, err = smtp.NewClient(conn, host)
		if err != nil {
			//Handle error
		}
	} else {
		var err error
		c, err = smtp.Dial(m.config.ConnectURL)
		if err != nil {
			return err
		}
	}

	c.StartTLS(tlsConfig)

	// Auth
	if err := c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err := c.Mail(from.Address); err != nil {
		return err
	}

	if err := c.Rcpt(to.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}
