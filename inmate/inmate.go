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

const ID_PREFIX = "i#"

type Inmate struct {
	ID       string `dynamodbav:"partition_key" json:"id"`
	LastName string `dynamodbav:"sort_key" json:"last_name"`
	Gender   Gender `dynamodbav:"inmate_gender" json:"gender"`
}

func (i Inmate) MarshalJSON() ([]byte, error) {
	type Alias Inmate
	return json.Marshal(&struct {
		ID string `json:"id"`
		Alias
	}{
		ID:    ID_PREFIX + i.ID,
		Alias: (Alias)(i),
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
