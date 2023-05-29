package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gonum/stat/distuv"
	"github.com/labstack/echo/v4"
)

type tbC2Inew struct {
	SCell    string
	NCell    string
	meanRSRP float32
	stdRSRP  float32
	P9       float32
	P6       float32
}
type tbMROData struct {
	TimeStamp         string
	ServingSector     string
	InterferingSector string
	LteScRSRP         float32
	LteNcRSRP         float32
	LteNcEarfcn       int
	LteNcPci          int
}

func C2Icalc(c echo.Context) error {
	//param
	minC := c.QueryParam("min")
	type dataTMP struct {
		ServingSector     string
		InterferingSector string
		mean              float32
		std               float32
	}
	var ans1 []dataTMP
	stmt := "select ServingSector,InterferingSector,AVG(LteScRSRP-LteNcRSRP) as mean,STDDEV(LteScRSRP-LteNcRSRP) as std" +
		"from tbMROData" +
		"group by ServingSector,InterferingSector" +
		"having count(*) >= " + minC + " )"
	err := db.Select(&ans1, stmt)
	if err != nil {
		return err
	}
	var prbc2i9 []float32
	var prbabs6 []float32
	for i := 0; i < len(ans1); i++ {
		f := distuv.Normal{Mu: float64(ans1[i].mean), Sigma: float64(ans1[i].std)}
		prbc2i9 = append(prbc2i9, float32(f.CDF(float64(9.0))))
		prbabs6 = append(prbabs6, float32(f.CDF(float64(6.0))-f.CDF(float64(-6.0))))
	}
	//查询是否存在
	var tableName string
	err = db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name = ?", dataDbName_, "tbC2Inew").Scan(&tableName)
	if err == sql.ErrNoRows {
		createTableStmt := "CREATE TABLE 'tbC2Inew' ('SCELL' nvarchar(255),'NCELL' nvarchar(255),'RSRPmean' float,'RSRPstd' float,'PrbC2I9' float,'PrbABS6' float,PRIMARY KEY ('SCELL','NCELL'));"
		_, err = db.Exec(createTableStmt)
		if err != nil {
			log.Fatal(err)
		}
	} else if err == nil {
		insertStmt, err := db.Prepare("insert into tbC2Inew values(?,?,?,?,?,?)")
		if err != nil {
			return err
		} else {
			for i := 0; i < len(ans1); i++ {
				_, err = insertStmt.Exec(ans1[i].ServingSector, ans1[i].InterferingSector, ans1[i].mean, ans1[i].std, prbc2i9[i], prbabs6[i])
				if err != nil {
					return err
				}
			}
		}
	} else {
		return err
	}
	return c.String(http.StatusOK, "success")
}
