package validators

import (
	"encoding/json"
	"fmt"
	"math"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func ValidateCreateNewsRequest(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return &ValidationError{Field: "body", Message: "invalid JSON format"}
	}

	if len(raw) == 0 {
		return &ValidationError{
			Field:   "body",
			Message: "must contain required fields (Title, Content)",
		}
	}

	if title, exists := raw["Title"]; exists {
		if _, ok := title.(string); !ok {
			return &ValidationError{
				Field:   "Title",
				Message: "must be string",
			}
		}
	}

	if content, exists := raw["Content"]; exists {
		if _, ok := content.(string); !ok {
			return &ValidationError{
				Field:   "Content",
				Message: "must be string",
			}
		}
	}

	if categories, exists := raw["Categories"]; exists && categories != nil {
		if err := validateCategoriesArray(categories); err != nil {
			return err
		}
	}

	return nil
}

func ValidateEditNewsRequest(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return &ValidationError{Field: "body", Message: "invalid JSON format"}
	}

	if len(raw) == 0 {
		return &ValidationError{
			Field:   "body",
			Message: "must contain at least one field to update (Title, Content, or Categories)",
		}
	}

	if title, exists := raw["Title"]; exists {
		if title != nil {
			if _, ok := title.(string); !ok {
				return &ValidationError{
					Field:   "Title",
					Message: fmt.Sprintf("must be string, got %T", title),
				}
			}
		}
	}

	if content, exists := raw["Content"]; exists {
		if content != nil {
			if _, ok := content.(string); !ok {
				return &ValidationError{
					Field:   "Content",
					Message: fmt.Sprintf("must be string, got %T", content),
				}
			}
		}
	}

	if categories, exists := raw["Categories"]; exists {
		if categories != nil {
			if err := validateCategoriesArray(categories); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateCategoriesArray(categories interface{}) error {
	catSlice, ok := categories.([]interface{})
	if !ok {
		return &ValidationError{
			Field:   "Categories",
			Message: "must be array numbers",
		}
	}

	for i, cat := range catSlice {
		num, ok := cat.(float64)
		if !ok {
			return &ValidationError{
				Field:   "Categories",
				Message: fmt.Sprintf("element at index %d must be number, got %T", i, cat),
			}
		}

		if num != math.Floor(num) {
			return &ValidationError{
				Field:   "Categories",
				Message: fmt.Sprintf("element at index %d must be integer, got %v", i, num),
			}
		}

		if num <= 0 {
			return &ValidationError{
				Field:   "Categories",
				Message: fmt.Sprintf("element at index %d must be positive, got %v", i, num),
			}
		}
	}

	return nil
}
