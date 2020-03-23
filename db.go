package main

import (
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

// Database is used to keep a (persistent) state of notifcations.
type Database struct {
	db *bbolt.DB
}

// DeliverySlot represents an ah.nl delivery slot in the database.
type DeliverySlot struct {
	From      time.Time
	To        time.Time
	Available bool
}

// NewDatabase returns a new Database.
func NewDatabase(path string) (*Database, error) {
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot open database: %v", err)
	}

	return &Database{db}, nil
}

// Close closes the underlying database.
func (db *Database) Close() {
	db.db.Close()
}
