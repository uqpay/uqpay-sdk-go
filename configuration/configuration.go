package configuration

import "net/http"

// Environment represents the UQPAY API environment
type Environment struct {
	BaseURL      string
	FilesBaseURL string
}

// Sandbox returns the sandbox environment
func Sandbox() *Environment {
	return &Environment{
		BaseURL:      "https://api-sandbox.uqpaytech.com/api",
		FilesBaseURL: "https://files.uqpaytech.com/api",
	}
}

// Production returns the production environment
func Production() *Environment {
	return &Environment{
		BaseURL:      "https://api.uqpay.com/api",
		FilesBaseURL: "https://files.uqpay.com/api",
	}
}

// Configuration holds SDK configuration
type Configuration struct {
	ClientID    string
	APIKey      string
	Environment *Environment
	HTTPClient  *http.Client
}
