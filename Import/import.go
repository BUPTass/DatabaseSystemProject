package Import

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type CellValues struct {
	City       string
	SectorID   string
	SectorName string
	EnodebID   string
	EnodebName string
	EARFCN     string
	PCI        string
	PSS        int
	SSS        int
	TAC        string
	Vendor     string
	Longitude  string
	Latitude   string
	Style      string
	Azimuth    string
	Height     int
	ElectTilt  int
	MechTilt   int
	TotlTilt   int
}

type KpiValues struct {
	StartTime                    string  `json:"StartTime"`                    // 起始时间
	ENODEB_NAME                  string  `json:"ENODEB_NAME"`                  // 网元/基站名称
	SECTOR_DESCRIPTION           string  `json:"SECTOR_DESCRIPTION"`           // 小区描述
	SECTOR_NAME                  string  `json:"SECTOR_NAME"`                  // 小区名称
	RCCConnSUCC                  int     `json:"RCCConnSUCC"`                  // RCC连接建立完成次数 (无)
	RCCConnATT                   int     `json:"RCCConnATT"`                   // RCC连接请求次数（包括重发） (无)
	RCCConnRATE                  float64 `json:"RCCConnRATE"`                  // RCC建立成功率qf (%) 或者：float 缺省值 null
	ERABConnSUCC                 int     `json:"ERABConnSUCC"`                 // E-RAB建立成功总次数 (无)
	ERABConnATT                  int     `json:"ERABConnATT"`                  // E-RAB建立尝试总次数 (无)
	ERABConnRATE                 float64 `json:"ERABConnRATE"`                 // E-RAB建立成功率2 (%)
	ENODEB_ERABRel               int     `json:"ENODEB_ERABRel"`               // eNodeB触发的E-RAB异常释放总次数 (无)
	SECTOR_ERABRel               int     `json:"SECTOR_ERABRel"`               // 小区切换出E-RAB异常释放总次数 (无)
	ERABDropRateNew              float64 `json:"ERABDropRateNew"`              // E-RAB掉线率(新) (%)
	WirelessAccessRateAY         float64 `json:"WirelessAccessRateAY"`         // 无线接通率ay (%)
	ENODEB_UECtxRel              int     `json:"ENODEB_UECtxRel"`              // eNodeB发起的S1 RESET导致的UE Context释放次数 (无)
	UEContextRel                 int     `json:"UEContextRel"`                 // UE Context异常释放次数 (无)
	UEContextSUCC                int     `json:"UEContextSUCC"`                // UE Context建立成功总次数 (无)
	WirelessDropRate             float64 `json:"WirelessDropRate"`             // 无线掉线率 (%)
	ENODEB_InterFreqHOOutSUCC    int     `json:"ENODEB_InterFreqHOOutSUCC"`    // eNodeB内异频切换出成功次数 (无)
	ENODEB_InterFreqHOOutATT     int     `json:"ENODEB_InterFreqHOOutATT"`     // eNodeB内异频切换出尝试次数 (无)
	ENODEB_IntraFreqHOOutSUCC    int     `json:"ENODEB_IntraFreqHOOutSUCC"`    // eNodeB内同频切换出成功次数 (无)
	ENODEB_IntraFreqHOOutATT     int     `json:"ENODEB_IntraFreqHOOutATT"`     // eNodeB内同频切换出尝试次数 (无)
	ENODEB_InterFreqHOInSUCC     int     `json:"ENODEB_InterFreqHOInSUCC"`     // eNodeB间异频切换出成功次数 (无)
	ENODEB_InterFreqHOInATT      int     `json:"ENODEB_InterFreqHOInATT"`      // eNodeB间异频切换出尝试次数 (无)
	ENODEB_IntraFreqHOInSUCC     int     `json:"ENODEB_IntraFreqHOInSUCC"`     // eNodeB间同频切换出成功次数 (无)
	ENODEB_IntraFreqHOInATT      int     `json:"ENODEB_IntraFreqHOInATT"`      // eNodeB间同频切换出尝试次数 (无)
	ENODEB_HOInRate              float64 `json:"ENODEB_HOInRate"`              // eNB内切换成功率 (%)
	ENODEB_HOOutRate             float64 `json:"ENODEB_HOOutRate"`             // eNB间切换成功率 (%)
	IntraFreqHOOutRateZSP        float64 `json:"IntraFreqHOOutRateZSP"`        // 同频切换成功率zsp (%)
	InterFreqHOOutRateZSP        float64 `json:"InterFreqHOOutRateZSP"`        // 异频切换成功率zsp (%)
	HOSuccessRate                float64 `json:"HOSuccessRate"`                // 切换成功率 (%)
	PDCP_UplinkThroughput        int64   `json:"PDCP_UplinkThroughput"`        // 小区PDCP层所接收到的上行数据的总吞吐量 (比特)
	PDCP_DownlinkThroughput      int64   `json:"PDCP_DownlinkThroughput"`      // 小区PDCP层所发送的下行数据的总吞吐量 (比特)
	RRCRebuildReq                int     `json:"RRCRebuildReq"`                // RRC重建请求次数 (无)
	RRCRebuildRate               float64 `json:"RRCRebuildRate"`               // RRC连接重建比率 (%)
	SourceENB_IntraFreqHOOutSUCC int     `json:"SourceENB_IntraFreqHOOutSUCC"` // 通过重建回源小区的eNodeB间同频切换出执行成功次数 (无)
	SourceENB_InterFreqHOOutSUCC int     `json:"SourceENB_InterFreqHOOutSUCC"` // 通过重建回源小区的eNodeB间异频切换出执行成功次数 (无)
	SourceENB_IntraFreqHOInSUCC  int     `json:"SourceENB_IntraFreqHOInSUCC"`  // 通过重建回源小区的eNodeB内同频切换出执行成功次数 (无)
	SourceENB_InterFreqHOInSUCC  int     `json:"SourceENB_InterFreqHOInSUCC"`  // 通过重建回源小区的eNodeB内异频切换出执行成功次数 (无)
	ENODEB_HOOutSUCC             int     `json:"ENODEB_HOOutSUCC"`             // eNB内切换出成功次数 (次)
	ENODEB_HOOutATT              int     `json:"ENODEB_HOOutATT"`              // eNB内切换出请求次数 (次)
}

