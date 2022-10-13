package schoolclasses

import "time"

type SchoolClass struct {
	ID             string    `json:"id"`
	Grade          int       `json:"grade"`
	Letter         string    `json:"letter"`
	NumberOfPupils int       `json:"numberOfPupils"`
	DateFrom       time.Time `json:"dateFrom"`
	DateTo         time.Time `json:"dateTo"`
}
