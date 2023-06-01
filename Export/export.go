package Export

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func TableAsCSV(db *sql.DB, outputPath, tableName string) (string, error) {
	randomName := fmt.Sprintf("%s-%d.csv", tableName, time.Now().UnixNano())
	isDefault := true
	if len(outputPath) == 0 {
		// default storage path
		outputPath = "/root/DatabaseSystemProject/download/" + randomName
	} else {
		outputPath = outputPath + "/" + randomName
		isDefault = false
	}

	switch tableName {
	case "tbCell":
		query := fmt.Sprintf(`select 'CITY', 'SECTOR_ID', 'SECTOR_NAME', 'ENODEBID', 'ENODEB_NAME', 'EARFCN', 
       		'PCI', 'PSS', 'SSS', 'TAC', 'VENDOR', 'LONGITUDE', 'LATITUDE', 'STYLE', 'AZIMUTH', 'HEIGHT', 'ELECTTILT',
        	'MECHTILT', 'TOTLETILT'
      		union ALL select * from tbCell
			INTO OUTFILE '%s'
    		FIELDS TERMINATED BY ','
    		OPTIONALLY ENCLOSED BY '"'
    		LINES TERMINATED BY '\n';`, outputPath)

		_, err := db.Exec(query)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "tbEnodeb":
		query := fmt.Sprintf(`select 
			'CITY', 'ENODEBID', 'ENODEB_NAME', 'VENDOR', 'LONGITUDE', 'LATITUDE', 'STYLE'
      		union ALL select * from tbEnodeb
			INTO OUTFILE '%s'
    		FIELDS TERMINATED BY ','
    		OPTIONALLY ENCLOSED BY '"'
    		LINES TERMINATED BY '\n';`, outputPath)

		_, err := db.Exec(query)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "tbKPI":
		query := fmt.Sprintf(`select 
			'StartTime', 'ENODEB_NAME', 'SECTOR_DESCRIPTION', 'SECTOR_NAME', 'RCCConnSUCC', 'RCCConnATT', 
			'RCCConnRATE', 'ERABConnSUCC', 'ERABConnATT', 'ERABConnRATE', 'ENODEB_ERABRel', 'SECTOR_ERABRel', 
			'ERABDropRateNew', 'WirelessAccessRateAY', 'ENODEB_UECtxRel', 'UEContextRel', 'UEContextSUCC',
			'WirelessDropRate', 'ENODEB_InterFreqHOOutSUCC', 'ENODEB_InterFreqHOOutATT', 'ENODEB_IntraFreqHOOutSUCC',
			'ENODEB_IntraFreqHOOutATT', 'ENODEB_InterFreqHOInSUCC', 'ENODEB_InterFreqHOInATT',
			'ENODEB_IntraFreqHOInSUCC', 'ENODEB_IntraFreqHOInATT', 'ENODEB_HOInRate', 'ENODEB_HOOutRate',
			'IntraFreqHOOutRateZSP', 'InterFreqHOOutRateZSP', 'HOSuccessRate', 'PDCP_UplinkThroughput',
			'PDCP_DownlinkThroughput', 'RRCRebuildReq', 'RRCRebuildRate', 'SourceENB_IntraFreqHOOutSUCC',
			'SourceENB_InterFreqHOOutSUCC', 'SourceENB_IntraFreqHOInSUCC', 'SourceENB_InterFreqHOInSUCC',
			'ENODEB_HOOutSUCC', 'ENODEB_HOOutATT'
      		union ALL select * from tbKPI
			INTO OUTFILE '%s'
    		FIELDS TERMINATED BY ','
    		OPTIONALLY ENCLOSED BY '"'
    		LINES TERMINATED BY '\n';`, outputPath)

		_, err := db.Exec(query)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "tbMROData":
		query := fmt.Sprintf(`select 
			'TimeStamp', 'ServingSector', 'InterferingSector', 'LteScRSRP', 'LteNcRSRP', 'LteNcEarfcn', 'LteNcPci'
      		union ALL select * from tbMROData
			INTO OUTFILE '%s'
    		FIELDS TERMINATED BY ','
    		OPTIONALLY ENCLOSED BY '"'
    		LINES TERMINATED BY '\n';`, outputPath)

		_, err := db.Exec(query)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "tbPRB":
		query := fmt.Sprintf(`select 
			'StartTime', 'ENODEB_NAME', 'SECTOR_DESCRIPTION', 'SECTOR_NAME', 'PRB00', 'PRB01', 'PRB02', 'PRB03', 
			'PRB04', 'PRB05', 'PRB06', 'PRB07', 'PRB08', 'PRB09', 'PRB10', 'PRB11', 'PRB12', 'PRB13', 'PRB14',
			'PRB15', 'PRB16', 'PRB17', 'PRB18', 'PRB19', 'PRB20', 'PRB21', 'PRB22', 'PRB23', 'PRB24', 'PRB25',
			'PRB26', 'PRB27', 'PRB28', 'PRB29', 'PRB30', 'PRB31', 'PRB32', 'PRB33', 'PRB34', 'PRB35', 'PRB36',
			'PRB37', 'PRB38', 'PRB39', 'PRB40', 'PRB41', 'PRB42', 'PRB43', 'PRB44', 'PRB45', 'PRB46', 'PRB47',
			'PRB48', 'PRB49', 'PRB50', 'PRB51', 'PRB52', 'PRB53', 'PRB54', 'PRB55', 'PRB56', 'PRB57', 'PRB58',
			'PRB59', 'PRB60', 'PRB61', 'PRB62', 'PRB63', 'PRB64', 'PRB65', 'PRB66', 'PRB67', 'PRB68', 'PRB69',
			'PRB70', 'PRB71', 'PRB72', 'PRB73', 'PRB74', 'PRB75', 'PRB76', 'PRB77', 'PRB78', 'PRB79', 'PRB80',
			'PRB81', 'PRB82', 'PRB83', 'PRB84', 'PRB85', 'PRB86', 'PRB87', 'PRB88', 'PRB89', 'PRB90', 'PRB91',
			'PRB92', 'PRB93', 'PRB94', 'PRB95', 'PRB96', 'PRB97', 'PRB98', 'PRB99'
      		union ALL select * from tbPRB
			INTO OUTFILE '%s'
    		FIELDS TERMINATED BY ','
    		OPTIONALLY ENCLOSED BY '"'
    		LINES TERMINATED BY '\n';`, outputPath)

		_, err := db.Exec(query)
		if err != nil {
			log.Println(err)
			return "", err
		}
	}

	if isDefault {
		return "/download/" + randomName, nil
	} else {
		return randomName, nil
	}
}
