package h

import "fmt"

type FunctionalError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *FunctionalError) Error() string {
	return fmt.Sprintf("FunctionalError %d: %s", e.Code, e.Message)
}

func NewFunctionalError(message string, code ...string) error {
	if len(code) > 0 {
		return &FunctionalError{Code: code[0], Message: message}
	} else {
		return &FunctionalError{Message: message}
	}
}

/// ----------------------------------------------------------------------------------------------------------------

type ResourceNotFoundError struct {
	Message string `json:"message,omitempty"`
}

func (e *ResourceNotFoundError) Error() string {
	return fmt.Sprintf("FunctionalError %d", e.Message)
}

func NewResourceNotFoundError(message string) error {
	return &ResourceNotFoundError{Message: message}
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

type UnauthorizedError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("FunctionalError %s", e.Message)
}

func NewUnauthorizedError(code string, message string) error {
	return &UnauthorizedError{Code: code, Message: message}
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
