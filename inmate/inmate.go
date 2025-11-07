package inmate

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Gender uint8

const (
	Female Gender = 0
	Male   Gender = 1
)

const INMATE_ID_PREFIX = "i#"
const INMATE_CONFIRM_ID_PREFIX = "i#ic#"
const INMATE_ATTEMPT_ID_PREFIX = "i#ia#"

type Inmate struct {
	ID       string `dynamodbav:"partition_key" json:"id"`
	LastName string `dynamodbav:"sort_key" json:"last_name"`
	Gender   Gender `dynamodbav:"inmate_gender" json:"gender"`
}

type InmateWithConfirm struct {
	Inmate        Inmate
	InmateConfirm InmateConfirm
}

type InmateWithAttempt struct {
	Inmate        Inmate
	InmateAttempt InmateAttempt
}

func (i *Inmate) MarshalJSON() ([]byte, error) {
	type Alias Inmate
	return json.Marshal(&struct {
		ID string `json:"id"`
		Alias
	}{
		ID:    INMATE_ID_PREFIX + i.ID,
		Alias: (Alias)(*i),
	})
}

func (g *Gender) UnmarshalText(text []byte) error {
	s := strings.ToLower(string(text))
	switch s {
	case "male", "1":
		*g = Male
	case "female", "0":
		*g = Female
	default:
		return fmt.Errorf("invalid gender: %s", text)
	}
	return nil
}
