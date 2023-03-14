package main

import (
	"os"
	"testing"
)

func TestWAL_WriteAndRead(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "wal")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	wal, err := New(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer wal.Close()

	err = wal.Write("key1", "value1")
	if err != nil {
		t.Fatalf("Failed to write: %v", err)
	}

	value, err := wal.Read("key1")
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}

	if value != "value1" {
		t.Errorf("Expected value1, got %s", value)
	}
}

func TestWAL_Recover(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "wal")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	wal, err := New(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	err = wal.Write("key1", "value1")
	if err != nil {
		t.Fatalf("Failed to write: %v", err)
	}
	wal.Close()

	walRecovered, err := New(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer walRecovered.Close()

	err = walRecovered.Recover()
	if err != nil {
		t.Fatalf("Failed to recover: %v", err)
	}

	valueRecovered, err := walRecovered.Read("key1")
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}

	if valueRecovered != "value1" {
		t.Errorf("Expected value1, got %s", valueRecovered)
	}
}