type MROData struct {
	TimeStamp         string  `json:"TimeStamp"`         // 测量时间点
	ServingSector     string  `json:"ServingSector"`     // 服务小区/主小区ID
	InterferingSector string  `json:"InterferingSector"` // 干扰小区ID
	LteScRSRP         float64 `json:"LteScRSRP"`         // 服务小区参考信号接收功率RSRP
	LteNcRSRP         float64 `json:"LteNcRSRP"`         // 干扰小区参考信号接收功率RSRP
	LteNcEarfcn       int     `json:"LteNcEarfcn"`       // 干扰小区频点
	LteNcPci          int16   `json:"LteNcPci"`          // 干扰小区PCI
}

type PRBData struct {
	StartTime          string `json:"StartTime"`
	ENODEB_NAME        string `json:"ENODEB_NAME"`
	SECTOR_DESCRIPTION string `json:"SECTOR_DESCRIPTION"`
	SECTOR_NAME        string `json:"SECTOR_NAME"`
	PRB00              int    `json:"PRB00"`
	PRB01              int    `json:"PRB01"`
	PRB02              int    `json:"PRB02"`
	PRB03              int    `json:"PRB03"`
	PRB04              int    `json:"PRB04"`
	PRB05              int    `json:"PRB05"`
	PRB06              int    `json:"PRB06"`
	PRB07              int    `json:"PRB07"`
	PRB08              int    `json:"PRB08"`
	PRB09              int    `json:"PRB09"`
	PRB10              int    `json:"PRB10"`
	PRB11              int    `json:"PRB11"`
	PRB12              int    `json:"PRB12"`
	PRB13              int    `json:"PRB13"`
	PRB14              int    `json:"PRB14"`
	PRB15              int    `json:"PRB15"`
	PRB16              int    `json:"PRB16"`
	PRB17              int    `json:"PRB17"`
	PRB18              int    `json:"PRB18"`
	PRB19              int    `json:"PRB19"`
	PRB20              int    `json:"PRB20"`
	PRB21              int    `json:"PRB21"`
	PRB22              int    `json:"PRB22"`
	PRB23              int    `json:"PRB23"`
	PRB24              int    `json:"PRB24"`
	PRB25              int    `json:"PRB25"`
	PRB26              int    `json:"PRB26"`
	PRB27              int    `json:"PRB27"`
	PRB28              int    `json:"PRB28"`
	PRB29              int    `json:"PRB29"`
	PRB30              int    `json:"PRB30"`
	PRB31              int    `json:"PRB31"`
	PRB32              int    `json:"PRB32"`
	PRB33              int    `json:"PRB33"`
	PRB34              int    `json:"PRB34"`
	PRB35              int    `json:"PRB35"`
	PRB36              int    `json:"PRB36"`
	PRB37              int    `json:"PRB37"`
	PRB38              int    `json:"PRB38"`
	PRB39              int    `json:"PRB39"`
	PRB40              int    `json:"PRB40"`
	PRB41              int    `json:"PRB41"`
	PRB42              int    `json:"PRB42"`
	PRB43              int    `json:"PRB43"`
	PRB44              int    `json:"PRB44"`
	PRB45              int    `json:"PRB45"`
	PRB46              int    `json:"PRB46"`
	PRB47              int    `json:"PRB47"`
	PRB48              int    `json:"PRB48"`
	PRB49              int    `json:"PRB49"`
	PRB50              int    `json:"PRB50"`
	PRB51              int    `json:"PRB51"`
	PRB52              int    `json:"PRB52"`
	PRB53              int    `json:"PRB53"`
	PRB54              int    `json:"PRB54"`
	PRB55              int    `json:"PRB55"`
	PRB56              int    `json:"PRB56"`
	PRB57              int    `json:"PRB57"`
	PRB58              int    `json:"PRB58"`
	PRB59              int    `json:"PRB59"`
	PRB60              int    `json:"PRB60"`
	PRB61              int    `json:"PRB61"`
	PRB62              int    `json:"PRB62"`
	PRB63              int    `json:"PRB63"`
	PRB64              int    `json:"PRB64"`
	PRB65              int    `json:"PRB65"`
	PRB66              int    `json:"PRB66"`
	PRB67              int    `json:"PRB67"`
	PRB68              int    `json:"PRB68"`
	PRB69              int    `json:"PRB69"`
	PRB70              int    `json:"PRB70"`
	PRB71              int    `json:"PRB71"`
	PRB72              int    `json:"PRB72"`
	PRB73              int    `json:"PRB73"`
	PRB74              int    `json:"PRB74"`
	PRB75              int    `json:"PRB75"`
	PRB76              int    `json:"PRB76"`
	PRB77              int    `json:"PRB77"`
	PRB78              int    `json:"PRB78"`
	PRB79              int    `json:"PRB79"`
	PRB80              int    `json:"PRB80"`
	PRB81              int    `json:"PRB81"`
	PRB82              int    `json:"PRB82"`
	PRB83              int    `json:"PRB83"`
	PRB84              int    `json:"PRB84"`
	PRB85              int    `json:"PRB85"`
	PRB86              int    `json:"PRB86"`
	PRB87              int    `json:"PRB87"`
	PRB88              int    `json:"PRB88"`
	PRB89              int    `json:"PRB89"`
	PRB90              int    `json:"PRB90"`
	PRB91              int    `json:"PRB91"`
	PRB92              int    `json:"PRB92"`
	PRB93              int    `json:"PRB93"`
	PRB94              int    `json:"PRB94"`
	PRB95              int    `json:"PRB95"`
	PRB96              int    `json:"PRB96"`
	PRB97              int    `json:"PRB97"`
	PRB98              int    `json:"PRB98"`
	PRB99              int    `json:"PRB99"`
}

