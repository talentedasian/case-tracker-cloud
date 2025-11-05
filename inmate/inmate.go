package inmate

type Gender uint8

const (
	Female Gender = 0
	Male   Gender = 1
)

type Inmate struct {
	ID       uint64 `dynamodbav:"inmate_id" json:"id"`
	LastName string `dynamodbav:"inmate_last_name" json:"last_name"`
	Gender   Gender `dynamodbav:"inmate_gender" json:"gender"`
}
