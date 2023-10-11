package db

import (
	"encoding/json"
	"os"
	"time"
)

type Item struct {
	Id     int    "json:\"id\""
	Name   string "json:\"name\""
	Status bool   "json:\"status\""
	Date   string "json:\"date\""
}

func GetJson(filePath string) ([]Item, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []Item

	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func GetDatabase(filePath string) ([]Item, error) {
	dbData, err := GetJson(filePath)
	if err != nil {
		return nil, err
	}
	return dbData, nil
}

func SaveJson(data []Item, filePath string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func NewEntry(name string, status bool, date string, filePath string) error {
	data, err := GetJson(filePath)
	if err != nil {
		return err
	}

	newItem := Item{
		Id:     len(data) + 1,
		Name:   name,
		Status: status,
		Date:   time.Now().Format("02/01/2006"),
	}

	data = append(data, newItem)
    
	if err := SaveJson(data, filePath); err != nil {
		return err
	}

	return nil
}

func DeleteEntry(id int, filePath string) error {
	data, err := GetJson(filePath)
	if err != nil {
		return err
	}

	var indextoremove = -1
	for i, item := range data {
		if item.Id == id {
			indextoremove = i
			break
		}
	}

	if indextoremove == -1 {
		return nil
	}

	data = append(data[:indextoremove], data[indextoremove+1:]...)

	if err := SaveJson(data, filePath); err != nil {
		return err
	}

	return nil
}
