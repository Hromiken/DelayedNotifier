package sender

import (
	"context"
	"crypto/tls"
	"delayedNotifier/cfg"
	"delayedNotifier/internal/entity"
	"fmt"
	"net/smtp"
)

type EmailSender struct {
	cfg cfg.EmailConfig
}

func NewEmailSender(cfg cfg.EmailConfig) *EmailSender {
	return &EmailSender{cfg: cfg}
}

func (s *EmailSender) Send(ctx context.Context, n entity.Notification) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)

	msg := []byte(
		"From: " + s.cfg.From + "\r\n" +
			"To: " + n.UserID + "\r\n" +
			"Subject: Notification\r\n\r\n" +
			n.Message + "\r\n",
	)

	tlsConfig := &tls.Config{
		ServerName:         s.cfg.Host,
		InsecureSkipVerify: true, // для gmail/mail.ru — работает
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, s.cfg.Host)
	if err != nil {
		return err
	}
	defer c.Close()

	if err = c.Auth(auth); err != nil {
		return err
	}

	if err = c.Mail(s.cfg.From); err != nil {
		return err
	}

	if err = c.Rcpt(n.UserID); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	return c.Quit()
}
