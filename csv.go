package rsvp

import "encoding/csv"

type Csv interface {
	MarshalCsv(w *csv.Writer) error
}
