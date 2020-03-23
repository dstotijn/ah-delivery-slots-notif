package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

// Service is a container for the HTTP client and database.
type Service struct {
	ahClient *AHClient
	db       *Database
	notifs   NotifConfigs
}

func main() {
	dbPath := flag.String("db", "ah-delivery-slots-notif.db", "Database path")
	cfgPath := flag.String("config", "config.toml", "Config path")
	flag.Parse()

	cfg, err := LoadConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	ahClient := NewAHClient()
	db, err := NewDatabase(*dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	svc := Service{
		ahClient: ahClient,
		db:       db,
		notifs:   cfg.Notifs,
	}

	if err := run(svc); err != nil {
		panic(err)
	}
}

func run(svc Service) error {
	// Immediately run first batch.
	if err := runBatch(svc, time.Now()); err != nil {
		log.Printf("[ERROR]: Could not run batch: %v", err)
	}

	// Run subsequent batches every 15 seconds.
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			if err := runBatch(svc, t); err != nil {
				log.Printf("[ERROR]: Could not run batch: %v", err)
			}
		}
	}
}

func runBatch(svc Service, ts time.Time) error {
	for postalCode := range svc.notifs {
		deliveryDates, err := svc.ahClient.GetDeliveryDates(postalCode)
		if err != nil {
			return fmt.Errorf("cannot get delivery dates for postal code (%v): %v", postalCode, err)
		}
		fmt.Printf("%+v", deliveryDates)

		// TODO: Check/store database, send notifications.
	}

	return nil
}
