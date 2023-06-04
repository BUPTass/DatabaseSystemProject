package Query

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

type communityMsg struct {
	Output string `json:"output"`
	Pic    string `json:"pic"`
}

func GetCommunity(db *sql.DB) ([]byte, error) {
	os.Mkdir("/tmp/tmp1", 0777)
	os.Chmod("/tmp/tmp1", 0777)
	if err := getCoordinate(db); err != nil {
		log.Println(err)
		return nil, err
	}
	if err := getC2I(db); err != nil {
		log.Println(err)
		return nil, err
	}
	louvainOut, err := runLouvain()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	randomName := fmt.Sprintf("louvain-%d.png", time.Now().UnixNano())
	if err := os.Rename("/tmp/tmp1/louvain.png", "./download/"+randomName); err != nil {
		log.Println(err)
		return nil, err
	}
	jsonData, _ := json.Marshal(communityMsg{louvainOut, "/download/" + randomName})
	return jsonData, nil

}

func runLouvain() (string, error) {
	cmd := exec.Command("python3", "./run.py")

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return stdout.String(), nil
}

func getCoordinate(db *sql.DB) error {
	// Execute the query to fetch data from tbCell
	rows, err := db.Query("SELECT SECTOR_ID, LONGITUDE, LATITUDE FROM tbCell")
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close()

	coordinateString := "{"

	// Iterate over the rows and populate the coordinateDict
	var sectorID string
	var longitude, latitude float32
	for rows.Next() {
		err := rows.Scan(&sectorID, &longitude, &latitude)
		if err != nil {
			log.Println(err)
			continue
		}
		coordinateString += fmt.Sprintf("'%s': [%f, %f], ", sectorID, longitude, latitude)
	}
	coordinateString = coordinateString[:len(coordinateString)-2] + "}"

	// Write the coordinateString to a file
	savePath := "/tmp/tmp1/coordinate.txt"
	err = ioutil.WriteFile(savePath, []byte(coordinateString), os.ModePerm)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func getC2I(db *sql.DB) error {
	// Remove tbC2I.txt first
	os.Remove("/tmp/tmp1/tbC2I.txt")

	exportSQL := `select SCELL, NCELL,C2I_Mean from tbC2I
			INTO OUTFILE '/tmp/tmp1/tbC2I.txt'
    		FIELDS TERMINATED BY ' '
    		OPTIONALLY ENCLOSED BY ''
    		LINES TERMINATED BY '\n';`

	_, err := db.Exec(exportSQL)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
