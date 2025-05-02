// internal/pixcelbattle/domain/validation.go
package domain

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	parts := make([]string, len(ve))
	for i, e := range ve {
		parts[i] = e.Error()
	}
	return strings.Join(parts, "; ")
}

func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}