func AddtbCell(db *sql.DB, path string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}

	var errorList []string
	executed := false
	var stmt *sql.Stmt

	// Prepare the batch insertion statement
	statementPre := `
        INSERT INTO tbCell (City, Sector_ID, Sector_Name, EnodebID, Enodeb_Name, EARFCN, PCI, PSS, SSS, TAC,
                            Vendor, Longitude, Latitude, Style, Azimuth, Height, ElectTilt, MechTilt, TOTLETILT)
        VALUES %s
        ON DUPLICATE KEY UPDATE
        City = VALUES(City),
        Sector_ID = VALUES(Sector_ID),
        Sector_Name = VALUES(Sector_Name),
        EnodebID = VALUES(EnodebID),
        Enodeb_Name = VALUES(Enodeb_Name),
        EARFCN = VALUES(EARFCN),
        PCI = VALUES(PCI),
        PSS = VALUES(PSS),
        SSS = VALUES(SSS),
        TAC = VALUES(TAC),
        Vendor = VALUES(Vendor),
        Longitude = VALUES(Longitude),
        Latitude = VALUES(Latitude),
        Style = VALUES(Style),
        Azimuth = VALUES(Azimuth),
        Height = VALUES(Height),
        ElectTilt = VALUES(ElectTilt),
        MechTilt = VALUES(MechTilt),
        TOTLETILT = VALUES(TOTLETILT)
    `

	count := 0
	// Define the batch size
	batchSize := 50
	values := make([]interface{}, 0, batchSize*19) // 19 is the number of columns in the table

	// Read only the 1st sheet
	sheet := f.GetSheetName(0)
	rows, err := f.Rows(sheet)
	if err != nil {
		fmt.Println(err)
		return err
	}
	valueStrings := make([]string, 0)

	// Skip the first scheme line
	rows.Next()

	for rows.Next() {
		count++
		row, err := rows.Columns()
		if err != nil || len(row) != 19 {
			newErr := "No sufficient columns when importing entry: " + strconv.Itoa(count) + " in " + path
			errorList = append(errorList, newErr)
			log.Println(newErr)
			for i := len(row); i < 19; i++ {
				row = append(row, "")
			}
		}
		// Extract the cell values from the row.

		PCI, err := strconv.Atoi(row[6])
		if err != nil {
			newErr := "Error when importing entry: " + strconv.Itoa(count) + " in " + path
			errorList = append(errorList, newErr)
			log.Println(newErr)
			continue
		}
		ElectTilt, err := strconv.ParseFloat(row[16], 64)
		MechTilt, err := strconv.ParseFloat(row[17], 64)
		TotlTilt, err := strconv.ParseFloat(row[18], 64)

		// Data Cleaning
		if TotlTilt == 0 || TotlTilt != ElectTilt+MechTilt {
			TotlTilt = ElectTilt + MechTilt
		}

		LONGITUDE, err := strconv.ParseFloat(row[11], 64)
		if err != nil || LONGITUDE > 180 || LONGITUDE < -180 {
			newErr := "Error when importing entry: " + strconv.Itoa(count) + " in " + path
			errorList = append(errorList, newErr)
			log.Println(newErr)
			continue
		}

		LATITUDE, err := strconv.ParseFloat(row[12], 64)
		if err != nil || LATITUDE > 90 || LATITUDE < -90 {
			newErr := "Error when importing entry: " + strconv.Itoa(count) + " in " + path
			errorList = append(errorList, newErr)
			log.Println(newErr)
			continue
		}

		cellValues := []interface{}{
			row[0],
			row[1],
			row[2],
			row[3],
			row[4],
			row[5],
			row[6],
			PCI % 3,
			PCI / 3,
			row[9],
			row[10],
			row[11],
			row[12],
			row[13],
			row[14],
			row[15],
			ElectTilt,
			MechTilt,
			TotlTilt,
		}

		// Batch insertion
		if !executed {
			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		}

		// Append the cell values to the batch
		values = append(values, cellValues...)

		// If the batch size is reached, execute the batch insertion
		if len(values) == batchSize*19 {
			if !executed {
				statement := fmt.Sprintf(statementPre, strings.Join(valueStrings, ","))
				stmt, err = db.Prepare(statement)
				if err != nil {
					log.Println(err)
					return err
				}
				defer stmt.Close()
				executed = true
			}
			_, err = stmt.Exec(values...)
			if err != nil {
				log.Println("Error executing batch insertion:", err)
				errorList = append(errorList, err.Error())
				// Fallback to single insertion when bulk insertion failed
				errorList = singleInsertion(db, values, statementPre, 19, errorList)
			}
			values = values[:0] // Clear the batch
		}
	}

	if err = rows.Close(); err != nil {
		fmt.Println(err)
	}
	if err = f.Close(); err != nil {
		fmt.Println(err)
	}

	// Insert the remaining values in the batch
	errorList = singleInsertion(db, values, statementPre, 19, errorList)

	if len(errorList) > 0 {
		return errors.New(strings.Join(errorList, "\n"))
	} else {
		return nil
	}
}

