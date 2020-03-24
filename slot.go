package main

import (
	"fmt"
	"time"
)

// StateSelectable is used for ah.nl delivery time slots that are available for
// reservation.
const StateSelectable = "selectable"

const dateLayout = "2006-01-02 15:04"

var loc = func() *time.Location {
	loc, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		panic(err)
	}
	return loc
}()

// ParseAvailableDeliverySlots returns an array of available delivery slots,
// parsed from an array of delivery dates returned from the ah.nl API.
func ParseAvailableDeliverySlots(d DeliveryDates) ([]DeliverySlot, error) {
	var slots []DeliverySlot
	for _, dd := range d {
		for _, timeSlot := range dd.DeliveryTimeSlots {
			if timeSlot.State != StateSelectable {
				continue
			}

			from, err := time.ParseInLocation(dateLayout, fmt.Sprintf("%v %v", dd.Date, timeSlot.From), loc)
			if err != nil {
				return nil, fmt.Errorf("cannot parse from date: %v", err)
			}

			to, err := time.ParseInLocation(dateLayout, fmt.Sprintf("%v %v", dd.Date, timeSlot.To), loc)
			if err != nil {
				return nil, fmt.Errorf("cannot parse to date: %v", err)
			}

			slots = append(slots, DeliverySlot{
				From:      from,
				To:        to,
				HRef:      timeSlot.NavItem.Link.HRef,
				Available: true,
			})
		}
	}

	return slots, nil
}
