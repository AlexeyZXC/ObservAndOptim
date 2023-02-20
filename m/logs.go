package m

import "time"

type Logs struct {
	Id        string    `json:"id,omitempty"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
