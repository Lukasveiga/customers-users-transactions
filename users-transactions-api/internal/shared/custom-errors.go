package shared

import "fmt"

type AlreadyExistsError struct {
	Object string
	Id     interface{}
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s already exists with id %v", e.Object, e.Id)
}

type ValidationError struct {
	Errors map[string]string
}

func (e *ValidationError) Error() string {
	return "Validation error"
}

func (e *ValidationError) AddError(field string, err string) {
	e.Errors[field] = err
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}

type EntityNotFoundError struct {
	Object string
	Id     interface{}
}

func (e *EntityNotFoundError) Error() string {
	return fmt.Sprintf("%s not found with id %v", e.Object, e.Id)
}

type InactiveAccountError struct{}

func (e *InactiveAccountError) Error() string {
	return "Card cannot be created to an inactive account"
}
