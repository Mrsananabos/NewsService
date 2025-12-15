package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

//go:generate reform
//reform:news
type News struct {
	ID      int64  `json:"Id" reform:"id,pk"`
	Title   string `json:"Title" reform:"title"`
	Content string `json:"Content" reform:"content"`
}

type NewsWithCategories struct {
	News
	Categories []int64 `json:"Categories"`
}

type NewsEditForm struct {
	Title      *string  `json:"Title" validate:"omitempty,min=1,max=255"`
	Content    *string  `json:"Content" validate:"omitempty,min=1"`
	Categories *[]int64 `json:"Categories" validate:"omitempty,dive,gt=0"`
}

type NewsCreateForm struct {
	Title      string   `json:"Title" validate:"required,min=1,max=255"`
	Content    string   `json:"Content" validate:"required,min=1"`
	Categories *[]int64 `json:"Categories" validate:"omitempty,dive,gt=0"`
}

var validate = validator.New()

func (n *NewsCreateForm) Validate() error {
	if err := validate.Struct(n); err != nil {
		return formatValidationError(err)
	}
	return nil
}

func (n *NewsCreateForm) Normalize() {
	n.Title = strings.TrimSpace(n.Title)
	n.Content = strings.TrimSpace(n.Content)
}

func (n *NewsEditForm) Validate() error {
	if n.Title == nil && n.Content == nil && n.Categories == nil {
		return errors.New("body must contain at least one field to update (Title, Content, or Categories)")
	}

	if err := validate.Struct(n); err != nil {
		return formatValidationError(err)
	}
	return nil
}

func (n *NewsEditForm) Normalize() {
	if n.Title != nil {
		trimmed := strings.TrimSpace(*n.Title)
		n.Title = &trimmed
	}

	if n.Content != nil {
		trimmed := strings.TrimSpace(*n.Content)
		n.Content = &trimmed
	}
}

func formatValidationError(err error) error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				return fmt.Errorf("%s: field is required", e.Field())
			case "min":

				return fmt.Errorf("%s: minimum length is %s", e.Field(), e.Param())
			case "max":
				return fmt.Errorf("%s: maximum length is %s", e.Field(), e.Param())
			case "gt":
				return fmt.Errorf("%s: must be greater than %s", e.Field(), e.Param())
			case "dive":
				return fmt.Errorf("%s: contains invalid element", e.Field())
			default:
				return fmt.Errorf("%s: validation failed (%s)", e.Field(), e.Tag())
			}
		}
	}
	return err
}
