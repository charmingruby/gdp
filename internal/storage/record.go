package storage

import (
	"encoding/json"
	"os"
)

type Record struct {
	Data interface{} `json:"data"`
}

func SaveRecord(filePath string, newRecord Record) error {
	var records []Record

	if _, err := os.Stat(filePath); err == nil {
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(fileContent, &records); err != nil {
			return err
		}
	}

	records = append(records, newRecord)

	updatedContent, err := json.Marshal(records)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, updatedContent, 0644); err != nil {
		return err
	}

	return nil
}
