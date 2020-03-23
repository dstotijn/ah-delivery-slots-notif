package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml"
)

// Config r
type Config struct {
	Notifs NotifConfigs
}

// NotifConfigs represents a notification config for a postal code, a map of
// postal codes ([\d{4}[A-Z]{2}) to phone numbers (E.164).
type NotifConfigs map[string][]string

var (
	postalCodeRegExp  = regexp.MustCompile(`^\d{4}[A-Z]{2}`)
	phoneNumberRegExp = regexp.MustCompile(`^\+?\d{11}`)
)

// LoadConfig loads a config from a TOML file.
func LoadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("cannot open config file: %v", err)
	}

	var cfg Config
	err = toml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("cannot parse config: %v", err)
	}

	err = cfg.Validate()
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Validate checks the validity of a config.
func (cfg Config) Validate() error {
	for postalCode, phoneNumbers := range cfg.Notifs {
		if !postalCodeRegExp.Match([]byte(postalCode)) {
			return fmt.Errorf("invalid postal code (%v)", postalCode)
		}

		for i, phoneNumber := range phoneNumbers {
			if !phoneNumberRegExp.Match([]byte(phoneNumber)) {
				return fmt.Errorf("invalid phone number (%v)", phoneNumber)
			}
			phoneNumbers[i] = strings.TrimPrefix(phoneNumber, "+")
		}
	}

	return nil
}
