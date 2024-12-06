package models

import "time"

type UatResult struct {
	WorkItemId    int       `json:"workItemId"`
	ExecutionDate time.Time `json:"executionDate"`
	IsSuccess     bool      `json:"isSuccess"`
}
