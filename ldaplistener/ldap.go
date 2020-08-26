// Listen to 10389 port for LDAP Request
// and route bind request to the handleBind func
package ldaplistener

import (
	"TukTuk/database"
	"TukTuk/emailalert"
	"TukTuk/telegrambot"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	ldap "github.com/vjeantet/ldapserver"
)

var LDAPserver *ldap.Server

func StartLDAP(db *sql.DB) error {
	//ldap logger
	log.Println(os.Stdout, "[LDAP Server] ", log.LstdFlags)

	//Create a new LDAP Server
	server := ldap.NewServer()

	routes := ldap.NewRouteMux()
	routes.Bind(handleBind)
	server.Handle(routes)
	LDAPserver = server
	// listen on 10389
	server.ListenAndServe("pwn.bar:10389")

	// When CTRL+C, SIGINT and SIGTERM signal occurs
	// Then stop server gracefully
	//ch := make(chan os.Signal)
	//	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	//	<-ch
	//(ch)

	server.Stop()
	return nil
}

// handleBind return Success if login == mysql
func handleBind(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetBindRequest()
	res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
	log.Println(m.Client.Addr())
	log.Println(r.Name())
	strname := fmt.Sprintf("%v", r.Name())
	logLDAP(strname, m.Client.Addr().String())
	w.Write(res)
	return

}

func logLDAP(dn, remouteAddr string) {
	var lastInsertId int64 = 0
	err := database.DNSDB.QueryRow("insert into ldap (data, source_ip, time) values ($1, $2, $3) RETURNING id", dn, remouteAddr, time.Now().String()).Scan(&lastInsertId)

	if err != nil {
		log.Println(err)
	}

	//Send Alert to telegram
	telegrambot.BotSendAlert(dn, remouteAddr, time.Now().String(), "LDAP", lastInsertId)
	//Send Alert to email
	emailalert.SendEmailAlert("LDAP Alert", "Remoute Address: "+remouteAddr+"\n"+dn+"\n"+time.Now().String())
}
