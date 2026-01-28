package rsvp

import "encoding/csv"

// Csv is used to provide rsvp with a way to render your type as text/csv.
type Csv interface {
	// MarshalCsv will be called if the Accept header contains text/csv and it is matched, or the URL path ends with .csv
	MarshalCsv(w *csv.Writer) error
}
