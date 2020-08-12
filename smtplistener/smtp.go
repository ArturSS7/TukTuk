package smtplistener

import (
	"TukTuk/telegrambot"
	"database/sql"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/emersion/go-smtp"
)

// The Backend implements SMTP server methods.
type Backend struct{}

// Login handles a login command with username and password.
func (bkd *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if username != "username" || password != "password" {
		return nil, errors.New("Invalid username or password")
	}
	return &Session{}, nil
}

// AnonymousLogin requires clients to authenticate using SMTP AUTH before sending emails
func (bkd *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return nil, smtp.ErrAuthRequired
}

// A Session is returned after successful login.
type Session struct{}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	log.Println("Mail from:", from)
	return nil
}

func (s *Session) Rcpt(to string) error {
	log.Println("Rcpt to:", to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if b, err := ioutil.ReadAll(r); err != nil {
		return err
	} else {
		log.Println("Data:", string(b))
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func StartSMTP(db *sql.DB, Domain string) {
	be := &Backend{}

	s := smtp.NewServer(be)

	s.Addr = ":587"
	s.Domain = Domain
	s.ReadTimeout = 100 * time.Second
	s.WriteTimeout = 100 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true
	log.Println("Starting server at", s.Addr)
	err, RemoteAddr := s.ListenAndServe()
	logSMTP(db, RemoteAddr, s.Domain)
	if err != nil {
		log.Fatal(err)
	}
}

func logSMTP(db *sql.DB, RemoteAddr, Domain string) {
	_, err := db.Exec("insert into smtp (data, source_ip, time) values ($1, $2, $3)", Domain, RemoteAddr, time.Now().String())

	if err != nil {
		log.Fatal(err)
	}
	telegrambot.BotSendAlert("", RemoteAddr, time.Now().String(), "SMTP")
}
