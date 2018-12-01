package server

import "gopkg.in/guregu/null.v3"

type RunResult struct {
	JobRunID string      `json:"jobRunId"`
	Status   string      `json:"status"`
	Error    null.String `json:"error"`
	Pending  bool        `json:"pending"`
	Data     Params      `json:"data"`
}

type Params struct {
	API []string `json:"api"`
	Paths []string `json:"paths"`
	AggregationType string `json:"aggregationType"`
	AggregateValue string `json:"aggregateValue"`
	FailedAPICount int `json:"failedApiCount"`
	APIErrors []string `json:"apiErrors"`
}
