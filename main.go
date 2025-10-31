package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/inmates", getInmates)

	router.Run()
}

type Inmate struct {
	ID          uint64    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	MiddleName  string    `json:"middle_name"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
	PhoneNumber string    `json:"phone_number"`
}

var Inmates = []Inmate{
	{
		ID:          1001,
		FirstName:   "John",
		LastName:    "Doe",
		MiddleName:  "Michael",
		UpdatedAt:   time.Now(),
		CreatedAt:   time.Date(2024, 3, 10, 9, 0, 0, 0, time.UTC),
		PhoneNumber: "+1-555-1001",
	},
	{
		ID:          1002,
		FirstName:   "Alice",
		LastName:    "Johnson",
		MiddleName:  "Marie",
		UpdatedAt:   time.Now(),
		CreatedAt:   time.Date(2024, 5, 21, 10, 15, 0, 0, time.UTC),
		PhoneNumber: "+1-555-1002",
	},
	{
		ID:          1003,
		FirstName:   "Robert",
		LastName:    "Smith",
		MiddleName:  "James",
		UpdatedAt:   time.Now(),
		CreatedAt:   time.Date(2024, 7, 5, 13, 45, 0, 0, time.UTC),
		PhoneNumber: "+1-555-1003",
	},
	{
		ID:          1004,
		FirstName:   "Maria",
		LastName:    "Lopez",
		MiddleName:  "Isabel",
		UpdatedAt:   time.Now(),
		CreatedAt:   time.Date(2024, 8, 18, 11, 30, 0, 0, time.UTC),
		PhoneNumber: "+1-555-1004",
	},
	{
		ID:          1005,
		FirstName:   "David",
		LastName:    "Brown",
		MiddleName:  "Andrew",
		UpdatedAt:   time.Now(),
		CreatedAt:   time.Date(2024, 9, 2, 16, 0, 0, 0, time.UTC),
		PhoneNumber: "+1-555-1005",
	},
}

func getInmates(c *gin.Context) {
	c.JSON(http.StatusOK, Inmates)
}
