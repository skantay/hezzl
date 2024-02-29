package schemas

import "github.com/skantay/hezzl/internal/entity"

// Requests

type CreateRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdatePriorityRequest struct {
	NewPriority int `json:"newPriority" validate:"required"`
}

type UpdateGoodRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description" validate:"required"`
}

// Responses

type ListResponse struct {
	Meta struct {
		Total   int `json:"total"`
		Removed int `json:"removed"`
		Limit   int `json:"limit"`
		Offset  int `json:"offset"`
	} `json:"meta"`
	Goods []struct {
		entity.Good
	} `json:"goods"`
}

type DeletedListResponse struct {
	Id         int  `json:"id"`
	CampaignID int  `json:"campignID"`
	Removed    bool `json:"removed"`
}
