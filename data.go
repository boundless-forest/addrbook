package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type Contract struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Note    string `json:"note"`
}

type WorkSpace struct {
	Contract map[string]Contract
}

type DataBase struct {
	Workspaces map[string]WorkSpace
}

var ErrWorkSpaceExists = errors.New("workspace already exists")
var ErrWorkSpaceNotFound = errors.New("workspace not found")

var ErrContractExists = errors.New("contract already exists")
var ErrContractNotFound = errors.New("contract not found")

func (db *DataBase) CreateWorkSpace(name string) error {
	if db.Workspaces == nil {
		db.Workspaces = make(map[string]WorkSpace)
	}

	if _, ok := db.Workspaces[name]; ok {
		return ErrWorkSpaceExists
	}

	db.Workspaces[name] = WorkSpace{Contract: make(map[string]Contract)}
	return nil
}

func (db *DataBase) DeleteWorkSpace(name string) error {
	if db.Workspaces == nil {
		return ErrWorkSpaceNotFound
	}

	delete(db.Workspaces, name)
	return nil
}

func (db *DataBase) ListWorkSpaces() []string {
	spaces := make([]string, 0, len(db.Workspaces))
	for name := range db.Workspaces {
		spaces = append(spaces, name)
	}
	return spaces
}

func (db *DataBase) Save(workspace string, contract string, address string, note string) error {
	if _, ok := db.Workspaces[workspace]; !ok {
		return ErrWorkSpaceNotFound
	}

	ws := db.Workspaces[workspace]
	if _, ok := ws.Contract[contract]; ok {
		return ErrContractExists
	}

	ws.Contract[contract] = Contract{
		Name:    contract,
		Address: address,
		Note:    note,
	}

	db.Workspaces[workspace] = ws
	return nil
}

func (db *DataBase) Update(workspace string, contract string, address string, note string) error {
	ws, ok := db.Workspaces[workspace]
	if !ok {
		return ErrWorkSpaceNotFound
	}

	existingContract, ok := ws.Contract[contract]
	if !ok {
		return ErrContractNotFound
	}

	if address == existingContract.Address && note == existingContract.Note {
		return errors.New("the new contract information is the same as the old one")
	}

	ws.Contract[contract] = Contract{
		Name:    contract,
		Address: address,
		Note:    note,
	}

	db.Workspaces[workspace] = ws
	return nil
}

func (db *DataBase) Delete(workspace string, contract string) error {
	ws, ok := db.Workspaces[workspace]
	if !ok {
		return ErrWorkSpaceNotFound
	}

	if _, ok := ws.Contract[contract]; !ok {
		return ErrContractNotFound
	}

	delete(ws.Contract, contract)

	db.Workspaces[workspace] = ws
	return nil
}

func LoadDB() (DataBase, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return DataBase{}, err
	}

	dataPath := filepath.Join(homeDir, ".addrbook", "data.json")
	dirPath := filepath.Dir(dataPath)

	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return DataBase{}, err
		}
	}

	_, err = os.Stat(dataPath)
	if os.IsNotExist(err) {
		defaultDB := DataBase{}
		file, err := os.Create(dataPath)
		if err != nil {
			return defaultDB, err
		}
		defer file.Close()

		byteValue, err := json.MarshalIndent(defaultDB, "", "    ")
		if err != nil {
			return defaultDB, err
		}

		_, err = file.Write(byteValue)
		if err != nil {
			return defaultDB, err
		}

		return defaultDB, nil
	} else {
		file, err := os.Open(dataPath)
		if err != nil {
			return DataBase{}, err
		}
		defer file.Close()

		byteValue, _ := io.ReadAll(file)
		db := DataBase{}
		err = json.Unmarshal(byteValue, &db)
		if err != nil {
			return DataBase{}, err
		}
		return db, nil

	}
}

func SaveToDB(db *DataBase) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dataPath := filepath.Join(homeDir, ".addrbook", "data.json")
	dirPath := filepath.Dir(dataPath)

	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	_, err = os.Stat(dataPath)
	if os.IsNotExist(err) {
		panic(err)
	}

	file, err := os.Create(dataPath)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(db, "", "    ")
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return err
}
