package shared

import "fmt"

type AlreadyExistsError struct {
	Object string
	Id int32
}

func (e AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s already exists with id %d", e.Object, e.Id)
}