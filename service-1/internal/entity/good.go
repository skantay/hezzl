package entity

import (
	"encoding/json"
	"time"
)

type Good struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"created_at"`
}

func (g Good) MarshalBinary() ([]byte, error) {
	return json.Marshal(g)
}