func AddtbKPI(db *sql.DB, path string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}

	var errorList []string
	executed := false
	var stmt *sql.Stmt

	// Prepare the batch insertion statement
	statementPre := `
    INSERT INTO tbKPI (StartTime, ENODEB_NAME, SECTOR_DESCRIPTION, SECTOR_NAME, RCCConnSUCC, RCCConnATT,
                       RCCConnRATE, ERABConnSUCC, ERABConnATT, ERABConnRATE, ENODEB_ERABRel, SECTOR_ERABRel,
                       ERABDropRateNew, WirelessAccessRateAY, ENODEB_UECtxRel, UEContextRel, UEContextSUCC,
                       WirelessDropRate, ENODEB_InterFreqHOOutSUCC, ENODEB_InterFreqHOOutATT,
                       ENODEB_IntraFreqHOOutSUCC, ENODEB_IntraFreqHOOutATT, ENODEB_InterFreqHOInSUCC,
                       ENODEB_InterFreqHOInATT, ENODEB_IntraFreqHOInSUCC, ENODEB_IntraFreqHOInATT,
                       ENODEB_HOInRate, ENODEB_HOOutRate, IntraFreqHOOutRateZSP, InterFreqHOOutRateZSP,
                       HOSuccessRate, PDCP_UplinkThroughput, PDCP_DownlinkThroughput, RRCRebuildReq,
                       RRCRebuildRate, SourceENB_IntraFreqHOOutSUCC, SourceENB_InterFreqHOOutSUCC,
                       SourceENB_IntraFreqHOInSUCC, SourceENB_InterFreqHOInSUCC, ENODEB_HOOutSUCC, ENODEB_HOOutATT)
    VALUES %s
    ON DUPLICATE KEY UPDATE
    StartTime = VALUES(StartTime),
    ENODEB_NAME = VALUES(ENODEB_NAME),
    SECTOR_DESCRIPTION = VALUES(SECTOR_DESCRIPTION),
    RCCConnSUCC = VALUES(RCCConnSUCC),
    RCCConnATT = VALUES(RCCConnATT),
    RCCConnRATE = VALUES(RCCConnRATE),
    ERABConnSUCC = VALUES(ERABConnSUCC),
    ERABConnATT = VALUES(ERABConnATT),
    ERABConnRATE = VALUES(ERABConnRATE),
    ENODEB_ERABRel = VALUES(ENODEB_ERABRel),
    SECTOR_ERABRel = VALUES(SECTOR_ERABRel),
    ERABDropRateNew = VALUES(ERABDropRateNew),
    WirelessAccessRateAY = VALUES(WirelessAccessRateAY),
    ENODEB_UECtxRel = VALUES(ENODEB_UECtxRel),
    UEContextRel = VALUES(UEContextRel),
    UEContextSUCC = VALUES(UEContextSUCC),
    WirelessDropRate = VALUES(WirelessDropRate),
    ENODEB_InterFreqHOOutSUCC = VALUES(ENODEB_InterFreqHOOutSUCC),
    ENODEB_InterFreqHOOutATT = VALUES(ENODEB_InterFreqHOOutATT),
    ENODEB_IntraFreqHOOutSUCC = VALUES(ENODEB_IntraFreqHOOutSUCC),
    ENODEB_IntraFreqHOOutATT = VALUES(ENODEB_IntraFreqHOOutATT),
    ENODEB_InterFreqHOInSUCC = VALUES(ENODEB_InterFreqHOInSUCC),
    ENODEB_InterFreqHOInATT = VALUES(ENODEB_InterFreqHOInATT),
    ENODEB_IntraFreqHOInSUCC = VALUES(ENODEB_IntraFreqHOInSUCC),
    ENODEB_IntraFreqHOInATT = VALUES(ENODEB_IntraFreqHOInATT),
    ENODEB_HOInRate = VALUES(ENODEB_HOInRate),
    ENODEB_HOOutRate = VALUES(ENODEB_HOOutRate),
    IntraFreqHOOutRateZSP = VALUES(IntraFreqHOOutRateZSP),
    InterFreqHOOutRateZSP = VALUES(InterFreqHOOutRateZSP),
    HOSuccessRate = VALUES(HOSuccessRate),
    PDCP_UplinkThroughput = VALUES(PDCP_UplinkThroughput),
    PDCP_DownlinkThroughput = VALUES(PDCP_DownlinkThroughput),
    RRCRebuildReq = VALUES(RRCRebuildReq),
    RRCRebuildRate = VALUES(RRCRebuildRate),
    SourceENB_IntraFreqHOOutSUCC = VALUES(SourceENB_IntraFreqHOOutSUCC),
    SourceENB_InterFreqHOOutSUCC = VALUES(SourceENB_InterFreqHOOutSUCC),
    SourceENB_IntraFreqHOInSUCC = VALUES(SourceENB_IntraFreqHOInSUCC),
    SourceENB_InterFreqHOInSUCC = VALUES(SourceENB_InterFreqHOInSUCC),
    ENODEB_HOOutSUCC = VALUES(ENODEB_HOOutSUCC),
    ENODEB_HOOutATT = VALUES(ENODEB_HOOutATT)
`

	count := 0
	// Define the batch size
	batchSize := 30
	values := make([]interface{}, 0, batchSize*41) // 41 is the number of columns in the table

	// Read only the 1st sheet
	sheet := f.GetSheetName(0)
	rows, err := f.Rows(sheet)
	if err != nil {
		fmt.Println(err)
		return err
	}
	valueStrings := make([]string, 0)

	// Skip the first scheme line
	rows.Next()

	for rows.Next() {
		count++
		row, err := rows.Columns()
		if err != nil || len(row) != 41 {
			newErr := "No sufficient columns when importing entry: " + strconv.Itoa(count) + " in " + path
			errorList = append(errorList, newErr)
			log.Println(newErr)
			for i := len(row); i < 41; i++ {
				row = append(row, "")
			}
		}

		layout := "01/02/2006 15:04:05" // The layout represents the format of the input string
		parsedTime, err := time.Parse(layout, row[0])
		if err != nil {
			log.Println(err)
			errorList = append(errorList, err.Error())
			continue
		}

		// Extract the cell values from the row.
		cellValues := []interface{}{
			parsedTime,
			row[1],
			row[2],
			row[3],
			row[4],
			row[5],
			row[6],
			row[7],
			row[8],
			row[9],
			row[10],
			row[11],
			row[12],
			row[13],
			row[14],
			row[15],
			row[16],
			row[17],
			row[18],
			row[19],
			row[20],
			row[21],
			row[22],
			row[23],
			row[24],
			row[25],
			row[26],
			row[27],
			row[28],
			row[29],
			row[30],
			row[31],
			row[32],
			row[33],
			row[34],
			row[35],
			row[36],
			row[37],
			row[38],
			row[39],
			row[40],
		}

		// Batch insertion
		if !executed {
			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, "+
				"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		}

		// Append the cell values to the batch
		values = append(values, cellValues...)

		// If the batch size is reached, execute the batch insertion
		if len(values) == batchSize*41 {
			if !executed {
				statement := fmt.Sprintf(statementPre, strings.Join(valueStrings, ","))
				stmt, err = db.Prepare(statement)
				if err != nil {
					log.Println(err)
					return err
				}
				defer stmt.Close()
				executed = true
			}
			_, err = stmt.Exec(values...)
			if err != nil {
				log.Println("Error executing batch insertion:", err)
				errorList = append(errorList, err.Error())
				// Fallback to single insertion when bulk insertion failed
				errorList = singleInsertion(db, values, statementPre, 41, errorList)
			}
			values = values[:0] // Clear the batch
		}
	}

	if err = rows.Close(); err != nil {
		fmt.Println(err)
	}
	if err = f.Close(); err != nil {
		fmt.Println(err)
	}

	// Insert the remaining values in the batch
	errorList = singleInsertion(db, values, statementPre, 41, errorList)

	if len(errorList) > 0 {
		return errors.New(strings.Join(errorList, "\n"))
	} else {
		return nil
	}
}

