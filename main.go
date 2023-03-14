package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type WAL struct {
	file     *os.File
	writer   *bufio.Writer
	filePath string
	data     map[string]string
	mu       sync.RWMutex
}

type logEntry struct {
	Key   string
	Value string
}

func New(filePath string) (*WAL, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &WAL{
		file:     file,
		writer:   bufio.NewWriter(file),
		filePath: filePath,
		data:     make(map[string]string),
	}, nil
}

func (w *WAL) Write(key, value string) error {
	entry := logEntry{Key: key, Value: value}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = w.writer.WriteString(fmt.Sprintf("%s\n", data))
	if err != nil {
		return err
	}

	err = w.writer.Flush()
	if err != nil {
		return err
	}

	w.mu.Lock()
	w.data[key] = value
	w.mu.Unlock()

	return nil
}

func (w *WAL) Read(key string) (string, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	value, ok := w.data[key]
	if !ok {
		return "", errors.New("key not found")
	}

	return value, nil
}

func (w *WAL) Close() error {
	return w.file.Close()
}

func (w *WAL) Recover() error {
	file, err := os.Open(w.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var entry logEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return err
		}

		w.mu.Lock()
		w.data[entry.Key] = entry.Value
		w.mu.Unlock()
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
