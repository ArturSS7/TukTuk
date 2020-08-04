package backend

import (
	"TukTuk/database"
	"database/sql"
	"github.com/labstack/echo"
	"log"
)

type Request struct {
	Data     string `json:"data"`
	SourceIp string `json:"source_ip"`
	Time     string `json:"time"`
}

type Result struct {
	Error string `json:"error"`
}

//стартуем наш бек
func StartBack(db *sql.DB) {
	e := echo.New()
	//делаем так чтобы поинтер на коннект можно было юзать внутри нттп хенделров
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &database.DBContext{Context: c, Db: db}
			return h(cc)
		}
	})
	e.GET("/api/http", getHTTP)
	e.Logger.Fatal(e.Start(":1234"))
}

//хендел для получения данных нттп с бека, селектим с базы с заданным лимитом
func getHTTP(c echo.Context) error {
	limit := c.FormValue("limit")
	cc := c.(*database.DBContext)
	rows, err := cc.Db.Query("select data, source_ip, time from http order by id limit $1", limit)
	if err != nil {
		log.Println(err)
		er := &Result{Error: "true"}
		return c.JSON(200, er)
	}
	rr := make([]Request, 0)
	for rows.Next() {
		r := Request{}
		err := rows.Scan(&r.Data, &r.SourceIp, &r.Time)
		if err != nil {
			log.Println(err)
			er := &Result{Error: "true"}
			return c.JSON(200, er)
		}
		rr = append(rr, r)
	}
	return c.JSON(200, rr)
}
