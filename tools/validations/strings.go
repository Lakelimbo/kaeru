package validations

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

// StringEq is a simple rule for ozzo-validation that will
// check if a string is the exact same as another
// (useful for, e.g password confirmation).
func StringEq(str string, errMsg string) validation.RuleFunc {
	return func(value any) error {
		s, _ := value.(string)
		if s != str {
			return errors.New(errMsg)
		}

		return nil
	}
}
