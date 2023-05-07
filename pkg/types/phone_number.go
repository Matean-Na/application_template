package types

import (
	"application_template/utils"
	"database/sql/driver"
	"encoding/json"
	"regexp"
)

type PhoneNumber struct {
	Number string
}

func isValidKyrgyzPhoneNumber(number string) bool {
	pattern := `(?:^| +)(?:\( *)?(?:(?:\+?996|0) *(?:- *)?)?(?:\( *)?(\d{3})(?: *\))? *(?:- *)?((?:\d *(?:- *)?){6})(?: +|$)`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(number)
}

func NewPhoneNumber(number string, validate bool) (*PhoneNumber, error) {
	if validate {
		if !(isValidKyrgyzPhoneNumber(number)) {
			return nil, utils.LocalizeError{
				Source:  nil,
				Message: "exception:wrong-phone-number-format",
			}
		}
	}

	pattern := `[^0-9]+`
	r := regexp.MustCompile(pattern)
	number = r.ReplaceAllString(number, "")
	return &PhoneNumber{
		Number: number,
	}, nil
}

func (p *PhoneNumber) Scan(value interface{}) error {
	number, ok := value.(string)
	if !ok {
		return utils.LocalizeError{
			Source:  nil,
			Message: "exception:failed-to-parse",
			Data: map[string]interface{}{
				"Value": value,
				"Type":  "string",
			},
		}
	}

	*p = PhoneNumber{
		Number: number,
	}
	return nil
}

func (p PhoneNumber) Value() (driver.Value, error) {
	return p.Number, nil
}

func (p *PhoneNumber) UnmarshalJSON(bytes []byte) error {
	var number string
	var baseError = utils.NewLocalizeError(nil, "exception:failed-to-unmarshall-phone-number", nil)

	// handle string
	err := json.Unmarshal(bytes, &number)
	if err != nil {
		// handle object
		var phoneNumber map[string]string
		err := json.Unmarshal(bytes, &phoneNumber)
		if err != nil {
			return err
		}

		if number, exists := phoneNumber["Number"]; exists {
			*p = PhoneNumber{
				Number: number,
			}
			return nil
		}
		return baseError
	}

	*p = PhoneNumber{
		Number: number,
	}

	return nil
}

func (p *PhoneNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Number)
}
