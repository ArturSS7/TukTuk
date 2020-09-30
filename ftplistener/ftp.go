package ftplistener

import (
	"TukTuk/discordbot"
	"TukTuk/emailalert"
	"TukTuk/telegrambot"
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	status426 = "426 Bye."
	status220 = "220 TukTuk callback server."
	status331 = "331 password required."
)

var FTPServer *Server

type Server struct {
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
}

func StartFTP(db *sql.DB) error {
	s, err := NewServer(db)
	if err != nil {
		return err
	}
	FTPServer = s
	return nil
}

func NewServer(db *sql.DB) (*Server, error) {
	s := &Server{
		quit: make(chan interface{}),
	}
	l, err := net.Listen("tcp", ":21")
	if err != nil {
		return nil, err
	}
	s.listener = l
	s.wg.Add(1)
	go s.serve(db)
	return s, nil
}

func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
}

func (s *Server) serve(db *sql.DB) {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Println("accept error", err)
			}
		} else {
			s.wg.Add(1)
			go func() {
				s.handleFTP(newConn(conn, db))
				s.wg.Done()
			}()
		}
	}
}

type Conn struct {
	conn net.Conn
	data *bytes.Buffer
	db   *sql.DB
}

//return new connection with parameters
func newConn(conn net.Conn, db *sql.DB) *Conn {
	return &Conn{
		conn: conn,
		data: new(bytes.Buffer),
		db:   db,
	}
}

//logging to database
func (c *Conn) log() {
	var lastInsertId int64 = 0
	err := c.db.QueryRow("insert into ftp (data, source_ip, time) values ($1, $2, $3) RETURNING id", c.data.String(), c.conn.RemoteAddr().String(), time.Now().String()).Scan(&lastInsertId)
	if err != nil {
		log.Println(err)
	}

	//Send alert to Telegram
	telegrambot.BotSendAlert(c.data.String(), c.conn.RemoteAddr().String(), time.Now().String(), "FTP", lastInsertId)
	//Send alert to email
	emailalert.SendEmailAlert("FTP Alert", "Remoute Address: "+c.conn.RemoteAddr().String()+"\n+"+c.data.String()+"\n"+time.Now().String())
	//Send alert to Discord
	discordbot.BotSendAlert(c.data.String(), c.conn.RemoteAddr().String(), time.Now().String(), "FTP", lastInsertId)

}

func (c *Conn) respond(s string) {
	fmt.Fprintf(c.data, ">>%s\n", s)
	_, err := fmt.Fprint(c.conn, s, "\n")
	if err != nil {
		log.Print(err)
	}
}

//main request handler
//if input differs from USER we just drop the connection and don't log it
func (s *Server) handleFTP(c *Conn) {
	c.respond(status220)
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		input := strings.Fields(scanner.Text())
		if len(input) == 0 {
			continue
		}
		command, args := input[0], input[1:]
		fmt.Fprintf(c.data, "<< %s %v\n", command, args)
		switch command {
		case "USER":
			c.respond(status331)
		case "PASS":
			c.respond(status426)
			c.log()
			return
		default:
			c.respond(status426)
			return
		}
	}
	if scanner.Err() != nil {
		log.Print(scanner.Err())
	}
}
