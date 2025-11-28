// Package email provides email operation related utilities
// Author: Done-0
// Created: 2025-08-24
package email

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/Done-0/gin-scaffold/configs"
)

// Email server configuration
var emailServers = map[string]struct {
	Server string
	Port   int
	SSL    bool
}{
	"qq":      {"smtp.qq.com", 465, true},         // QQ email uses SSL encryption
	"gmail":   {"smtp.gmail.com", 465, true},      // Gmail uses SSL encryption
	"outlook": {"smtp.office365.com", 587, false}, // Outlook uses TLS encryption
}

// SendEmail sends email to specified email addresses with specified content type
func SendEmail(subject, content string, toEmails []string, contentType string) (bool, error) {
	cfgs, err := configs.GetConfig()
	if err != nil {
		return false, fmt.Errorf("failed to load email config: %v", err)
	}

	// Get SMTP configuration
	emailType := cfgs.AppConfig.Email.EmailType
	serverConfig := emailServers[emailType]

	// Create email
	m := gomail.NewMessage()
	m.SetHeader("From", cfgs.AppConfig.Email.FromEmail)
	m.SetHeader("To", toEmails...)
	m.SetHeader("Subject", subject)
	m.SetBody(contentType, content)

	// Configure sender
	d := gomail.NewDialer(
		serverConfig.Server,
		serverConfig.Port,
		cfgs.AppConfig.Email.FromEmail,
		cfgs.AppConfig.Email.EmailSmtp,
	)

	// Configure security options based on port
	if serverConfig.SSL {
		d.SSL = true
	} else {
		d.TLSConfig = &tls.Config{
			ServerName: serverConfig.Server,
			MinVersion: tls.VersionTLS12,
		}
	}

	if err := d.DialAndSend(m); err != nil {
		return false, fmt.Errorf("failed to send email: %v", err)
	}

	return true, nil
}

// NewRand generates six-digit random verification code
func NewRand() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(900000) + 100000
}
