package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var defaultBaseURL = "https://www.ah.nl/service/rest/delegate"

// DeliveryTimeSlot represents an ah.nl delivery time slot.
type DeliveryTimeSlot struct {
	From    string  `json:"from"`
	To      string  `json:"to"`
	Value   float64 `json:"value"`
	State   string  `json:"state"`
	NavItem struct {
		Link struct {
			HRef string `json:"href"`
		} `json:"link"`
	} `json:"navItem"`
	Eco bool `json:"eco"`
}

// DeliveryDate represents a delivery date (day) on ah.nl.
type DeliveryDate struct {
	Date              string             `json:"date"`
	DeliveryTimeSlots []DeliveryTimeSlot `json:"deliveryTimeSlots"`
}

// DeliveryDates represents an array of ah.nl delivery dates.
type DeliveryDates []DeliveryDate

// DeliveryDatesResponse represents an API response with delivery dates.
type deliveryDatesResponse struct {
	Embedded struct {
		Lanes []struct {
			Embedded struct {
				Items []struct {
					Type     string          `json:"type"`
					Embedded json.RawMessage `json:"_embedded"`
				} `json:"items"`
			} `json:"_embedded"`
		} `json:"lanes"`
	} `json:"_embedded"`
}

// AHClient defines an HTTP client for accessing the AH.nl REST API.
type AHClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewAHClient returns a new AHClient.
func NewAHClient() *AHClient {
	return &AHClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: defaultBaseURL,
	}
}

// GetDeliveryDates queries the ah.nl API to get delivery slots for a given
// postal code.
func (ah *AHClient) GetDeliveryDates(postalCode string) ([]DeliveryDate, error) {
	req, err := http.NewRequest("GET", ah.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %v", err)
	}

	q := req.URL.Query()
	q.Set("url", "kies-moment/bezorgen/"+postalCode)
	req.URL.RawQuery = q.Encode()

	resp, err := ah.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received unexpected status code: %v", resp.StatusCode)
	}

	var deliveryDates DeliveryDates
	if err = json.NewDecoder(resp.Body).Decode(&deliveryDates); err != nil {
		return nil, fmt.Errorf("cannot parse HTTP response body: %v", err)
	}

	return deliveryDates, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (dd *DeliveryDates) UnmarshalJSON(b []byte) error {
	var resp deliveryDatesResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return err
	}

	for _, lane := range resp.Embedded.Lanes {
		for _, item := range lane.Embedded.Items {
			if item.Type != "DeliveryDateSelector" {
				continue
			}

			var embed struct {
				DeliveryDates []DeliveryDate
			}
			if err := json.Unmarshal(item.Embedded, &embed); err != nil {
				return err
			}
			*dd = embed.DeliveryDates

			return nil
		}
	}

	return errors.New("items with type `DeliveryDateSelector` are missing")
}
