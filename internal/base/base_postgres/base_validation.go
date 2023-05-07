package base_postgres

import (
	"application_template/pkg/types"
	"encoding/json"
	"gorm.io/gorm/utils"
	"io"
)

func ValidateBody(model HasId, requestBody io.ReadCloser) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	err := json.NewDecoder(requestBody).Decode(&m)
	if err != nil {
		return nil, err
	}

	for k, v := range m {

		value := v
		if k == "PhoneNumber" {
			switch phoneNumber := value.(type) {
			case types.PhoneNumber:
				value = phoneNumber.Number
			}
		}

		m[k] = value

		if utils.Contains(model.GetReadOnlyFields(), k) {
			delete(m, k)
		}
	}

	return m, nil
}