func AddtbPRB(db *sql.DB, path string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}

	var errorList []string
	executed := false
	var stmt *sql.Stmt

	// Prepare the batch insertion statement
	statementPre := `
	INSERT INTO tbPRB (StartTime, ENODEB_NAME, SECTOR_DESCRIPTION, SECTOR_NAME,
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
	VALUES %s
	ON DUPLICATE KEY UPDATE
		StartTime = VALUES(StartTime),
		ENODEB_NAME = VALUES(ENODEB_NAME),
		SECTOR_DESCRIPTION = VALUES(SECTOR_DESCRIPTION),
		SECTOR_NAME = VALUES(SECTOR_NAME),
		PRB00 = VALUES(PRB00),
		PRB01 = VALUES(PRB01),
		PRB02 = VALUES(PRB02),
		PRB03 = VALUES(PRB03),
		PRB04 = VALUES(PRB04),
		PRB05 = VALUES(PRB05),
		PRB06 = VALUES(PRB06),
		PRB07 = VALUES(PRB07),
		PRB08 = VALUES(PRB08),
		PRB09 = VALUES(PRB09),
		PRB10 = VALUES(PRB10),
		PRB11 = VALUES(PRB11),
		PRB12 = VALUES(PRB12),
		PRB13 = VALUES(PRB13),
		PRB14 = VALUES(PRB14),
		PRB15 = VALUES(PRB15),
		PRB16 = VALUES(PRB16),
		PRB17 = VALUES(PRB17),
		PRB18 = VALUES(PRB18),
		PRB19 = VALUES(PRB19),
		PRB20 = VALUES(PRB20),
		PRB21 = VALUES(PRB21),
		PRB22 = VALUES(PRB22),
		PRB23 = VALUES(PRB23),
		PRB24 = VALUES(PRB24),
		PRB25 = VALUES(PRB25),
		PRB26 = VALUES(PRB26),
		PRB27 = VALUES(PRB27),
		PRB28 = VALUES(PRB28),
		PRB29 = VALUES(PRB29),
		PRB30 = VALUES(PRB30),
		PRB31 = VALUES(PRB31),
		PRB32 = VALUES(PRB32),
		PRB33 = VALUES(PRB33),
		PRB34 = VALUES(PRB34),
		PRB35 = VALUES(PRB35),
		PRB36 = VALUES(PRB36),
		PRB37 = VALUES(PRB37),
		PRB38 = VALUES(PRB38),
		PRB39 = VALUES(PRB39),
		PRB40 = VALUES(PRB40),
		PRB41 = VALUES(PRB41),
		PRB42 = VALUES(PRB42),
		PRB43 = VALUES(PRB43),
		PRB44 = VALUES(PRB44),
		PRB45 = VALUES(PRB45),
		PRB46 = VALUES(PRB46),
		PRB47 = VALUES(PRB47),
		PRB48 = VALUES(PRB48),
		PRB49 = VALUES(PRB49),
		PRB50 = VALUES(PRB50),
		PRB51 = VALUES(PRB51),
		PRB52 = VALUES(PRB52),
		PRB53 = VALUES(PRB53),
		PRB54 = VALUES(PRB54),
		PRB55 = VALUES(PRB55),
		PRB56 = VALUES(PRB56),
		PRB57 = VALUES(PRB57),
		PRB58 = VALUES(PRB58),
		PRB59 = VALUES(PRB59),
		PRB60 = VALUES(PRB60),
		PRB61 = VALUES(PRB61),
		PRB62 = VALUES(PRB62),
		PRB63 = VALUES(PRB63),
		PRB64 = VALUES(PRB64),
		PRB65 = VALUES(PRB65),
		PRB66 = VALUES(PRB66),
		PRB67 = VALUES(PRB67),
		PRB68 = VALUES(PRB68),
		PRB69 = VALUES(PRB69),
		PRB70 = VALUES(PRB70),
		PRB71 = VALUES(PRB71),
		PRB72 = VALUES(PRB72),
		PRB73 = VALUES(PRB73),
		PRB74 = VALUES(PRB74),
		PRB75 = VALUES(PRB75),
		PRB76 = VALUES(PRB76),
		PRB77 = VALUES(PRB77),
		PRB78 = VALUES(PRB78),
		PRB79 = VALUES(PRB79),
		PRB80 = VALUES(PRB80),
		PRB81 = VALUES(PRB81),
		PRB82 = VALUES(PRB82),
		PRB83 = VALUES(PRB83),
		PRB84 = VALUES(PRB84),
		PRB85 = VALUES(PRB85),
		PRB86 = VALUES(PRB86),
		PRB87 = VALUES(PRB87),
		PRB88 = VALUES(PRB88),
		PRB89 = VALUES(PRB89),
		PRB90 = VALUES(PRB90),
		PRB91 = VALUES(PRB91),
		PRB92 = VALUES(PRB92),
		PRB93 = VALUES(PRB93),
		PRB94 = VALUES(PRB94),
		PRB95 = VALUES(PRB95),
		PRB96 = VALUES(PRB96),
		PRB97 = VALUES(PRB97),
		PRB98 = VALUES(PRB98),
		PRB99 = VALUES(PRB99)
`

	count := 0
	// Define the batch size
	batchSize := 50
	values := make([]interface{}, 0, batchSize*104) // 104 is the number of columns in the table

	// Read only the 1st sheet
	sheet := f.GetSheetName(0)
	rows, err := f.Rows(sheet)
	if err != nil {
		fmt.Println(err)
		return err
	}
	valueStrings := make([]string, 0)

	// Skip the first scheme line
	rows.Next()

	for rows.Next() {
		count++
		row, err := rows.Columns()
		if err != nil || len(row) != 104 {
			newErr := "No sufficient columns when importing entry: " + strconv.Itoa(count) + " in " + path
			errorList = append(errorList, newErr)
			log.Println(newErr)
			for i := len(row); i < 104; i++ {
				row = append(row, "")
			}
		}
		layout := "01/02/2006 15:04:05" // The layout represents the format of the input string
		parsedTime, err := time.Parse(layout, row[0])
		if err != nil {
			log.Println(err)
			errorList = append(errorList, err.Error())
			continue
		}
		// Extract the cell values from the row.
		cellValues := []interface{}{
			parsedTime,
			row[1],
			row[2],
			row[3],
			row[4],
			row[5],
			row[6],
			row[7],
			row[8],
			row[9],
			row[10],
			row[11],
			row[12],
			row[13],
			row[14],
			row[15],
			row[16],
			row[17],
			row[18],
			row[19],
			row[20],
			row[21],
			row[22],
			row[23],
			row[24],
			row[25],
			row[26],
			row[27],
			row[28],
			row[29],
			row[30],
			row[31],
			row[32],
			row[33],
			row[34],
			row[35],
			row[36],
			row[37],
			row[38],
			row[39],
			row[40],
			row[41],
			row[42],
			row[43],
			row[44],
			row[45],
			row[46],
			row[47],
			row[48],
			row[49],
			row[50],
			row[51],
			row[52],
			row[53],
			row[54],
			row[55],
			row[56],
			row[57],
			row[58],
			row[59],
			row[60],
			row[61],
			row[62],
			row[63],
			row[64],
			row[65],
			row[66],
			row[67],
			row[68],
			row[69],
			row[70],
			row[71],
			row[72],
			row[73],
			row[74],
			row[75],
			row[76],
			row[77],
			row[78],
			row[79],
			row[80],
			row[81],
			row[82],
			row[83],
			row[84],
			row[85],
			row[86],
			row[87],
			row[88],
			row[89],
			row[90],
			row[91],
			row[92],
			row[93],
			row[94],
			row[95],
			row[96],
			row[97],
			row[98],
			row[99],
			row[100],
			row[101],
			row[102],
			row[103],
		}

		// Batch insertion
		if !executed {
			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,"+
				"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,"+
				"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,"+
				"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,"+
				"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		}

		// Append the cell values to the batch
		values = append(values, cellValues...)

		// If the batch size is reached, execute the batch insertion
		if len(values) == batchSize*104 {
			if !executed {
				statement := fmt.Sprintf(statementPre, strings.Join(valueStrings, ","))
				stmt, err = db.Prepare(statement)
				if err != nil {
					log.Println(err)
					return err
				}
				defer stmt.Close()
				executed = true
			}
			_, err = stmt.Exec(values...)
			if err != nil {
				log.Println("Error executing batch insertion:", err)
				errorList = append(errorList, err.Error())
				// Fallback to single insertion when bulk insertion failed
				errorList = singleInsertion(db, values, statementPre, 104, errorList)
			}
			values = values[:0] // Clear the batch
		}
	}

	if err = rows.Close(); err != nil {
		fmt.Println(err)
	}
	if err = f.Close(); err != nil {
		fmt.Println(err)
	}

	// Insert the remaining values in the batch
	errorList = singleInsertion(db, values, statementPre, 104, errorList)

	if len(errorList) > 0 {
		return errors.New(strings.Join(errorList, "\n"))
	} else {
		return nil
	}
}

