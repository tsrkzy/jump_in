package mail

import (
	"crypto/tls"
	"fmt"
	"github.com/tsrkzy/jump_in/helper"
	"log"
	"net"
	"net/mail"
	"net/smtp"
)

type Content struct {
	MailTo  string
	NameTo  string
	Subject string
	Body    string
}

// SendMailSSL https:gist.github.com/chrisgillis/10888032
func SendMailSSL(m *Content) error {
	domainName := helper.MustGetenv("SMTP_SERVER_NAME")
	port := helper.MustGetenv("SMTP_SERVER_PORT")
	mailerAddress := helper.MustGetenv("SMTP_SERVER_MAIL_ADDRESS")
	password := helper.MustGetenv("SMTP_SERVER_PASSWORD")
	mailerName := helper.MustGetenv("SMTP_SERVER_MAILER_NAME")

	from := mail.Address{Name: mailerName, Address: mailerAddress}
	to := mail.Address{Name: m.NameTo, Address: m.MailTo}

	subj := m.Subject
	body := m.Body
	servername := domainName + ":" + port
	username := mailerAddress

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	host, _, err := net.SplitHostPort(servername)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", username, password, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
		return err
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
		return err
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
		return err
	}

	err = c.Quit()
	if err != nil {
		return err

	}
	return nil
}

//func main() {
//	//sendMail()
//	m := MailContent{
//		MailTo:  "tsrmix+jump_in@gmail.com",
//		Subject: "test! 2022/10/27",
//		Body:    "body1\n2\n\n4",
//	}
//	SendMailSSL(m)
//}
