package models

import "time"

type WorkflowScenario struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	SystemName string    `json:"system_name"`
	WorkflowId int64     `json:"workflow_id"`
	Scripts    []*Script `json:"scripts"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}