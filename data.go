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
	Cs map[string]Contract
}

type DataBase struct {
	Workspaces map[string]WorkSpace
}

func (db *DataBase) CreateWorkSpace(name string) error {
	if db.Workspaces == nil {
		db.Workspaces = make(map[string]WorkSpace)
	}

	if _, ok := db.Workspaces[name]; ok {
		return errors.New("workspace already exists")
	}

	db.Workspaces[name] = WorkSpace{Cs: make(map[string]Contract)}
	return nil
}

func (db *DataBase) DeleteWorkSpace(name string) error {
	if db.Workspaces == nil {
		return errors.New("workspaces not found")
	}

	delete(db.Workspaces, name)
	return errors.New("workspace not found")
}

func (db *DataBase) ListWorkSpaces() {
	for name := range db.Workspaces {
		fmt.Println(name)
	}
}

func (db *DataBase) Save(workspace string, contract string, address string, note string) error {
	if _, ok := db.Workspaces[workspace]; !ok {
		return errors.New("workspace not found")
	}

	ws := db.Workspaces[workspace]
	if _, ok := ws.Cs[contract]; ok {
		return errors.New("contract already exists")
	}
	ws.Cs = map[string]Contract{
		contract: {
			Name:    contract,
			Address: address,
			Note:    note,
		},
	}

	db.Workspaces[workspace] = ws
	return nil
}

func (db *DataBase) Update(workspace string, contract string, address string, note string) error {
	if _, ok := db.Workspaces[workspace]; ok {
		return errors.New("workspace not found")
	}

	ws := db.Workspaces[workspace]
	if _, ok := ws.Cs[contract]; !ok {
		return errors.New("contract not found")
	}

	ws.Cs[contract] = Contract{
		Name:    contract,
		Address: address,
		Note:    note,
	}

	db.Workspaces[workspace] = ws
	return nil
}

func (db *DataBase) Delete(workspace string, contract string) error {
	if _, ok := db.Workspaces[workspace]; ok {
		return errors.New("workspace not found")
	}

	ws := db.Workspaces[workspace]
	if _, ok := ws.Cs[contract]; !ok {
		return errors.New("contract not found")
	}
	delete(ws.Cs, contract)

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

	fmt.Printf("In the Save function: %+v\n", db)
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
