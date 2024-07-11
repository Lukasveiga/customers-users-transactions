package shared

import "fmt"

type AlreadyExistsError struct {
	Object string
	Id interface{}
}

func (e AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s already exists with id %v", e.Object, e.Id)
}