package inmate

import "time"

type InmateConfirm struct {
	Creation time.Time `dynamodbav:"creation_date_and_time" json:"creation_date_and_time"`
}

type InmateAttempt struct {
	Creation time.Time `dynamodbav:"creation_date_and_time" json:"creation_date_and_time"`
}
