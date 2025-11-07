package inmate

import "time"

type InmateConfirm struct {
	InmateId string    `dynamodbav:"partition_key" json:"id"`
	ID       string    `dynamodbav:"sort_key"`
	Creation time.Time `dynamodbav:"creation_date_and_time" json:"creation_date_and_time"`
}

type InmateAttempt struct {
	InmateId string    `dynamodbav:"partition_key" json:"id"`
	ID       string    `dynamodbav:"sort_key"`
	Creation time.Time `dynamodbav:"creation_date_and_time" json:"creation_date_and_time"`
	Reason   string    `dynamodbav:"reason" json:"reason"`
	Attempts int8      `dynamodbav:"attempts" json:"attempts"`
}
