package utils

import (
	"agentgo/pkg/conf"
	"gopkg.in/gomail.v2"
	"log"
)

const (
	CodeMsg     = "AgentGo Captcha Code below, valid for 2 minutes: "
	UserNameMsg = "AgentGo has generated a random username for you: "
)

func SendEmail(toEmail, code, msg string) error {
	c := conf.Config.Email

	m := gomail.NewMessage()
	m.SetHeader("From", c.User)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "System Notification from AgentGo")
	m.SetBody("text/plain", msg+"\n\n"+code)

	d := gomail.NewDialer(c.Host, c.Port, c.User, c.Password)

	// Send the email and log the result
	if err := d.DialAndSend(m); err != nil {
		log.Printf("[Email Service] Send email to %s failed: %v\n", toEmail, err)
		return err
	}
	
	log.Printf("[Email Service] Successfully sent email to %s\n", toEmail)
	return nil
}