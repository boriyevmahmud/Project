package main

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/mail.v2"
)

func main() {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "boriyevmahmud@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", "umarovbaxrom10@gmail.com")
	
	// Set E-Mail subject
	m.SetHeader("code:", "salom oka qaliysiz")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", "salom oka qaliysiz")

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "boriyevmahmud@gmail.com", "mahmUd3253")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	return

}
