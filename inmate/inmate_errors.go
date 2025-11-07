package inmate

type UnableToParse struct {
	err string
}

type UnableToWriteItem struct {
	err string
}

func NewUnableToParseRequestError() *UnableToParse {
	return &UnableToParse{"failed to parse inmate into dynamodb request"}
}

func NewUnableToParseResponseError() *UnableToParse {
	return &UnableToParse{"failed to parse dynamodb response"}
}

func NewUnableToWriteItem() *UnableToWriteItem {
	return &UnableToWriteItem{"failed to write to dynamodb item"}
}

func (u *UnableToParse) Error() string {
	return u.err
}

func (u *UnableToWriteItem) Error() string {
	return u.err
}
