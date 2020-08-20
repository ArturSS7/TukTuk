package plaintcplistener

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"html"
	"log"
	"net"
	"sync"
	"time"
)

var TCPServers = make(map[string]*Server, 0)

type Server struct {
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
}

func StartTCP(db *sql.DB, message string, port string) error {
	s, err := NewServer(db, message, port)
	if err != nil {
		return err
	}
	TCPServers[port] = s
	return nil
}

func NewServer(db *sql.DB, message string, port string) (*Server, error) {
	s := &Server{
		quit: make(chan interface{}),
	}
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}
	s.listener = l
	s.wg.Add(1)
	go s.serve(db, message, port)
	return s, nil
}

func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
}

func (s *Server) serve(db *sql.DB, message string, port string) {
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
				s.handleTCP(newConn(conn, db, message, port))
				s.wg.Done()
			}()
		}
	}
}

/*
func StartTCP(db *sql.DB, message string, port string) error{
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println(err)
		return err
	}
	TCPServers[port] = listener
	for {
		conn, err := listener.Accept()
		log.Printf("[*] Connection Accepted from [%s]\n", conn.RemoteAddr().String())
		if err != nil {
			log.Print(err)
			continue
		}
		go handleTCP(conn, db, message, port)
	}
}

*/

func (s *Server) handleTCP(c *Conn) {
	defer c.log()
	defer c.conn.Close()
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}
		fmt.Fprintf(c.data, scanner.Text())
		c.respond(c.message)
	}
	if scanner.Err() != nil {
		log.Print(scanner.Err())
	}
}

type Conn struct {
	conn    net.Conn
	data    *bytes.Buffer
	db      *sql.DB
	message string
	port    string
}

//return new connection with parameters
func newConn(conn net.Conn, db *sql.DB, message string, port string) *Conn {
	return &Conn{
		conn:    conn,
		data:    new(bytes.Buffer),
		db:      db,
		message: message,
		port:    port,
	}
}

//logging to database
func (c *Conn) log() {
	var lastInsertId int64 = 0
	err := c.db.QueryRow("insert into tcp (data, source_ip, time, port) values ($1, $2, $3, $4) RETURNING id", html.EscapeString(c.data.String()), c.conn.RemoteAddr().String(), time.Now().String(), c.port).Scan(&lastInsertId)
	if err != nil {
		log.Println(err)
	}

	//Send Alert to telegram
	//telegrambot.BotSendAlert(c.data.String(), c.conn.RemoteAddr().String(), time.Now().String(), "FTP", lastInsertId)
}

func (c *Conn) respond(s string) {
	fmt.Fprintf(c.data, ">>%s\n", s)
	_, err := fmt.Fprint(c.conn, s, "\n")
	if err != nil {
		log.Print(err)
	}
}
