package main

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateWorkSpace(t *testing.T) {
	db := DataBase{}
	db.CreateWorkSpace("test")

	if len(db.Workspaces) != 1 {
		t.Errorf("Expected 1 workspace, got %d", len(db.Workspaces))
	}
	if _, ok := db.Workspaces["test"]; !ok {
		t.Errorf("Expected workspace 'test' to exist")
	}
	assert.NotNil(t, db.Workspaces["test"], "workspace 'test' should exist")
	if err := db.CreateWorkSpace("test"); !errors.Is(err, ErrWorkSpaceExists) {
		t.Errorf("Expected error 'workspace already exists', got '%v'", err)
	}

	// Test multiple workspaces
	db.CreateWorkSpace("test1")
	db.CreateWorkSpace("test2")
	db.CreateWorkSpace("test3")
	if len(db.Workspaces) != 4 {
		t.Errorf("Expected 4 workspace, got %d", len(db.Workspaces))
	}
}

func TestDeleteWorkSpace(t *testing.T) {
	db := DataBase{}
	db.CreateWorkSpace("test0")
	db.CreateWorkSpace("test1")
	db.CreateWorkSpace("test2")
	if len(db.Workspaces) != 3 {
		t.Errorf("Expected 3 workspace, got %d", len(db.Workspaces))
	}

	if err := db.DeleteWorkSpace("test5"); errors.Is(err, ErrWorkSpaceNotFound) {
		t.Errorf("Expected error 'workspace not found', got '%v'", err)
	}

	db.DeleteWorkSpace("test1")
	if len(db.Workspaces) != 2 {
		t.Errorf("Expected 2 workspace, got %d", len(db.Workspaces))
	}
	if _, ok := db.Workspaces["test"]; ok {
		t.Errorf("workspace 'test1' should be deleted")
	}
}

func TestListWorkSpaces(t *testing.T) {
	db := DataBase{}
	db.CreateWorkSpace("test0")
	db.CreateWorkSpace("test1")
	db.CreateWorkSpace("test2")

	if len(db.ListWorkSpaces()) != 3 {
		t.Errorf("Expected 3 workspace, got %d", len(db.ListWorkSpaces()))
	}
	if reflect.DeepEqual(db.ListWorkSpaces(), []string{"test2", "test1", "test0"}) {
		t.Errorf("Expected workspaces ['test2', 'test1', 'test0'], got %v", db.ListWorkSpaces())
	}
}

func TestSave(t *testing.T) {
	db := DataBase{}
	if err := db.Save("test", "test", "test", "test"); !errors.Is(err, ErrWorkSpaceNotFound) {
		t.Errorf("Expected error 'workspace not found', got '%v'", err)
	}

	db.CreateWorkSpace("test0")
	db.Save("test0", "test contract", "test address", "test note")
	db.Save("test0", "test contract 1", "test address 1", "test note")

	if len(db.Workspaces["test0"].Contract) != 2 {
		t.Errorf("Expected 2 contract, got %d", len(db.Workspaces["test0"].Contract))
	}

	contract := db.Workspaces["test0"].Contract["test contract"]
	assert.Equal(t, "test contract", contract.Name)
	assert.Equal(t, "test address", contract.Address)
	assert.Equal(t, "test note", contract.Note)
}

func TestUpdate(t *testing.T) {
	db := DataBase{}
	if err := db.Update("test", "test", "test", "test"); !errors.Is(err, ErrWorkSpaceNotFound) {
		t.Errorf("Expected error 'workspace not found', got '%v'", err)
	}

	db.CreateWorkSpace("test0")
	db.Save("test0", "test contract", "test address", "test note")
	contract := db.Workspaces["test0"].Contract["test contract"]
	assert.Equal(t, "test address", contract.Address)

	if err := db.Update("test0", "test contract 2", "test address updated", "test note"); !errors.Is(err, ErrContractNotFound) {
		t.Errorf("Expect error 'contract not found', got '%v'", err)
	}

	if err := db.Update("test0", "test contract", "test address updated", "test note"); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	contract = db.Workspaces["test0"].Contract["test contract"]
	assert.Equal(t, "test address updated", contract.Address)

}

func TestDelete(t *testing.T) {
	db := DataBase{}
	if err := db.Delete("test", "test"); !errors.Is(err, ErrWorkSpaceNotFound) {
		t.Errorf("Expected error 'workspace not found', got '%v'", err)
	}

	db.CreateWorkSpace("test0")
	db.Save("test0", "test contract 1", "test address 1", "test note")
	db.Save("test0", "test contract 2", "test address 2", "test note")
	db.Save("test0", "test contract 3", "test address 3", "test note")
	db.Save("test0", "test contract 4", "test address 4", "test note")
	workspace := db.Workspaces["test0"]
	if len(workspace.Contract) != 4 {
		t.Errorf("Expected 4 contract, got %d", len(workspace.Contract))
	}

	if err := db.Delete("test0", "test contract 9"); !errors.Is(err, ErrContractNotFound) {
		t.Errorf("Expect error 'contract not found', got '%v'", err)
	}
	if err := db.Delete("test0", "test contract 2"); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(workspace.Contract) != 3 {
		t.Errorf("Expected 3 contract, got %d", len(workspace.Contract))
	}
	if contract, ok := db.Workspaces["test0"].Contract["test contract 2"]; ok {
		t.Errorf("Expected contract 'test contract 2' to be deleted, got %v", contract)
	}
}
