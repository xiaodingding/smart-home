package responses

import (
	"github.com/e154/smart-home/api/server/v1/models"
)

// swagger:response TemplateItemList
type TemplateItemList struct {
	// in:body
	Body struct {
		Items []*models.TemplateItem `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}

// swagger:response TemplateItemSortedList
type TemplateItemSortedList struct {
	// in:body
	Body struct {
		Items []string `json:"items"`
		Total int      `json:"total"`
	}
}