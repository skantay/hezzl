package entity

import (
	"encoding/json"
	"time"
)

type Good struct {
	ID           int       `json:"id"`
	ProjectID    int       `json:"project_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Priority     int       `json:"priority"`
	Removed      bool      `json:"removed"`
	CreatedAt    time.Time `json:"-"`
	CreatedAtStr string    `json:"created_at"`
}

type Collection struct {
	Goods []Good `json:"goods"`
}

func (c *Collection) UnmarshalJSON(data []byte) error {
	type Alias Collection
	aux := &struct {
		Goods []Good `json:"goods"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	for i := range aux.Goods {
		aux.Goods[i].CreatedAt, _ = time.Parse(time.RFC3339Nano, aux.Goods[i].CreatedAtStr)
		c.Goods = append(c.Goods, aux.Goods[i])
	}
	return nil
}
