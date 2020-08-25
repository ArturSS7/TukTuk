package backend

import (
	"TukTuk/database"
	"TukTuk/ftplistener"
	"github.com/labstack/echo"
	"log"
	"time"
)

func startFTP(c echo.Context) error {
	cc := c.(*database.DBContext)
	e := make(chan error)
	go func(r chan error) {
		err := ftplistener.StartFTP(cc.Db)
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
		res := &TcpResult{Port: "21", Success: true}
		close(e)
		return c.JSON(200, res)
	}
}

func shutdownFTP(c echo.Context) error {
	if ftplistener.FTPServer != nil {
		ftplistener.FTPServer.Stop()
		ftplistener.FTPServer = nil
		return c.JSON(200, TcpResult{
			Port:    "21",
			Success: true,
		})
	}
	return c.JSON(200, TcpResult{
		Port:    "21",
		Success: false,
	})
}
