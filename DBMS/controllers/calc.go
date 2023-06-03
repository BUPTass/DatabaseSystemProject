package controllers

import (
	"compress/gzip"
	"database/sql"
	"encoding/xml"
	"io/fs"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gonum/stat/distuv"
	"github.com/labstack/echo/v4"
)

type tbC2Inew struct {
	SCELL    string  `db:"SCELL"`
	NCELL    string  `db:"NCELL"`
	RSRPmean float32 `db:"RSRPmean"`
	RSRPstd  float32 `db:"RSRPstd"`
	PrbC2I9  float32 `db:"PrbC2I9"`
	PrbABS6  float32 `db:"PrbABS6"`
}
type tbC2I3 struct {
	A string `db:"a"`
	B string `db:"b"`
	C string `db:"c"`
}
type tbMRODatanew struct {
	MroID             int
	ServingSector     string
	InterferingSector string
	LteScRSRP         int
	LteNcRSRP         int
	LteNcEarfcn       int
	LteNcPci          int
}

func C2InewCalc(c echo.Context) error {
	//param
	minC := c.QueryParam("min")
	type dataTMP struct {
		ServingSector     string  `db:"ServingSector"`
		InterferingSector string  `db:"InterferingSector"`
		Mean              float32 `db:"mean"`
		Std               float32 `db:"std"`
	}
	var ans1 []dataTMP
	stmt := "select ServingSector,InterferingSector,AVG(LteScRSRP-LteNcRSRP) as mean,STDDEV(LteScRSRP-LteNcRSRP) as std " +
		"from tbMROData " +
		"group by ServingSector,InterferingSector " +
		"having count(*) >= ?"
	err := db.Select(&ans1, stmt, minC)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	var prbc2i9 []float32
	var prbabs6 []float32
	for i := 0; i < len(ans1); i++ {
		f := distuv.Normal{Mu: float64(ans1[i].Mean), Sigma: float64(ans1[i].Std)}
		prbc2i9 = append(prbc2i9, float32(f.CDF(float64(9.0))))
		prbabs6 = append(prbabs6, float32(f.CDF(float64(6.0))-f.CDF(float64(-6.0))))
	}
	//查询是否存在
	var tableName string
	err = db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name = ?", dataDbName_, "tbC2Inew").Scan(&tableName)
	if err == sql.ErrNoRows {
		createTableStmt := "CREATE TABLE tbC2Inew (`SCELL` nvarchar(255),`NCELL` nvarchar(255), `RSRPmean` float, `RSRPstd` float, `PrbC2I9` float, `PrbABS6` float, PRIMARY KEY (`SCELL`,`NCELL`));"
		_, err = db.Exec(createTableStmt)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	} else if err != nil {
		log.Println(err.Error())
		return err
	}
	var ans []tbC2Inew
	if err == nil {
		db.Exec("delete from tbC2Inew")
		insertStmt, err := db.Prepare("insert into tbC2Inew values(?,?,?,?,?,?)")
		if err != nil {
			log.Println(err.Error())
			return err
		} else {
			for i := 0; i < len(ans1); i++ {
				ans = append(ans, tbC2Inew{ans1[i].ServingSector, ans1[i].InterferingSector, ans1[i].Mean, ans1[i].Std, prbc2i9[i], prbabs6[i]})
				_, err = insertStmt.Exec(ans1[i].ServingSector, ans1[i].InterferingSector, ans1[i].Mean, ans1[i].Std, prbc2i9[i], prbabs6[i])
				if err != nil {
					log.Println(err.Error())
					return err
				}
			}
		}
	} else {
		log.Println(err.Error())
		return err
	}
	return c.JSON(http.StatusOK, ans)
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
		log.Println(err.Error())
		return err
	}
	//检查表tbC2I3是否存在
	err = db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name = ?", dataDbName_, "tbC2I3").Scan(&tableName)
	if err == sql.ErrNoRows {
		createTableStmt := "CREATE TABLE tbC2I3 (`a` nvarchar(255),`b` nvarchar(255),`c` nvarchar(255),PRIMARY KEY (`a`,`b`,`c`));"
		_, err = db.Exec(createTableStmt)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	} else if err != nil {
		log.Println(err.Error())
		return err
	} else {
		db.Exec("delete from tbC2I3")
	}
	//计算三元组
	var tb3tmp []tbC2I3
	SelectStmt := "select A.SCELL as a,B.SCELL as b,B.NCELL as c from (tbC2Inew as A join tbC2Inew as B on A.NCELL = B.SCELL) join tbC2Inew as C on (C.SCELL = B.NCELL and C.NCELL = A.SCELL) or (C.SCELL = A.SCELL and C.NCELL = B.NCELL) where A.PrbABS6 >= ? and B.PrbABS6 >= ? and C.PrbABS6 >= ?"
	err = db.Select(&tb3tmp, SelectStmt, x, x, x)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	//去重
	s := make(map[tbC2I3]bool)
	for i := 0; i < len(tb3tmp); i++ {
		if tb3tmp[i].A > tb3tmp[i].B {
			tmp := tb3tmp[i].A
			tb3tmp[i].A = tb3tmp[i].B
			tb3tmp[i].B = tmp
		}
		if tb3tmp[i].B > tb3tmp[i].C {
			tmp := tb3tmp[i].B
			tb3tmp[i].B = tb3tmp[i].C
			tb3tmp[i].C = tmp
		}
		if tb3tmp[i].A > tb3tmp[i].B {
			tmp := tb3tmp[i].A
			tb3tmp[i].A = tb3tmp[i].B
			tb3tmp[i].B = tmp
		}
		s[tb3tmp[i]] = true
	}
	insertStmt, err := db.Prepare("insert into tbC2I3 values(?,?,?)")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	var ans []tbC2I3
	for k, v := range s {
		if !v {
			return c.String(http.StatusOK, "dynamic error")
		}
		_, err = insertStmt.Exec(k.A, k.B, k.C)
		ans = append(ans, k)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}
	return c.JSON(http.StatusOK, ans)
}
func MROMREcalc(c echo.Context) error {
	filePath := c.QueryParam("filePath")

	//step1
	var files []fs.DirEntry
	filesTmp, err := os.ReadDir(filePath)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	reg, _ := regexp.Compile("MRO")
	for _, file := range filesTmp {
		if !file.IsDir() {
			res := reg.MatchString(file.Name())
			log.Println(file.Name())
			if res {
				files = append(files, file)
			}
		}
	}
	var table []tbMRODatanew
	//step2
	for _, file := range files {
		//解压
		gzFile, err := os.Open(filePath + "\\" + file.Name())
		if err != nil {
			log.Println(err.Error())
			continue
		}
		defer gzFile.Close()
		reader, err := gzip.NewReader(gzFile)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		defer reader.Close()
		//提取xml文件信息
		type XMLdata struct {
			ENBid      xml.Name `xml:"eNB"`
			Mesurement []struct {
				Smr    string   `xml:"smr"`
				Object []string `xml:"v"`
			} `xml:"mesurement"`
		}
		XMLdecoder := xml.NewDecoder(reader)
		data := XMLdata{}
		err = XMLdecoder.Decode(&data)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		type Result struct {
			Id      int
			SPci    int
			NPci    int
			NEarfcn int
			SRSRP   int
			NRSRP   int
		}
		minlen := len(data.Mesurement[0].Object)
		for _, i := range data.Mesurement {
			if minlen < len(i.Object) {
				minlen = len(i.Object)
			}
		}
		var result []Result
		for i := 0; i < minlen; i++ {
			tmp := Result{}
			for _, j := range data.Mesurement {
				if j.Smr == "MR.LteScPci" {
					tmp.SPci, _ = strconv.Atoi(j.Object[i])
				} else if j.Smr == "MR.LteNcPci" {
					tmp.NPci, _ = strconv.Atoi(j.Object[i])
				} else if j.Smr == "MR.LteNcEarfcn" {
					tmp.NEarfcn, _ = strconv.Atoi(j.Object[i])
				} else if j.Smr == "MR.LteScRSRP" {
					tmp.SRSRP, _ = strconv.Atoi(j.Object[i])
				} else if j.Smr == "MR.LteNcRSRP" {
					tmp.NRSRP, _ = strconv.Atoi(j.Object[i])
				} else if j.Smr == "eNBid" {
					tmp.Id, _ = strconv.Atoi(j.Object[i])
				}
			}
			result = append(result, tmp)
		}
		//获取主临小区ID
		Idtable := make(map[int]string)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		for _, i := range result {
			if Idtable[i.SPci] == "" {
				var Sid []string
				err = db.Select(&Sid, "select SECTOR_ID from tbCell where PCI = ?", i.SPci)
				if err != nil {
					log.Println(err.Error())
					return err
				}
				if len(Sid) != 1 {
					Idtable[i.SPci] = "error"
					continue
				} else {
					Idtable[i.SPci] = Sid[0]
				}
			}
			if Idtable[i.NPci] == "" {
				var Nid []string
				err = db.Select(&Nid, "select SECTOR_ID from tbCell where PCI = ?", i.NPci)
				if err != nil {
					log.Println(err.Error())
					return err
				}
				if len(Nid) != 1 {
					Idtable[i.NPci] = "error"
					continue
				}
				Idtable[i.NPci] = Nid[0]
			}
		}
		//过滤
		var EarfcnList []int
		Earfcn := make(map[int]bool)
		err = db.Select(&EarfcnList, "select distinct EARFCN from tbCell")
		if err != nil {
			log.Println(err.Error())
			return err
		}
		for _, i := range EarfcnList {
			Earfcn[i] = true
		}
		var tmp tbMRODatanew
		for _, i := range result {
			if i.NPci < 0 || i.NPci > 503 || i.NRSRP < 0 || i.NRSRP > 97 || i.SRSRP < 0 || i.SRSRP > 97 || !Earfcn[i.NEarfcn] || Idtable[i.SPci] == "" || Idtable[i.SPci] == "error" || Idtable[i.NPci] == "" || Idtable[i.NPci] == "error" {
				continue
			}
			tmp.MroID = i.Id
			tmp.ServingSector = Idtable[i.SPci]
			tmp.InterferingSector = Idtable[i.NPci]
			tmp.LteScRSRP = i.SRSRP
			tmp.LteNcRSRP = i.NRSRP
			tmp.LteNcEarfcn = i.NEarfcn
			tmp.LteNcPci = i.NPci
			table = append(table, tmp)
		}
	}
	//检查表tbMROdatanew是否存在
	var tableName string
	err = db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name = ?", dataDbName_, "tbMROdatanew").Scan(&tableName)
	if err == sql.ErrNoRows {
		createTableStmt := "CREATE TABLE tbMROdatanew (`MroID` int,`ServingSector` varchar(15), `InterferingSector` varchar(15), `LteScRSRP` int, `LteNcRSRP` int, `LteNcEarfcn` int, `LteNcPci` smallint,PRIMARY KEY (`MroID`,`ServingSector`,`InterferingSector`));"
		_, err = db.Exec(createTableStmt)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	} else if err != nil {
		log.Println(err.Error())
		return err
	} else {
		db.Exec("delete from tbMROdatanew")
	}
	//写入数据库
	insertStmt, err := db.Prepare("insert into tbMROdatanew values(?,?,?,?,?,?,?)")
	if err != nil {
		log.Println(err.Error())
		return err
	} else {
		for _, i := range table {
			_, err = insertStmt.Exec(i.MroID, i.ServingSector, i.InterferingSector, i.LteScRSRP, i.LteNcRSRP, i.LteNcEarfcn, i.LteNcPci)
			if err != nil {
				log.Println(err.Error())
				return err
			}
		}
	}
	return c.String(http.StatusOK, "success")
}