func AddtbMROData(db *sql.DB, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	ch := processCSV(file)

	var errorList []string
	executed := false
	var stmt *sql.Stmt

	// Prepare the batch insertion statement
	statementPre := `
    INSERT INTO tbMROData (TimeStamp, ServingSector, InterferingSector, LteScRSRP, LteNcRSRP, LteNcEarfcn, LteNcPci)
    VALUES %s
    ON DUPLICATE KEY UPDATE
    TimeStamp = VALUES(TimeStamp),
    ServingSector = VALUES(ServingSector),
    InterferingSector = VALUES(InterferingSector),
    LteScRSRP = VALUES(LteScRSRP),
    LteNcRSRP = VALUES(LteNcRSRP),
    LteNcEarfcn = VALUES(LteNcEarfcn),
    LteNcPci = VALUES(LteNcPci)
`

	count := 0
	// Define the batch size
	batchSize := 50
	values := make([]interface{}, 0, batchSize*7) // 7 is the number of columns in the table
	valueStrings := make([]string, 0)

	for row := range ch {
		count++
		if len(row) != 7 {
			newErr := "No sufficient columns when importing entry: " + strconv.Itoa(count) + " in " + path
			errorList = append(errorList, newErr)
			log.Println(newErr)
			for i := len(row); i < 7; i++ {
				row = append(row, "")
			}
		}
		// Extract the cell values from the row.
		cellValues := []interface{}{
			row[0],
			row[1],
			row[2],
			row[3],
			row[4],
			row[5],
			row[6],
		}

		// Batch insertion
		if !executed {
			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?)")
		}

		// Append the cell values to the batch
		values = append(values, cellValues...)

		// If the batch size is reached, execute the batch insertion
		if len(values) == batchSize*7 {
			if !executed {
				statement := fmt.Sprintf(statementPre, strings.Join(valueStrings, ","))
				stmt, err = db.Prepare(statement)
				if err != nil {
					log.Println(err)
					return err
				}
				defer stmt.Close()
				executed = true
			}
			_, err = stmt.Exec(values...)
			if err != nil {
				log.Println("Error executing batch insertion:", err)
				errorList = append(errorList, err.Error())
				// Fallback to single insertion when bulk insertion failed
				errorList = singleInsertion(db, values, statementPre, 7, errorList)
			}
			values = values[:0] // Clear the batch
		}
	}

	if err = file.Close(); err != nil {
		fmt.Println(err)
	}

	// Insert the remaining values in the batch
	errorList = singleInsertion(db, values, statementPre, 7, errorList)

	if len(errorList) > 0 {
		return errors.New(strings.Join(errorList, "\n"))
	} else {
		return nil
	}
}

