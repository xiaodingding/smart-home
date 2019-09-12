package responses

import (
	"github.com/e154/smart-home/api/mobile/v1/models"
)

// swagger:response MapElementList
type MapElementList struct {
	// in:body
	Body struct {
		Items []*models.MapActiveElement `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}