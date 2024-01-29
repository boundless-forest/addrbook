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
