// Package smtpsender sends real email via a configured SMTP server. It's
// wired in by main.go only when SMTP_HOST is set; otherwise the app falls
// back to logsender.
package smtpsender

import (
	"context"
	"fmt"
	"net/smtp"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type Sender struct {
	cfg Config
}

func New(cfg Config) *Sender {
	return &Sender{cfg: cfg}
}

func (s *Sender) Send(_ context.Context, to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)

	var auth smtp.Auth
	if s.cfg.Username != "" {
		auth = smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	}

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.cfg.From, to, subject, body,
	)

	if err := smtp.SendMail(addr, auth, s.cfg.From, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("sending email via smtp: %w", err)
	}
	return nil
}
