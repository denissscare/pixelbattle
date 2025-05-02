package domain

import (
	"fmt"
	"regexp"
	"time"
)

var (
	MaxWidth  = 1000
	MaxHeight = 1000
	re        = regexp.MustCompile(`(?i)^#[0-9A-F]{6}$`)
)

type Pixel struct {
	X         int       `json:"x"`
	Y         int       `json:"y"`
	Color     string    `json:"color"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
}

func (p Pixel) Validate() error {
	var errs ValidationErrors

	if p.X < 0 || p.X >= MaxWidth {
		errs = append(errs, ValidationError{
			Field:   "X",
			Message: fmt.Sprintf("out of range [0,%d)", MaxWidth),
		})
	}
	if p.Y < 0 || p.Y >= MaxHeight {
		errs = append(errs, ValidationError{
			Field:   "Y",
			Message: fmt.Sprintf("out of range [0,%d)", MaxHeight),
		})
	}
	if !re.MatchString(p.Color) {
		errs = append(errs, ValidationError{
			Field:   "Color",
			Message: "must be hex #RRGGBB",
		})
	}
	if p.Author == "" {
		errs = append(errs, ValidationError{
			Field:   "Author",
			Message: "cannot be empty",
		})
	}
	if p.Timestamp.IsZero() {
		errs = append(errs, ValidationError{
			Field:   "Timestamp",
			Message: "must be set",
		})
	} else if p.Timestamp.After(time.Now().Add(time.Minute)) {
		errs = append(errs, ValidationError{
			Field:   "Timestamp",
			Message: "too far in the future",
		})
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
