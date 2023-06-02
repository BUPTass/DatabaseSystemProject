package Query

import (
	"DatabaseSystemProject/Import"
	"database/sql"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"time"
)

type Enodeb struct {
	City       string  `json:"city"`
	EnodebID   int     `json:"enodebID"`
	EnodebName string  `json:"enodebName"`
	Vendor     string  `json:"vendor"`
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
	Style      string  `json:"style"`
}

// GetCellInfo gets cell information from the tbCell table based on the provided input
func GetCellInfo(db *sql.DB, query string) ([]Import.CellValues, error) {
	// Prepare the SQL statement
	stmt, err := db.Prepare(`
		SELECT * FROM tbCell
		WHERE Sector_ID = ? OR Sector_Name = ?
	`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()

	// Execute the query with the provided input
	rows, err := stmt.Query(query, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and retrieve cell information
	var cellInfoList []Import.CellValues
	for rows.Next() {
		var cellInfo Import.CellValues
		err := rows.Scan(
			&cellInfo.City,
			&cellInfo.SectorID,
			&cellInfo.SectorName,
			&cellInfo.EnodebID,
			&cellInfo.EnodebName,
			&cellInfo.EARFCN,
			&cellInfo.PCI,
			&cellInfo.PSS,
			&cellInfo.SSS,
			&cellInfo.TAC,
			&cellInfo.Vendor,
			&cellInfo.Longitude,
			&cellInfo.Latitude,
			&cellInfo.Style,
			&cellInfo.Azimuth,
			&cellInfo.Height,
			&cellInfo.ElectTilt,
			&cellInfo.MechTilt,
			&cellInfo.TotlTilt,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		cellInfoList = append(cellInfoList, cellInfo)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return cellInfoList, nil
}

// GetAllSectorNames retrieves all the SectorNames from the tbCell table
func GetAllSectorNames(db *sql.DB) ([]string, error) {
	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT SECTOR_NAME FROM tbCell")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and retrieve the SectorNames
	var sectorNames []string
	for rows.Next() {
		var sectorName string
		err := rows.Scan(&sectorName)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		sectorNames = append(sectorNames, sectorName)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return sectorNames, nil
}

func GetEnodeb(db *sql.DB, query string) ([]Enodeb, error) {
	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT * FROM tbEnodeb WHERE ENODEBID = ? OR ENODEB_NAME = ?")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query(query, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	enodebs := make([]Enodeb, 0)
	for rows.Next() {
		enodeb := Enodeb{}
		err := rows.Scan(&enodeb.City, &enodeb.EnodebID, &enodeb.EnodebName, &enodeb.Vendor, &enodeb.Longitude,
			&enodeb.Latitude, &enodeb.Style)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		enodebs = append(enodebs, enodeb)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return enodebs, nil
}

// GetAllEnodebNames retrieves all the ENODEB_NAME from the tbEnodeb table
func GetAllEnodebNames(db *sql.DB) ([]string, error) {
	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT ENODEB_NAME FROM tbEnodeb ")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and retrieve the ENODEBNames
	var enodebNames []string
	for rows.Next() {
		var enodebName string
		err := rows.Scan(&enodebName)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		enodebNames = append(enodebNames, enodebName)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return enodebNames, nil
}

// GetKPIInfoBySectorName retrieves all the KPI information from the tbKPI table based on SECTOR_NAME
func GetKPIInfoBySectorName(db *sql.DB, sectorName string) ([]Import.KpiValues, error) {
	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT * FROM tbKPI WHERE SECTOR_NAME = ?")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query(sectorName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var results []Import.KpiValues

	for rows.Next() {
		kpiData := Import.KpiValues{}

		// Scan the row values into the KPIData struct fields
		err = rows.Scan(
			&kpiData.StartTime,
			&kpiData.ENODEB_NAME,
			&kpiData.SECTOR_DESCRIPTION,
			&kpiData.SECTOR_NAME,
			&kpiData.RCCConnSUCC,
			&kpiData.RCCConnATT,
			&kpiData.RCCConnRATE,
			&kpiData.ERABConnSUCC,
			&kpiData.ERABConnATT,
			&kpiData.ERABConnRATE,
			&kpiData.ENODEB_ERABRel,
			&kpiData.SECTOR_ERABRel,
			&kpiData.ERABDropRateNew,
			&kpiData.WirelessAccessRateAY,
			&kpiData.ENODEB_UECtxRel,
			&kpiData.UEContextRel,
			&kpiData.UEContextSUCC,
			&kpiData.WirelessDropRate,
			&kpiData.ENODEB_InterFreqHOOutSUCC,
			&kpiData.ENODEB_InterFreqHOOutATT,
			&kpiData.ENODEB_IntraFreqHOOutSUCC,
			&kpiData.ENODEB_IntraFreqHOOutATT,
			&kpiData.ENODEB_InterFreqHOInSUCC,
			&kpiData.ENODEB_InterFreqHOInATT,
			&kpiData.ENODEB_IntraFreqHOInSUCC,
			&kpiData.ENODEB_IntraFreqHOInATT,
			&kpiData.ENODEB_HOInRate,
			&kpiData.ENODEB_HOOutRate,
			&kpiData.IntraFreqHOOutRateZSP,
			&kpiData.InterFreqHOOutRateZSP,
			&kpiData.HOSuccessRate,
			&kpiData.PDCP_UplinkThroughput,
			&kpiData.PDCP_DownlinkThroughput,
			&kpiData.RRCRebuildReq,
			&kpiData.RRCRebuildRate,
			&kpiData.SourceENB_IntraFreqHOOutSUCC,
			&kpiData.SourceENB_InterFreqHOOutSUCC,
			&kpiData.SourceENB_IntraFreqHOInSUCC,
			&kpiData.SourceENB_InterFreqHOInSUCC,
			&kpiData.ENODEB_HOOutSUCC,
			&kpiData.ENODEB_HOOutATT,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		results = append(results, kpiData)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return results, nil
}

// GetKPISectorNames retrieves all the SectorNames from the tbKPI table
func GetKPISectorNames(db *sql.DB) ([]string, error) {
	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT UNIQUE SECTOR_NAME FROM tbKPI")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and retrieve the SectorNames
	var sectorNames []string
	for rows.Next() {
		var sectorName string
		err := rows.Scan(&sectorName)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		sectorNames = append(sectorNames, sectorName)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return sectorNames, nil
}

func GeneratePRBNewTable(db *sql.DB, outputPath string) (string, error) {
	// Create tbPRBNew table
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS tbPRBNew (
			SECTOR_NAME        VARCHAR(255) NOT NULL,
			Hour               DATETIME NOT NULL,
			PRB_Interference   DECIMAL(7, 4) NOT NULL,
			PRB00              DECIMAL(7, 4)          null,
    		PRB01              DECIMAL(7, 4)          null,
    		PRB02              DECIMAL(7, 4)          null,
    		PRB03              DECIMAL(7, 4)          null,
    		PRB04              DECIMAL(7, 4)          null,
    		PRB05              DECIMAL(7, 4)          null,
    		PRB06              DECIMAL(7, 4)          null,
    		PRB07              DECIMAL(7, 4)          null,
    		PRB08              DECIMAL(7, 4)          null,
    		PRB09              DECIMAL(7, 4)          null,
    		PRB10              DECIMAL(7, 4)          null,
    		PRB11              DECIMAL(7, 4)          null,
    		PRB12              DECIMAL(7, 4)          null,
    		PRB13              DECIMAL(7, 4)          null,
    		PRB14              DECIMAL(7, 4)          null,
    		PRB15              DECIMAL(7, 4)          null,
    		PRB16              DECIMAL(7, 4)          null,
    		PRB17              DECIMAL(7, 4)          null,
    		PRB18              DECIMAL(7, 4)          null,
    		PRB19              DECIMAL(7, 4)          null,
    		PRB20              DECIMAL(7, 4)          null,
    		PRB21              DECIMAL(7, 4)          null,
    		PRB22              DECIMAL(7, 4)          null,
    		PRB23              DECIMAL(7, 4)          null,
    		PRB24              DECIMAL(7, 4)          null,
    		PRB25              DECIMAL(7, 4)          null,
    		PRB26              DECIMAL(7, 4)          null,
    		PRB27              DECIMAL(7, 4)          null,
    		PRB28              DECIMAL(7, 4)          null,
    		PRB29              DECIMAL(7, 4)          null,
    		PRB30              DECIMAL(7, 4)          null,
    		PRB31              DECIMAL(7, 4)          null,
    		PRB32              DECIMAL(7, 4)          null,
    		PRB33              DECIMAL(7, 4)          null,
    		PRB34              DECIMAL(7, 4)          null,
    		PRB35              DECIMAL(7, 4)          null,
    		PRB36              DECIMAL(7, 4)          null,
    		PRB37              DECIMAL(7, 4)          null,
    		PRB38              DECIMAL(7, 4)          null,
    		PRB39              DECIMAL(7, 4)          null,
    		PRB40              DECIMAL(7, 4)          null,
    		PRB41              DECIMAL(7, 4)          null,
    		PRB42              DECIMAL(7, 4)          null,
    		PRB43              DECIMAL(7, 4)          null,
    		PRB44              DECIMAL(7, 4)          null,
    		PRB45              DECIMAL(7, 4)          null,
    		PRB46              DECIMAL(7, 4)          null,
    		PRB47              DECIMAL(7, 4)          null,
    		PRB48              DECIMAL(7, 4)          null,
    		PRB49              DECIMAL(7, 4)          null,
    		PRB50              DECIMAL(7, 4)          null,
    		PRB51              DECIMAL(7, 4)          null,
    		PRB52              DECIMAL(7, 4)          null,
    		PRB53              DECIMAL(7, 4)          null,
    		PRB54              DECIMAL(7, 4)          null,
    		PRB55              DECIMAL(7, 4)          null,
    		PRB56              DECIMAL(7, 4)          null,
    		PRB57              DECIMAL(7, 4)          null,
    		PRB58              DECIMAL(7, 4)          null,
    		PRB59              DECIMAL(7, 4)          null,
    		PRB60              DECIMAL(7, 4)          null,
    		PRB61              DECIMAL(7, 4)          null,
    		PRB62              DECIMAL(7, 4)          null,
    		PRB63              DECIMAL(7, 4)          null,
    		PRB64              DECIMAL(7, 4)          null,
    		PRB65              DECIMAL(7, 4)          null,
    		PRB66              DECIMAL(7, 4)          null,
    		PRB67              DECIMAL(7, 4)          null,
    		PRB68              DECIMAL(7, 4)          null,
    		PRB69              DECIMAL(7, 4)          null,
    		PRB70              DECIMAL(7, 4)          null,
    		PRB71              DECIMAL(7, 4)          null,
    		PRB72              DECIMAL(7, 4)          null,
    		PRB73              DECIMAL(7, 4)          null,
    		PRB74              DECIMAL(7, 4)          null,
    		PRB75              DECIMAL(7, 4)          null,
    		PRB76              DECIMAL(7, 4)          null,
    		PRB77              DECIMAL(7, 4)          null,
    		PRB78              DECIMAL(7, 4)          null,
    		PRB79              DECIMAL(7, 4)          null,
    		PRB80              DECIMAL(7, 4)          null,
    		PRB81              DECIMAL(7, 4)          null,
    		PRB82              DECIMAL(7, 4)          null,
    		PRB83              DECIMAL(7, 4)          null,
    		PRB84              DECIMAL(7, 4)          null,
    		PRB85              DECIMAL(7, 4)          null,
    		PRB86              DECIMAL(7, 4)          null,
    		PRB87              DECIMAL(7, 4)          null,
    		PRB88              DECIMAL(7, 4)          null,
    		PRB89              DECIMAL(7, 4)          null,
    		PRB90              DECIMAL(7, 4)          null,
    		PRB91              DECIMAL(7, 4)          null,
    		PRB92              DECIMAL(7, 4)          null,
    		PRB93              DECIMAL(7, 4)          null,
    		PRB94              DECIMAL(7, 4)          null,
    		PRB95              DECIMAL(7, 4)          null,
    		PRB96              DECIMAL(7, 4)          null,
    		PRB97              DECIMAL(7, 4)          null,
    		PRB98              DECIMAL(7, 4)          null,
    		PRB99              DECIMAL(7, 4)          null,
			PRIMARY KEY (SECTOR_NAME, Hour)
		);
	`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Drop all before insertion
	_, err = db.Exec(`delete from tbPRBNew`)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Perform aggregation and insert into tbPRBnew
	aggregateSQL := `
		INSERT INTO tbPRBNew (SECTOR_NAME, Hour, 
		PRB00, PRB01, PRB02, PRB03, PRB04, PRB05, PRB06, PRB07, PRB08, PRB09,
		PRB10, PRB11, PRB12, PRB13, PRB14, PRB15, PRB16, PRB17, PRB18, PRB19,
		PRB20, PRB21, PRB22, PRB23, PRB24, PRB25, PRB26, PRB27, PRB28, PRB29,
		PRB30, PRB31, PRB32, PRB33, PRB34, PRB35, PRB36, PRB37, PRB38, PRB39,
		PRB40, PRB41, PRB42, PRB43, PRB44, PRB45, PRB46, PRB47, PRB48, PRB49,
		PRB50, PRB51, PRB52, PRB53, PRB54, PRB55, PRB56, PRB57, PRB58, PRB59,
		PRB60, PRB61, PRB62, PRB63, PRB64, PRB65, PRB66, PRB67, PRB68, PRB69,
		PRB70, PRB71, PRB72, PRB73, PRB74, PRB75, PRB76, PRB77, PRB78, PRB79,
		PRB80, PRB81, PRB82, PRB83, PRB84, PRB85, PRB86, PRB87, PRB88, PRB89,
		PRB90, PRB91, PRB92, PRB93, PRB94, PRB95, PRB96, PRB97, PRB98, PRB99)
		SELECT SECTOR_NAME, DATE_FORMAT(StartTime, '%Y-%m-%d %H:00:00') AS Hour,
		AVG(PRB00) AS PRB00,
		AVG(PRB01) AS PRB01,
		AVG(PRB02) AS PRB02,
		AVG(PRB03) AS PRB03,
		AVG(PRB04) AS PRB04,
		AVG(PRB05) AS PRB05,
		AVG(PRB06) AS PRB06,
		AVG(PRB07) AS PRB07,
		AVG(PRB08) AS PRB08,
		AVG(PRB09) AS PRB09,
		AVG(PRB10) AS PRB10,
		AVG(PRB11) AS PRB11,
		AVG(PRB12) AS PRB12,
		AVG(PRB13) AS PRB13,
		AVG(PRB14) AS PRB14,
		AVG(PRB15) AS PRB15,
		AVG(PRB16) AS PRB16,
		AVG(PRB17) AS PRB17,
		AVG(PRB18) AS PRB18,
		AVG(PRB19) AS PRB19,
		AVG(PRB20) AS PRB20,
		AVG(PRB21) AS PRB21,
		AVG(PRB22) AS PRB22,
		AVG(PRB23) AS PRB23,
		AVG(PRB24) AS PRB24,
		AVG(PRB25) AS PRB25,
		AVG(PRB26) AS PRB26,
		AVG(PRB27) AS PRB27,
		AVG(PRB28) AS PRB28,
		AVG(PRB29) AS PRB29,
		AVG(PRB30) AS PRB30,
		AVG(PRB31) AS PRB31,
		AVG(PRB32) AS PRB32,
		AVG(PRB33) AS PRB33,
		AVG(PRB34) AS PRB34,
		AVG(PRB35) AS PRB35,
		AVG(PRB36) AS PRB36,
		AVG(PRB37) AS PRB37,
		AVG(PRB38) AS PRB38,
		AVG(PRB39) AS PRB39,
		AVG(PRB40) AS PRB40,
		AVG(PRB41) AS PRB41,
		AVG(PRB42) AS PRB42,
		AVG(PRB43) AS PRB43,
		AVG(PRB44) AS PRB44,
		AVG(PRB45) AS PRB45,
		AVG(PRB46) AS PRB46,
		AVG(PRB47) AS PRB47,
		AVG(PRB48) AS PRB48,
		AVG(PRB49) AS PRB49,
		AVG(PRB50) AS PRB50,
		AVG(PRB51) AS PRB51,
		AVG(PRB52) AS PRB52,
		AVG(PRB53) AS PRB53,
		AVG(PRB54) AS PRB54,
		AVG(PRB55) AS PRB55,
		AVG(PRB56) AS PRB56,
		AVG(PRB57) AS PRB57,
		AVG(PRB58) AS PRB58,
		AVG(PRB59) AS PRB59,
		AVG(PRB60) AS PRB60,
		AVG(PRB61) AS PRB61,
		AVG(PRB62) AS PRB62,
		AVG(PRB63) AS PRB63,
		AVG(PRB64) AS PRB64,
		AVG(PRB65) AS PRB65,
		AVG(PRB66) AS PRB66,
		AVG(PRB67) AS PRB67,
		AVG(PRB68) AS PRB68,
		AVG(PRB69) AS PRB69,
		AVG(PRB70) AS PRB70,
		AVG(PRB71) AS PRB71,
		AVG(PRB72) AS PRB72,
		AVG(PRB73) AS PRB73,
		AVG(PRB74) AS PRB74,
		AVG(PRB75) AS PRB75,
		AVG(PRB76) AS PRB76,
		AVG(PRB77) AS PRB77,
		AVG(PRB78) AS PRB78,
		AVG(PRB79) AS PRB79,
		AVG(PRB80) AS PRB80,
		AVG(PRB81) AS PRB81,
		AVG(PRB82) AS PRB82,
		AVG(PRB83) AS PRB83,
		AVG(PRB84) AS PRB84,
		AVG(PRB85) AS PRB85,
		AVG(PRB86) AS PRB86,
		AVG(PRB87) AS PRB87,
		AVG(PRB88) AS PRB88,
		AVG(PRB89) AS PRB89,
		AVG(PRB90) AS PRB90,
		AVG(PRB91) AS PRB91,
		AVG(PRB92) AS PRB92,
		AVG(PRB93) AS PRB93,
		AVG(PRB94) AS PRB94,
		AVG(PRB95) AS PRB95,
		AVG(PRB96) AS PRB96,
		AVG(PRB97) AS PRB97,
		AVG(PRB98) AS PRB98,
		AVG(PRB99) AS PRB99
		FROM tbPRB
		WHERE StartTime BETWEEN '2020-07-17' AND '2020-07-20'
		GROUP BY SECTOR_NAME, Hour;
	`
	_, err = db.Exec(aggregateSQL)
	if err != nil {
		log.Println(err)
		return "", err
	}

	averageSQL := `UPDATE tbPRBNew
		SET PRB_Interference = (
		PRB00 + PRB01 + PRB02 + PRB03 + PRB04 + PRB05 + PRB06 + PRB07 + PRB08 + PRB09 +
		PRB10 + PRB11 + PRB12 + PRB13 + PRB14 + PRB15 + PRB16 + PRB17 + PRB18 + PRB19 +
		PRB20 + PRB21 + PRB22 + PRB23 + PRB24 + PRB25 + PRB26 + PRB27 + PRB28 + PRB29 +
		PRB30 + PRB31 + PRB32 + PRB33 + PRB34 + PRB35 + PRB36 + PRB37 + PRB38 + PRB39 +
		PRB40 + PRB41 + PRB42 + PRB43 + PRB44 + PRB45 + PRB46 + PRB47 + PRB48 + PRB49 +
		PRB50 + PRB51 + PRB52 + PRB53 + PRB54 + PRB55 + PRB56 + PRB57 + PRB58 + PRB59 +
		PRB60 + PRB61 + PRB62 + PRB63 + PRB64 + PRB65 + PRB66 + PRB67 + PRB68 + PRB69 +
		PRB70 + PRB71 + PRB72 + PRB73 + PRB74 + PRB75 + PRB76 + PRB77 + PRB78 + PRB79 +
		PRB80 + PRB81 + PRB82 + PRB83 + PRB84 + PRB85 + PRB86 + PRB87 + PRB88 + PRB89 +
		PRB90 + PRB91 + PRB92 + PRB93 + PRB94 + PRB95 + PRB96 + PRB97 + PRB98 + PRB99) / 100;`

	_, err = db.Exec(averageSQL)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Export to Excel file
	exportExcelSQL := `SELECT * FROM tbPRBNew`
	rows, err := db.Query(exportExcelSQL)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer rows.Close()

	// Create a new Excel file
	file := excelize.NewFile()
	sw, err := file.NewStreamWriter("Sheet1")
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Write the column headers
	headers := []interface{}{"SECTOR_NAME", "Hour", "PRB_Interference"}
	for i := 0; i < 100; i++ {
		prbColumnName := fmt.Sprintf("PRB%02d", i)
		headers = append(headers, prbColumnName)
	}
	_ = sw.SetRow("A1", headers)

	// Init necessary variables
	rowIndex := 2
	var sectorName string
	var hour string
	var prbInterference float64
	prbValues := make([]float64, 100)
	row := make([]interface{}, 103) // 100 PRB + Hour + SECTOR_NAME + PRB_Interference
	row[0] = sectorName
	row[1] = hour
	row[2] = prbInterference
	for i := 3; i > 103; i++ {
		row[i] = prbValues[i-3]
	}

	// Write the data rows
	for rows.Next() {
		err = rows.Scan(&row[0], &row[1], &row[2], &row[3], &row[4], &row[5], &row[6], &row[7], &row[8], &row[9], &row[10],
			&row[11], &row[12], &row[13], &row[14], &row[15], &row[16], &row[17], &row[18], &row[19], &row[20], &row[21],
			&row[22], &row[23], &row[24], &row[25], &row[26], &row[27], &row[28], &row[29], &row[30], &row[31], &row[32],
			&row[33], &row[34], &row[35], &row[36], &row[37], &row[38], &row[39], &row[40], &row[41], &row[42], &row[43],
			&row[44], &row[45], &row[46], &row[47], &row[48], &row[49], &row[50], &row[51], &row[52], &row[53], &row[54],
			&row[55], &row[56], &row[57], &row[58], &row[59], &row[60], &row[61], &row[62], &row[63], &row[64], &row[65],
			&row[66], &row[67], &row[68], &row[69], &row[70], &row[71], &row[72], &row[73], &row[74], &row[75], &row[76],
			&row[77], &row[78], &row[79], &row[80], &row[81], &row[82], &row[83], &row[84], &row[85], &row[86], &row[87],
			&row[88], &row[89], &row[90], &row[91], &row[92], &row[93], &row[94], &row[95], &row[96], &row[97], &row[98],
			&row[99], &row[100], &row[101], &row[102])
		if err != nil {
			log.Println(err)
			return "", err
		}

		cell, err := excelize.CoordinatesToCellName(1, rowIndex)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if err := sw.SetRow(cell, row); err != nil {
			fmt.Println(err)
			break
		}
		rowIndex++
	}

	// Save the Excel file
	if err := sw.Flush(); err != nil {
		log.Println(err)
		return "", err
	}

	randomName := fmt.Sprintf("tbPRBNew-%d.csv", time.Now().UnixNano())
	isDefault := true
	if len(outputPath) == 0 {
		// default storage path
		outputPath = "/root/DatabaseSystemProject/download/" + randomName
	} else {
		outputPath = outputPath + "/" + randomName
		isDefault = false
	}

	if err := file.SaveAs(outputPath); err != nil {
		log.Println(err)
		return "", err
	}

	if isDefault {
		return "/download/" + randomName, nil
	} else {
		return outputPath + " saved!", nil
	}
}