func processCSV(rc io.Reader) (ch chan []string) {
	// 10 channels in total
	// Approximately 10-threaded stream reading
	ch = make(chan []string, 10)
	go func() {
		r := csv.NewReader(rc)
		if _, err := r.Read(); err != nil { // Read header and ignore
			log.Println(err)
		}
		defer close(ch)
		for {
			rec, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println(err)
			}
			ch <- rec
		}
	}()
	return
}

func singleInsertion(db *sql.DB, values []interface{}, statementPre string, pCount int, errorList []string) []string {
	if len(values) > 0 {
		statement := fmt.Sprintf(statementPre, generatePattern(pCount))
		stmt, err := db.Prepare(statement)
		if err != nil {
			log.Println(err)
			errorList = append(errorList, err.Error())
			return errorList
		}
		defer stmt.Close()

		for len(values) > 0 {
			_, err = stmt.Exec(values[0:pCount]...)
			if err != nil {
				log.Println("Error executing insertion:", err)
				errorList = append(errorList, err.Error())
			}
			values = values[pCount:]
		}
	}
	return errorList
}

func generatePattern(pCount int) string {
	// Create a slice of "?" strings with length pCount
	placeholders := make([]string, pCount)
	for i := 0; i < pCount; i++ {
		placeholders[i] = "?"
	}

	// Join the "?" strings with ", " separator
	pattern := strings.Join(placeholders, ", ")

	// Enclose the pattern with parentheses
	pattern = "(" + pattern + ")"

	return pattern
}

