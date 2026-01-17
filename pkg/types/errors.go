// Package types contains error definitions for GibRAM
package types

import (
	"fmt"
)

// ErrorCode represents a structured error code
type ErrorCode int

const (
	// Success
	ErrOK ErrorCode = 0

	// Client errors (1xxx)
	ErrBadRequest      ErrorCode = 1000
	ErrUnauthorized    ErrorCode = 1001
	ErrForbidden       ErrorCode = 1002
	ErrNotFound        ErrorCode = 1003
	ErrConflict        ErrorCode = 1004 // duplicate
	ErrRateLimited     ErrorCode = 1005
	ErrPayloadTooLarge ErrorCode = 1006
	ErrInvalidInput    ErrorCode = 1007

	// Server errors (2xxx)
	ErrInternal    ErrorCode = 2000
	ErrUnavailable ErrorCode = 2001
	ErrTimeout     ErrorCode = 2002
	ErrShuttingDown ErrorCode = 2003

	// Data errors (3xxx)
	ErrInvalidVector    ErrorCode = 3000
	ErrInvalidEntity    ErrorCode = 3001
	ErrInvalidQuery     ErrorCode = 3002
	ErrCorruptedData    ErrorCode = 3003
	ErrInvalidDocument  ErrorCode = 3004
	ErrInvalidTextUnit  ErrorCode = 3005
	ErrInvalidRelation  ErrorCode = 3006
	ErrInvalidCommunity ErrorCode = 3007
)

// String returns the string representation of the error code
func (e ErrorCode) String() string {
	switch e {
	case ErrOK:
		return "OK"
	case ErrBadRequest:
		return "BAD_REQUEST"
	case ErrUnauthorized:
		return "UNAUTHORIZED"
	case ErrForbidden:
		return "FORBIDDEN"
	case ErrNotFound:
		return "NOT_FOUND"
	case ErrConflict:
		return "CONFLICT"
	case ErrRateLimited:
		return "RATE_LIMITED"
	case ErrPayloadTooLarge:
		return "PAYLOAD_TOO_LARGE"
	case ErrInvalidInput:
		return "INVALID_INPUT"
	case ErrInternal:
		return "INTERNAL_ERROR"
	case ErrUnavailable:
		return "UNAVAILABLE"
	case ErrTimeout:
		return "TIMEOUT"
	case ErrShuttingDown:
		return "SHUTTING_DOWN"
	case ErrInvalidVector:
		return "INVALID_VECTOR"
	case ErrInvalidEntity:
		return "INVALID_ENTITY"
	case ErrInvalidQuery:
		return "INVALID_QUERY"
	case ErrCorruptedData:
		return "CORRUPTED_DATA"
	case ErrInvalidDocument:
		return "INVALID_DOCUMENT"
	case ErrInvalidTextUnit:
		return "INVALID_TEXTUNIT"
	case ErrInvalidRelation:
		return "INVALID_RELATION"
	case ErrInvalidCommunity:
		return "INVALID_COMMUNITY"
	default:
		return "UNKNOWN"
	}
}

// GibRAMError represents a structured error response
type GibRAMError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

// Error implements the error interface
func (e *GibRAMError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code.String(), e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code.String(), e.Message)
}

// NewError creates a new GibRAMError
func NewError(code ErrorCode, message string) *GibRAMError {
	return &GibRAMError{
		Code:    code,
		Message: message,
	}
}

// NewErrorWithDetails creates a new GibRAMError with details
func NewErrorWithDetails(code ErrorCode, message, details string) *GibRAMError {
	return &GibRAMError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// IsClientError returns true if the error is a client error (1xxx)
func (e *GibRAMError) IsClientError() bool {
	return e.Code >= 1000 && e.Code < 2000
}

// IsServerError returns true if the error is a server error (2xxx)
func (e *GibRAMError) IsServerError() bool {
	return e.Code >= 2000 && e.Code < 3000
}

// IsDataError returns true if the error is a data error (3xxx)
func (e *GibRAMError) IsDataError() bool {
	return e.Code >= 3000 && e.Code < 4000
}

// Common error instances for reuse
var (
	ErrEntityNotFound      = NewError(ErrNotFound, "Entity not found")
	ErrDocumentNotFound    = NewError(ErrNotFound, "Document not found")
	ErrTextUnitNotFound    = NewError(ErrNotFound, "TextUnit not found")
	ErrRelationNotFound    = NewError(ErrNotFound, "Relationship not found")
	ErrCommunityNotFound   = NewError(ErrNotFound, "Community not found")
	ErrDuplicateExternalID = NewError(ErrConflict, "External ID already exists")
	ErrDuplicateTitle      = NewError(ErrConflict, "Title already exists")
	ErrVectorDimMismatch   = NewError(ErrInvalidVector, "Vector dimension mismatch")
	ErrEmptyVector         = NewError(ErrInvalidVector, "Vector cannot be empty")
	ErrServerShuttingDown  = NewError(ErrShuttingDown, "Server is shutting down")
)
