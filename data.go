package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
	cs map[string]Contract
}

type DataBase struct {
	workspaces map[string]WorkSpace
}

func (db *DataBase) CreateWorkSpace(name string) error {
	if db.workspaces == nil {
		db.workspaces = make(map[string]WorkSpace)
	}

	if _, ok := db.workspaces[name]; ok {
		return errors.New("workspace already exists")
	}

	db.workspaces[name] = WorkSpace{cs: make(map[string]Contract)}
	return nil
}

func (db *DataBase) DeleteWorkSpace(name string) error {
	if db.workspaces == nil {
		return errors.New("workspaces not found")
	}

	if _, ok := db.workspaces[name]; ok {
		delete(db.workspaces, name)
	}
	return errors.New("workspace not found")
}

func (db *DataBase) ListWorkSpaces() {
	for name, _ := range db.workspaces {
		fmt.Println(name)
	}
}

func (db *DataBase) Save(workspace string, contract string, address string, note string) error {
	if _, ok := db.workspaces[workspace]; ok {
		return errors.New("workspace not found")
	}

	ws := db.workspaces[workspace]
	if _, ok := ws.cs[contract]; ok {
		return errors.New("contract already exists")
	}

	ws.cs[contract] = Contract{
		Name:    contract,
		Address: address,
		Note:    note,
	}

	db.workspaces[workspace] = ws
	return nil
}

func (db *DataBase) Update(workspace string, contract string, address string, note string) error {
	if _, ok := db.workspaces[workspace]; ok {
		return errors.New("workspace not found")
	}

	ws := db.workspaces[workspace]
	if _, ok := ws.cs[contract]; !ok {
		return errors.New("contract not found")
	}

	ws.cs[contract] = Contract{
		Name:    contract,
		Address: address,
		Note:    note,
	}

	db.workspaces[workspace] = ws
	return nil
}

func (db *DataBase) Delete(workspace string, contract string) error {
	if _, ok := db.workspaces[workspace]; ok {
		return errors.New("workspace not found")
	}

	ws := db.workspaces[workspace]
	if _, ok := ws.cs[contract]; !ok {
		return errors.New("contract not found")
	}
	delete(ws.cs, contract)

	db.workspaces[workspace] = ws
	return nil
}

func LoadDB() (DataBase, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return DataBase{}, err
	}

	dataPath := filepath.Join(homeDir, ".addrbook", "data.json")
	_, err = os.Stat(dataPath)
	if os.IsNotExist(err) {
		return DataBase{}, nil
	}

	file, err := os.Open(dataPath)
	if err != nil {
		return DataBase{}, err
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)
	var db DataBase
	err = json.Unmarshal(byteValue, &db)
	if err != nil {
		return DataBase{}, err
	}

	return db, nil
}

func SaveToDB(db *DataBase) error {
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

	byteValue, err := json.MarshalIndent(db, "", "    ")
	if err != nil {
		return err
	}

	_, err = file.Write(byteValue)
	return err
}
