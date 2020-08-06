package ftplistener

import (
	"TukTuk/telegrambot"
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	status426 = "426 Bye."
	status220 = "220 TukTuk callback server."
)

//starting ftp server
func StartFTP(db *sql.DB) {
	server := fmt.Sprintf(":%d", 21)
	listener, err := net.Listen("tcp", server)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		log.Printf("[*] Connection Accepted from [%s]\n", conn.RemoteAddr().String())
		if err != nil {
			log.Print(err)
			continue
		}
		go handleFTP(conn, db)
	}
}

func handleFTP(c net.Conn, db *sql.DB) {
	defer c.Close()
	ServeFTP(newConn(c, db))
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
	_, err := c.db.Exec("insert into ftp (data, source_ip, time) values ($1, $2, $3)", c.data.String(), c.conn.RemoteAddr().String(), time.Now().String())
	if err != nil {
		log.Println(err)
	}
	//Send Alert to telegram
	telegrambot.BotSendAlert(c.data.String(), c.conn.RemoteAddr().String(), time.Now().String(), "FTP")
}

func (c *Conn) respond(s string) {
	log.Print(">> ", s)
	fmt.Fprintf(c.data, ">>%s\n", s)
	_, err := fmt.Fprint(c.conn, s, "\n")
	if err != nil {
		log.Print(err)
	}
}

//main request handler
//if input differs from USER we just drop the connection and don't log it
func ServeFTP(c *Conn) {
	c.respond(status220)
	s := bufio.NewScanner(c.conn)
	for s.Scan() {
		input := strings.Fields(s.Text())
		if len(input) == 0 {
			continue
		}
		command, args := input[0], input[1:]
		log.Printf("<< %s %v", command, args)
		fmt.Fprintf(c.data, "<< %s %v\n", command, args)
		switch command {
		case "USER":
			c.respond(status426)
			c.log()
			return
		default:
			c.respond(status426)
			return
		}
	}
	if s.Err() != nil {
		log.Print(s.Err())
	}
}
