package csvparser

import (
	"encoding/csv"
	"os"
)

func ReWriteCSV(path string, data []string) error {
	metaFile, err := os.OpenFile(path, os.O_RDWR, 0o666)
	if err != nil {
		return err
	}
	defer metaFile.Close()

	reader := csv.NewReader(metaFile)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for i := 0; i < len(records); i++ {
		if len(records[i]) < 4 {
			continue
		}
		if records[i][0] == data[0] && records[i][3] == "Active" {
			records[i][2] = data[2]
			records[i][3] = data[3]
		}
	}
	metaFile.Seek(0, 0)
	metaFile.Truncate(0)

	writer := csv.NewWriter(metaFile)
	defer writer.Flush()

	err = writer.WriteAll(records)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRecord(path string, name string) error {
	metaFile, err := os.OpenFile(path, os.O_RDWR, 0o666)
	if err != nil {
		return err
	}
	defer metaFile.Close()
	reader := csv.NewReader(metaFile)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for i := 0; i < len(records); i++ {
		if len(records[i]) < 4 {
			continue
		}
		if records[i][0] == name {
			records[i] = []string{}
		}
	}

	metaFile.Seek(0, 0)
	metaFile.Truncate(0)

	writer := csv.NewWriter(metaFile)
	defer writer.Flush()
	err = writer.WriteAll(records)
	if err != nil {
		return err
	}
	return nil
}

func WriteCSV(metaFile *os.File, data []string) error {
	writer := csv.NewWriter(metaFile)
	defer writer.Flush()

	err := writer.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func CheckEmpty(pathToCSV string) (bool, error) {
	metaFile, err := os.OpenFile(pathToCSV, os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		return false, err
	}
	defer metaFile.Close()

	reader := csv.NewReader(metaFile)
	records, err := reader.ReadAll()
	if err != nil {
		return false, err
	}

	if len(records) <= 1 {
		return true, nil
	}
	return false, nil
}
