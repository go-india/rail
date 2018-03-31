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
	Runs bool   // `json:"runs"`
	Code string `json:"code"`
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
	Available bool   // `json:"available"`
	Name      string `json:"name,omitempty"`
	Code      string `json:"code"`
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
	c.Available = t.Avail == "Y"
	return nil
}

// Train holds train details
type Train struct {
	Name    string  `json:"name"`
	Number  int     `json:"number,string"`
	Classes []Class `json:"classes"`
	Days    []Day   `json:"days"`
}

// Passenger holds passenger details
type Passenger struct {
	Number        int    `json:"no"`
	CurrentStatus string `json:"current_status"`
	BookingStatus string `json:"booking_status"`
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
	ActualArrivalDate    time.Time //`json:"actarr_date"`
	ScheduledArrivalDate time.Time //`json:"scharr_date"`

	ScheduledArrival   time.Time //`json:"scharr"`
	ScheduledDeparture time.Time //`json:"schdep"`
	ActualDeparture    time.Time //`json:"actdep"`
	ActualArrival      time.Time //`json:"actarr"`

	HasArrived  bool `json:"has_arrived"`
	HasDeparted bool `json:"has_departed"`

	Station       Station `json:"station"`
	Status        string  `json:"status"`
	LateByMinutes int     `json:"latemin"`
	Distance      float64 `json:"distance"`
	Day           int     `json:"day"`
	Number        int     `json:"no,omitempty"`
	Halt          int     `json:"halt,omitempty"`
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
		r.ScheduledArrival = sa
	}

	if len(t.ScheduledDeparture) == 5 {
		sd, err := time.Parse("15:04", t.ScheduledDeparture)
		if err != nil {
			return errors.Wrap(err, "parse ScheduledDeparture failed")
		}
		r.ScheduledDeparture = sd
	}

	if len(t.ActualDeparture) == 5 {
		ad, err := time.Parse("15:04", t.ActualDeparture)
		if err != nil {
			return errors.Wrap(err, "parse ActualDeparture failed")
		}
		r.ActualDeparture = ad
	}

	if len(t.ActualArrival) == 5 {
		aa, err := time.Parse("15:04", t.ActualArrival)
		if err != nil {
			return errors.Wrap(err, "parse ActualArrival failed")
		}
		r.ActualArrival = aa
	}

	if t.ActualArrivalDate != "" {
		aad, err := time.Parse("2 Jan 2006", t.ActualArrivalDate)
		if err != nil {
			return errors.Wrap(err, "parse ActualArrivalDate failed")
		}
		r.ActualArrivalDate = aad
	}

	if t.ScheduledArrivalDate != "" {
		sad, err := time.Parse("2 Jan 2006", t.ScheduledArrivalDate)
		if err != nil {
			return errors.Wrap(err, "parse ScheduledArrivalDate failed")
		}
		r.ScheduledArrivalDate = sad
	}

	return nil
}
