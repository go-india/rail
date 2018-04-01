package rail

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type (
	// WindowHour defines window time in hours to search.
	WindowHour uint8
)

// Window Hours
const (
	WindowHour2 = 1 + iota
	WindowHour4
)

// date return compatible date value
func date(t time.Time) string { return t.Format("02-01-2006") }

// Response is the stand response object that comes with every response
type Response struct {
	Debit        int `json:"debit"`
	ResponseCode int `json:"response_code"`
}

// Available holds an available item
type Available struct {
	Status string    `json:"status"`
	Date   time.Time `json:"date"`
}

// UnmarshalJSON convert JSON data to struct
func (a *Available) UnmarshalJSON(data []byte) error {
	type Alias Available
	t := struct {
		Alias
		Date string `json:"date"`
	}{}
	if err := json.Unmarshal(data, &t); err != nil {
		return errors.Wrap(err, "UnmarshalJSON failed")
	}

	*a = Available(t.Alias)

	if t.Date != "" {
		date, err := time.Parse("2-1-2006", t.Date)
		if err != nil {
			return errors.Wrap(err, "parse Date failed")
		}
		a.Date = date
	}
	return nil
}

// Day holds day details
type Day struct {
	Runs bool   // `json:"runs,omitempty"`
	Code string `json:"code,omitempty"`
}

// UnmarshalJSON convert JSON data to struct
func (d *Day) UnmarshalJSON(data []byte) error {
	type Alias Day
	t := struct {
		Alias
		Run string `json:"runs"`
	}{}
	if err := json.Unmarshal(data, &t); err != nil {
		return errors.Wrap(err, "UnmarshalJSON failed")
	}

	*d = Day(t.Alias)
	d.Runs = t.Run == "Y"
	return nil
}

// Quota holds quota details
type Quota struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Class holds class details
type Class struct {
	Available *bool  // `json:"available,omitempty"`
	Name      string `json:"name,omitempty"`
	Code      string `json:"code,omitempty"`
}

// UnmarshalJSON convert JSON data to struct
func (c *Class) UnmarshalJSON(data []byte) error {
	type Alias Class
	t := struct {
		Alias
		Avail string `json:"available"`
	}{}
	if err := json.Unmarshal(data, &t); err != nil {
		return errors.Wrap(err, "UnmarshalJSON failed")
	}

	*c = Class(t.Alias)

	b := t.Avail == "Y"
	c.Available = &b
	return nil
}

// Train holds train details
type Train struct {
	Name    string  `json:"name"`
	Number  uint32  `json:"number,string"`
	Classes []Class `json:"classes,omitempty"`
	Days    []Day   `json:"days,omitempty"`
}

// Passenger holds passenger details
type Passenger struct {
	Number        *uint16 `json:"no,omitempty"`
	CurrentStatus *string `json:"current_status,omitempty"`
	BookingStatus *string `json:"booking_status,omitempty"`
}

// Station holds station details
type Station struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Name      string  `json:"name"`
	Code      string  `json:"code"`
}

// Route holds route details
type Route struct {
	ActualArrivalDate    *time.Time //`json:"actarr_date,omitempty"`
	ScheduledArrivalDate *time.Time //`json:"scharr_date,omitempty"`

	ScheduledArrival   *time.Time //`json:"scharr,omitempty"`
	ScheduledDeparture *time.Time //`json:"schdep,omitempty"`
	ActualDeparture    *time.Time //`json:"actdep,omitempty"`
	ActualArrival      *time.Time //`json:"actarr,omitempty"`

	HasArrived  *bool `json:"has_arrived,omitempty"`
	HasDeparted *bool `json:"has_departed,omitempty"`

	Station       *Station `json:"station,omitempty"`
	Status        *string  `json:"status,omitempty"`
	LateByMinutes *int     `json:"latemin,omitempty"`
	Distance      *float64 `json:"distance,omitempty"`
	Day           *int     `json:"day,omitempty"`
	Number        *int     `json:"no,omitempty"`
	Halt          *int     `json:"halt,omitempty"`
}

// UnmarshalJSON convert JSON data to struct
func (r *Route) UnmarshalJSON(data []byte) error {
	type Alias Route
	t := struct {
		Alias
		ActualArrivalDate    string `json:"actarr_date"`
		ScheduledArrivalDate string `json:"scharr_date"`

		ScheduledArrival   string `json:"scharr"`
		ScheduledDeparture string `json:"schdep"`
		ActualDeparture    string `json:"actdep"`
		ActualArrival      string `json:"actarr"`
	}{}
	if err := json.Unmarshal(data, &t); err != nil {
		return errors.Wrap(err, "UnmarshalJSON failed")
	}

	*r = Route(t.Alias)

	if len(t.ScheduledArrival) == 5 {
		sa, err := time.Parse("15:04", t.ScheduledArrival)
		if err != nil {
			return errors.Wrap(err, "parse ScheduledArrival failed")
		}
		r.ScheduledArrival = &sa
	}

	if len(t.ScheduledDeparture) == 5 {
		sd, err := time.Parse("15:04", t.ScheduledDeparture)
		if err != nil {
			return errors.Wrap(err, "parse ScheduledDeparture failed")
		}
		r.ScheduledDeparture = &sd
	}

	if len(t.ActualDeparture) == 5 {
		ad, err := time.Parse("15:04", t.ActualDeparture)
		if err != nil {
			return errors.Wrap(err, "parse ActualDeparture failed")
		}
		r.ActualDeparture = &ad
	}

	if len(t.ActualArrival) == 5 {
		aa, err := time.Parse("15:04", t.ActualArrival)
		if err != nil {
			return errors.Wrap(err, "parse ActualArrival failed")
		}
		r.ActualArrival = &aa
	}

	if t.ActualArrivalDate != "" {
		aad, err := time.Parse("2 Jan 2006", t.ActualArrivalDate)
		if err != nil {
			return errors.Wrap(err, "parse ActualArrivalDate failed")
		}
		r.ActualArrivalDate = &aad
	}

	if t.ScheduledArrivalDate != "" {
		sad, err := time.Parse("2 Jan 2006", t.ScheduledArrivalDate)
		if err != nil {
			return errors.Wrap(err, "parse ScheduledArrivalDate failed")
		}
		r.ScheduledArrivalDate = &sad
	}

	return nil
}
