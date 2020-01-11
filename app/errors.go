package app

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Path   string
	Reason string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("Failed to validate %s due to %s", e.Path, e.Reason)
}

type AggregatedValidationError struct {
	errors []ValidationError
}

func NewAggregatedValidationError() AggregatedValidationError {
	err := AggregatedValidationError{}
	err.errors = make([]ValidationError, 0)
	return err
}

func (e *AggregatedValidationError) AddError(validationError ValidationError) {
	if e.errors == nil {
		e.errors = make([]ValidationError, 0)
	}
	e.errors = append(e.errors, validationError)
}

func (e AggregatedValidationError) Errors() []ValidationError {
	return e.errors
}

func (e AggregatedValidationError) Error() string {
	sb := strings.Builder{}
	if e.errors == nil {
		return ""
	}
	for _, e := range e.errors {
		sb.WriteString(fmt.Sprintf("error at %s \n", e.Path))
	}
	return sb.String()
}
