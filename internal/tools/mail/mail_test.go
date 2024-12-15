package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSend1(t *testing.T) {
	cfg := ConfigMailer{
		User:       "noreply@clive.tk",
		Password:   "U74f7Bea1NM5",
		MailBox:    "noreply@clive.tk",
		ConnectURL: "mail.clive.tk:25", // 25
	}
	//require.True(t, len(cfg.Password) == 0)
	mailer := NewMailer(cfg)
	require.NoError(t, mailer.Send("rerednaw69@gmail.com", "spam", "spam message", "admin"))
}

func TestSend2(t *testing.T) {
	cfg := ConfigMailer{
		User:       "noreply@clive.tk",
		Password:   "U74f7Bea1NM5",
		MailBox:    "noreply@clive.tk",
		ConnectURL: "mail.clive.tk:465", // 25
		IsUseTLS:   true,
	}
	//require.True(t, len(cfg.Password) == 0)
	mailer := NewMailer(cfg)
	require.NoError(t, mailer.Send("rerednaw69@gmail.com", "spam", "spam message", "admin"))
}

func TestSendFail(t *testing.T) {
	cfg := ConfigMailer{
		User:       "noreply@clive.tk",
		Password:   "",
		MailBox:    "noreply@clive.tk",
		ConnectURL: "mail.clive.tk:25",
	}
	require.True(t, len(cfg.Password) == 0)
	mailer := NewMailer(cfg)
	require.ErrorContains(t, mailer.Send("rerednaw69@gmail.com", "spam", "spam message", "admin"), "authentication failed")
}
