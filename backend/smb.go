package backend

import (
	"TukTuk/smblistener"
	"fmt"
	"github.com/labstack/echo"
	"log"
	"time"
)

func startSMBServer(c echo.Context) error {
	e := make(chan error)
	go func(r chan error) {
		err := smblistener.CMD.Run()
		r <- err
	}(e)
	time.Sleep(1 * time.Second)
	fmt.Println(smblistener.CMD.Process)
	select {
	case err := <-e:
		log.Println(err)
		return c.JSON(200, struct {
			Error error
		}{Error: err})
	default:
		return c.JSON(200, struct {
			Success bool
		}{Success: true})
	}
}

func stopSMBServer(c echo.Context) error {
	err := smblistener.CMD.Process.Kill()
	if err != nil {
		log.Println(err)
		return c.JSON(200, struct {
			error error
		}{error: err})
	}
	return c.JSON(200, struct {
		success bool
	}{success: true})
}
