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
			Error string `json:"error"`
		}{Error: fmt.Sprintf("%v", err)})
	default:
		return c.JSON(200, struct {
			Success bool `json:"success"`
		}{Success: true})
	}
}

func stopSMBServer(c echo.Context) error {
	if smblistener.CMD.Process == nil {
		return c.JSON(200, struct {
			Success bool `json:"success"`
		}{Success: false})
	}
	err := smblistener.CMD.Process.Kill()
	if err != nil {
		log.Println(err)
		return c.JSON(200, struct {
			Error string `json:"error"`
		}{Error: fmt.Sprintf("%v", err)})
	}
	return c.JSON(200, struct {
		Success bool `json:"success"`
	}{Success: true})
}
