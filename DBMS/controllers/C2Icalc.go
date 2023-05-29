package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gonum/stat/distuv"
	"github.com/labstack/echo/v4"
)

type tbC2Inew struct {
	SCELL    string  `db:"SCELL"`
	NCELL    string  `db:"NCELL"`
	RSRPmean int     `db:"RSRPmean"`
	RSRPstd  int     `db:"RSRPstd"`
	PrbC2I9  float32 `db:"PrbC2I9"`
	PrbABS6  float32 `db:"PrbABS6"`
}
type tbC2I3 struct {
	a string `db:"a"`
	b string `db:"b"`
	c string `db:"c"`
}

func C2InewCalc(c echo.Context) error {
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
		createTableStmt := "CREATE TABLE tbC2Inew ('SCELL' nvarchar(255),'NCELL' nvarchar(255),'RSRPmean' float,'RSRPstd' float,'PrbC2I9' float,'PrbABS6' float,PRIMARY KEY ('SCELL','NCELL'));"
		_, err = db.Exec(createTableStmt)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	if err == nil {
		db.Exec("delete from tbC2Inew")
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
func C2I3Calc(c echo.Context) error {
	x := c.QueryParam("x")
	//检查表tbC2Inew是否存在
	var tableName string
	err := db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name = ?", dataDbName_, "tbC2Inew").Scan(&tableName)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.String(http.StatusOK, "tbC2Inew not exist, please calculate tbC2Inew first")
		}
		return err
	}
	//检查表tbC2I3是否存在
	err = db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name = ?", dataDbName_, "tbC2I3").Scan(&tableName)
	if err == sql.ErrNoRows {
		createTableStmt := "CREATE TABLE tbC2I3 ('a' nvarchar(255),'b' nvarchar(255),'c' nvarchar(255),PRIMARY KEY ('a','b','c'));"
		_, err = db.Exec(createTableStmt)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		db.Exec("delete from tbC2I3")
	}
	//计算三元组
	var tbnew []tbC2Inew
	err = db.Select(&tbnew, "select * from tbC2Inew")
	if err != nil {
		return err
	}
	var tb3tmp []tbC2I3
	SelectStmt := "select A.SCELL as a,B.SCELL as b,C.SCELL as c from (tbC2Inew as A right join tbC2Inew as B on A.NCELL = B.SCELL) right join tbC2Inew as C on B.NCELL = C.SCELL and (C.NCELL = A.SCELL or A.NCELL = C.SCELL) where A.PrbABS6 >= ? and B.PrbABS6 >= ? and C.PrbABS6 >= ?"
	err = db.Select(&tb3tmp, SelectStmt, x, x, x)
	if err != nil {
		return err
	}
	//去重
	var s map[tbC2I3]bool
	for i := 0; i < len(tb3tmp); i++ {
		if tb3tmp[i].a > tb3tmp[i].b {
			tmp := tb3tmp[i].a
			tb3tmp[i].a = tb3tmp[i].b
			tb3tmp[i].b = tmp
		}
		if tb3tmp[i].b > tb3tmp[i].c {
			tmp := tb3tmp[i].b
			tb3tmp[i].b = tb3tmp[i].c
			tb3tmp[i].c = tmp
		}
		if tb3tmp[i].a > tb3tmp[i].b {
			tmp := tb3tmp[i].a
			tb3tmp[i].a = tb3tmp[i].b
			tb3tmp[i].b = tmp
		}
		s[tb3tmp[i]] = true
	}
	insertStmt, err := db.Prepare("insert into tbC2I3 values(?,?,?)")
	for k, _ := range s {
		_, err = insertStmt.Exec(k.a, k.b, k.c)
		if err != nil {
			return err
		}
	}
	return c.String(http.StatusOK, "success")
}
