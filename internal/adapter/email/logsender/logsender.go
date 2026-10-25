// Package logsender is the default EmailSender used when no SMTP server is
// configured (SMTP_HOST unset). It logs the email instead of delivering it,
// so password-reset/email-verification flows are fully exercisable — and
// their tokens visible for manual testing — without any mail setup.
package logsender

import (
	"context"
	"log/slog"
)

type Sender struct{}

func New() *Sender {
	return &Sender{}
}

func (s *Sender) Send(_ context.Context, to, subject, body string) error {
	slog.Info("email not delivered — no SMTP configured, logging instead",
		"to", to, "subject", subject, "body", body)
	return nil
}