func AddtbC2I(db *sql.DB, path string) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}

	var errorList []string
	executed := false
	var stmt *sql.Stmt

	// Prepare the batch insertion statement
	statementPre := `INSERT INTO tbC2I (CITY, SCELL, NCELL, PrC2I9, C2I_Mean, Std, SampleCount, WeightedC2I) 
		VALUES %s ON DUPLICATE KEY UPDATE
		CITY = VALUES(CITY), SCELL = VALUES(SCELL), NCELL = VALUES(NCELL), PrC2I9 = VALUES(PrC2I9),
		C2I_Mean = VALUES(C2I_Mean), Std = VALUES(Std), SampleCount = VALUES(SampleCount),
		WeightedC2I = VALUES(WeightedC2I)`

	count := 0
	// Define the batch size
	batchSize := 50
	values := make([]interface{}, 0, batchSize*8)

	// Read only the 1st sheet
	sheet := f.GetSheetName(0)
	rows, err := f.Rows(sheet)
	if err != nil {
		fmt.Println(err)
		return err
	}
	valueStrings := make([]string, 0)

	// Skip the first scheme line
	rows.Next()

	for rows.Next() {
		count++
		row, err := rows.Columns()
		if err != nil || len(row) != 8 {
			newErr := "No sufficient columns when importing entry: " + strconv.Itoa(count) + " in " + path
			errorList = append(errorList, newErr)
			log.Println(newErr)
			for i := len(row); i < 8; i++ {
				row = append(row, "")
			}
		}
		// Extract the cell values from the row.

		cellValues := []interface{}{
			row[0],
			row[1],
			row[2],
			row[3],
			row[4],
			row[5],
			row[6],
			row[7],
		}

		// Batch insertion
		if !executed {
			valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?)")
		}

		// Append the cell values to the batch
		values = append(values, cellValues...)

		// If the batch size is reached, execute the batch insertion
		if len(values) == batchSize*19 {
			if !executed {
				statement := fmt.Sprintf(statementPre, strings.Join(valueStrings, ","))
				stmt, err = db.Prepare(statement)
				if err != nil {
					log.Println(err)
					return err
				}
				defer stmt.Close()
				executed = true
			}
			_, err = stmt.Exec(values...)
			if err != nil {
				log.Println("Error executing batch insertion:", err)
				errorList = append(errorList, err.Error())
				// Fallback to single insertion when bulk insertion failed
				errorList = singleInsertion(db, values, statementPre, 8, errorList)
			}
			values = values[:0] // Clear the batch
		}
	}

	if err = rows.Close(); err != nil {
		fmt.Println(err)
	}
	if err = f.Close(); err != nil {
		fmt.Println(err)
	}

	// Insert the remaining values in the batch
	errorList = singleInsertion(db, values, statementPre, 8, errorList)

	if len(errorList) > 0 {
		return errors.New(strings.Join(errorList, "\n"))
	} else {
		return nil
	}
}
