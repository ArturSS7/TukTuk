package backend

import (
	"TukTuk/database"
	"TukTuk/ldaplistener"

	"log"
	"time"

	"github.com/labstack/echo"
)

func startLDAP(c echo.Context) error {
	cc := c.(*database.DBContext)
	e := make(chan error)
	go func(r chan error) {
		err := ldaplistener.StartLDAP(cc.Db)
		if err != nil {
			r <- err
		}
	}(e)
	time.Sleep(2 * time.Second)
	select {
	case err := <-e:
		log.Println(err)
		close(e)
		er := &TcpErr{Error: err}
		return c.JSON(200, er)
	default:
		res := &TcpResult{Port: "10389", Success: true}
		close(e)
		return c.JSON(200, res)
	}
}

func shutdownLDAP(c echo.Context) error {
	if ldaplistener.LDAPserver != nil {
		ldaplistener.LDAPserver.Stop()
		ldaplistener.LDAPserver = nil
		return c.JSON(200, TcpResult{
			Port:    "10389",
			Success: true,
		})
	}
	return c.JSON(200, TcpResult{
		Port:    "10389",
		Success: false,
	})
}
