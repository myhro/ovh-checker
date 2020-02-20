package hardware

import (
	"time"
)

// LatestOffers DB structure for latest hardware offers
type LatestOffers struct {
	ID        int        `json:"id"`
	Server    string     `json:"server"`
	Country   string     `json:"country"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}
