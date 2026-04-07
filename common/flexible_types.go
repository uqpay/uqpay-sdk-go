package common

import "encoding/json"

// FlexibleString accepts both JSON string and JSON number, storing the value as a string.
// Used for API fields that inconsistently return string or number types.
type FlexibleString string

func (s *FlexibleString) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*s = ""
		return nil
	}
	// JSON string: strip quotes
	if data[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		*s = FlexibleString(str)
		return nil
	}
	// JSON number or anything else: store raw
	*s = FlexibleString(data)
	return nil
}

func (s FlexibleString) String() string { return string(s) }

// FlexibleStringMap is a map[string]string that tolerates non-object JSON values
// (e.g. "", null, []) returned by the API when the field is empty.
type FlexibleStringMap map[string]string

func (m *FlexibleStringMap) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*m = nil
		return nil
	}
	if data[0] == '"' || data[0] == '[' {
		*m = nil
		return nil
	}
	type Alias map[string]string
	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		*m = nil
		return nil
	}
	*m = FlexibleStringMap(alias)
	return nil
}
