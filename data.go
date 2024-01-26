package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type ContractItem struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Note    string `json:"note"`
}

type WorkSpace = []ContractItem
type WorkSpaces = map[string]WorkSpace

func LoadWorkSpaces() (WorkSpaces, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dataPath := filepath.Join(homeDir, ".addrbook", "data.json")
	_, err = os.Stat(dataPath)
	if os.IsNotExist(err) {
		return WorkSpaces{}, nil
	}

	file, err := os.Open(dataPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)
	var workspaces WorkSpaces
	err = json.Unmarshal(byteValue, &workspaces)
	if err != nil {
		return nil, err
	}

	return workspaces, nil
}

func SaveWorkSpaces(data *WorkSpaces) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dataPath := filepath.Join(homeDir, ".addrbook", "data.json")
	_, err = os.Stat(dataPath)
	if os.IsNotExist(err) {
		panic(err)
	}

	file, err := os.Open(dataPath)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	_, err = file.Write(byteValue)
	return err

}
