package shared

import "fmt"

type EntityNotFoundError struct {
	Message string
}

func (e *EntityNotFoundError) Error() string {
	return fmt.Sprint(e.Message)
}
