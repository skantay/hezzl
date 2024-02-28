package model

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

// func (g *Good) UnmarshalBinary(data []byte) error {
// 	b := bytes.NewBuffer(data)
// 	_, err := fmt.Fscanln(
// 		b,
// 		g.ID,
// 		g.ProjectID,
// 		g.Name,
// 		g.Description,
// 		g.Priority,
// 		g.Removed,
// 		g.CreatedAt)

// 	return err
// }
