package Import

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tealeg/xlsx"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
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

func AddtbCell(EsiFile *multipart.FileHeader) error {
	tmpFile, err := EsiFile.Open()
	if err != nil {
		return err
	}
	file, err := xlsx.OpenReaderAt(tmpFile, EsiFile.Size)
	if err != nil {
		return err
	}
	// Connect to the MySQL database.
	db, err := sql.Open("mysql", "root:1taNWY1vXdTc4_-j@tcp(127.0.0.1:3306)/LTE")
	if err != nil {
		log.Println("Error connecting to database:", err)
		return err
	}
	defer db.Close()

	var errorList []string

	// Iterate over each sheet in the xlsx file.
	for _, sheet := range file.Sheets {
		// Iterate over each row in the sheet. Ignore the first row.
		for _, row := range sheet.Rows[1:] {
			// Extract the cell values from the row.
			PCI, _ := strconv.Atoi(row.Cells[6].Value)
			HEIGHT, _ := strconv.Atoi(row.Cells[15].Value)
			ElectTilt, _ := strconv.Atoi(row.Cells[16].Value)
			MechTilt, _ := strconv.Atoi(row.Cells[17].Value)
			TotlTilt, _ := strconv.Atoi(row.Cells[18].Value)
			cellValues := CellValues{
				City:       row.Cells[0].Value,
				SectorID:   row.Cells[1].Value,
				SectorName: row.Cells[2].Value,
				EnodebID:   row.Cells[3].Value,
				EnodebName: row.Cells[4].Value,
				EARFCN:     row.Cells[5].Value,
				PCI:        row.Cells[6].Value,
				PSS:        PCI % 3,
				SSS:        PCI / 3,
				TAC:        row.Cells[9].Value,
				Vendor:     row.Cells[10].Value,
				Longitude:  row.Cells[11].Value,
				Latitude:   row.Cells[12].Value,
				Style:      row.Cells[13].Value,
				Azimuth:    row.Cells[14].Value,
				Height:     HEIGHT,
				ElectTilt:  ElectTilt,
				MechTilt:   MechTilt,
				TotlTilt:   TotlTilt,
			}

			// Insert or update the cell values in the MySQL database.
			_, err = db.Exec(`
                INSERT INTO tbCell (City, Sector_ID, Sector_Name, EnodebID, Enodeb_Name, EARFCN, PCI, PSS, SSS, TAC, Vendor, Longitude, Latitude, Style, Azimuth, Height, ElectTilt, MechTilt, TOTLETILT)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
        `, cellValues.City, cellValues.SectorID, cellValues.SectorName, cellValues.EnodebID, cellValues.EnodebName, cellValues.EARFCN, cellValues.PCI, cellValues.PSS, cellValues.SSS, cellValues.TAC, cellValues.Vendor, cellValues.Longitude, cellValues.Latitude, cellValues.Style, cellValues.Azimuth, cellValues.Height, cellValues.ElectTilt, cellValues.MechTilt, cellValues.TotlTilt)
			if err != nil {
				log.Println("Error inserting data:", err)
				errorList = append(errorList, err.Error())
			}
		}
	}
	if len(errorList) > 0 {
		return errors.New(strings.Join(errorList, "\n"))
	} else {
		return nil
	}
}
