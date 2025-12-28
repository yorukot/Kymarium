package config

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

var (
	smtpAuth smtp.Auth
	initOnce sync.Once
	initErr  error
)

func InitSMTP() error {
	initOnce.Do(func() {
		initErr = initSMTPInternal()
	})
	return initErr
}

func initSMTPInternal() error {
	if !Env().SMTPEnabled {
		zap.L().Info("SMTP is disabled in the configuration.")
		return nil
	}

	host := strings.TrimSpace(Env().SMTPHost)
	port := strings.TrimSpace(Env().SMTPPort)
	user := strings.TrimSpace(Env().SMTPUsername)
	pass := Env().SMTPPassword

	if host == "" {
		return fmt.Errorf("missing SMTP host")
	}
	if port == "" {
		return fmt.Errorf("missing SMTP port")
	}

	if user != "" {
		if strings.TrimSpace(pass) == "" {
			return fmt.Errorf("missing SMTP password")
		}
		smtpAuth = smtp.PlainAuth("", user, pass, host)
	} else {
		smtpAuth = nil
	}

	zap.L().Info("SMTP initialized",
		zap.String("host", host),
		zap.String("port", port),
		zap.Bool("auth_enabled", smtpAuth != nil),
	)
	return nil
}

func SendEmail(to string, cc, bcc []string, subject string, body string) error {
	if err := InitSMTP(); err != nil {
		return err
	}

	if !Env().SMTPEnabled {
		zap.L().Warn("Attempted to send email while SMTP is disabled.")
		return nil
	}

	to = strings.TrimSpace(to)
	if to == "" {
		return fmt.Errorf("missing recipient: to")
	}

	host := strings.TrimSpace(Env().SMTPHost)
	port := strings.TrimSpace(Env().SMTPPort)
	from := strings.TrimSpace(Env().SMTPFrom)
	if from == "" {
		return fmt.Errorf("missing SMTPFrom")
	}

	recipients := uniqueNonEmpty(append([]string{to}, append(cc, bcc...)...))
	if len(recipients) == 0 {
		return fmt.Errorf("no recipients")
	}

	msg, ccClean := buildPlainTextMessage(from, to, cc, subject, body)

	if err := sendSMTP(host, port, smtpAuth, from, recipients, msg); err != nil {
		zap.L().Error("Failed to send email", zap.Error(err))
		return err
	}

	zap.L().Info("Email sent successfully",
		zap.String("to", to),
		zap.Strings("cc", ccClean),
		zap.Int("bcc_count", len(uniqueNonEmpty(bcc))),
	)
	return nil
}

func buildPlainTextMessage(from, to string, cc []string, subject, body string) ([]byte, []string) {
	var buf bytes.Buffer

	writeHeader(&buf, "From", from)
	writeHeader(&buf, "To", to)

	ccClean := uniqueNonEmpty(cc)
	if len(ccClean) > 0 {
		writeHeader(&buf, "Cc", strings.Join(ccClean, ", "))
	}

	writeHeader(&buf, "Subject", mime.QEncoding.Encode("utf-8", subject))
	writeHeader(&buf, "Date", time.Now().Format(time.RFC1123Z))
	writeHeader(&buf, "MIME-Version", "1.0")
	writeHeader(&buf, "Content-Type", `text/plain; charset="UTF-8"`)
	writeHeader(&buf, "Content-Transfer-Encoding", "quoted-printable")
	buf.WriteString("\r\n")

	qp := quotedprintable.NewWriter(&buf)
	_, _ = qp.Write([]byte(body))
	_ = qp.Close()
	buf.WriteString("\r\n")

	return buf.Bytes(), ccClean
}

func sendSMTP(host, port string, auth smtp.Auth, from string, recipients []string, msg []byte) error {
	addr := net.JoinHostPort(host, port)

	if portNum, _ := strconv.Atoi(port); portNum == 465 {
		tlsConn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host})
		if err != nil {
			return fmt.Errorf("tls dial failed: %w", err)
		}
		defer tlsConn.Close()

		c, err := smtp.NewClient(tlsConn, host)
		if err != nil {
			return fmt.Errorf("new smtp client failed: %w", err)
		}
		defer c.Quit()

		if err := maybeAuth(c, auth); err != nil {
			return err
		}
		return sendData(c, from, recipients, msg)
	}

	c, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("smtp dial failed: %w", err)
	}
	defer c.Close()

	starttlsOK, _ := c.Extension("STARTTLS")
	if starttlsOK {
		if err := c.StartTLS(&tls.Config{ServerName: host}); err != nil {
			return fmt.Errorf("starttls failed: %w", err)
		}
		_, _ = c.Extension("8BITMIME")
	} else {
		if auth != nil {
			return fmt.Errorf("server does not support STARTTLS; refusing to auth over insecure connection")
		}
		zap.L().Warn("SMTP server does not advertise STARTTLS; sending without TLS")
	}

	if err := maybeAuth(c, auth); err != nil {
		return err
	}

	if err := sendData(c, from, recipients, msg); err != nil {
		return err
	}

	return c.Quit()
}

func maybeAuth(c *smtp.Client, auth smtp.Auth) error {
	if auth == nil {
		return nil
	}
	if ok, _ := c.Extension("AUTH"); ok {
		if err := c.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth failed: %w", err)
		}
	}
	return nil
}

func sendData(c *smtp.Client, from string, recipients []string, msg []byte) error {
	if err := c.Mail(from); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}
	for _, rcpt := range recipients {
		if err := c.Rcpt(rcpt); err != nil {
			return fmt.Errorf("RCPT TO failed (%s): %w", rcpt, err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("DATA failed: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		_ = w.Close()
		return fmt.Errorf("write message failed: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("close DATA failed: %w", err)
	}
	return nil
}

func writeHeader(buf *bytes.Buffer, key, value string) {
	value = strings.ReplaceAll(value, "\r", "")
	value = strings.ReplaceAll(value, "\n", "")
	buf.WriteString(key)
	buf.WriteString(": ")
	buf.WriteString(value)
	buf.WriteString("\r\n")
}

func uniqueNonEmpty(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
