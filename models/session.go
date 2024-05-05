package models

import "time"

type Session struct {
	ID      string    `json:"id"`
	Expires time.Time `json:"expires"`
}

func (s Session) Expired() bool {
	return s.Expires.Before(time.Now())
}
