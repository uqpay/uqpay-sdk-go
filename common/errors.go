package common

import (
	"encoding/json"
	"fmt"
)

// APIError represents an API error response
type APIError struct {
	Code       FlexibleCode `json:"code"`    // Error code for programmatic handling, e.g. "invalid_number". May be string or number
	Message    string       `json:"message"` // Human-readable error description. Not intended for parsing; format may change
	StatusCode int          `json:"-"`       // HTTP status code: 400 (Bad Request), 401 (Unauthorized), 404 (Not Found), 500 (Server Error)
}

// FlexibleCode handles API code fields that may be string or number
type FlexibleCode string

func (c *FlexibleCode) UnmarshalJSON(data []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*c = FlexibleCode(s)
		return nil
	}
	// Try number
	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		*c = FlexibleCode(n.String())
		return nil
	}
	*c = FlexibleCode(string(data))
	return nil
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s (HTTP %d)", string(e.Code), e.Message, e.StatusCode)
}

// IsNotFound returns true if the error is a 404
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == 404
}

// IsUnauthorized returns true if the error is a 401
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == 401
}

// IsBadRequest returns true if the error is a 400
func (e *APIError) IsBadRequest() bool {
	return e.StatusCode == 400
}
