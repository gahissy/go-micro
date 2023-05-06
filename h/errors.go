package h

import "fmt"

type FunctionalError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *FunctionalError) Error() string {
	return fmt.Sprintf("FunctionalError %d: %s", e.Code, e.Message)
}

func NewFunctionalError(message string) error {
	return &FunctionalError{Message: message}
}

/// ----------------------------------------------------------------------------------------------------------------

type ForbiddenError struct {
	Message string `json:"message,omitempty"`
}

func (e *ForbiddenError) Error() string {
	return fmt.Sprintf("FunctionalError %s", e.Message)
}

func NewForbiddenError(message string) error {
	return &ForbiddenError{Message: message}
}

/// ----------------------------------------------------------------------------------------------------------------

/// ----------------------------------------------------------------------------------------------------------------

type TechnicalError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *TechnicalError) Error() string {
	return fmt.Sprintf("TechnicalError %d: %s", e.Code, e.Message)
}

func NewTechnicalError(code string, message string) error {
	return &TechnicalError{Code: code, Message: message}
}

/// ----------------------------------------------------------------------------------------------------------------
