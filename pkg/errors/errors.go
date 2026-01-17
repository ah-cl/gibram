// Package errors provides standardized error types for GibRAM
package errors

import (
	"errors"
	"fmt"
)

// Error codes for GibRAM operations
const (
	// General errors
	CodeInternal        = "INTERNAL_ERROR"
	CodeInvalidInput    = "INVALID_INPUT"
	CodeNotFound        = "NOT_FOUND"
	CodeAlreadyExists   = "ALREADY_EXISTS"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeForbidden       = "FORBIDDEN"
	CodeTimeout         = "TIMEOUT"
	CodeCancelled       = "CANCELLED"

	// Resource errors
	CodeQuotaExceeded   = "QUOTA_EXCEEDED"
	CodeResourceExhausted = "RESOURCE_EXHAUSTED"
	CodeOutOfMemory     = "OUT_OF_MEMORY"
	CodeTooManyRequests = "TOO_MANY_REQUESTS"

	// Data integrity errors
	CodeCorruption      = "DATA_CORRUPTION"
	CodeChecksumMismatch = "CHECKSUM_MISMATCH"
	CodeVersionMismatch = "VERSION_MISMATCH"
	CodeIntegrityViolation = "INTEGRITY_VIOLATION"

	// Concurrency errors
	CodeDeadlock        = "DEADLOCK"
	CodeConflict        = "CONFLICT"
	CodeRetryable       = "RETRYABLE"

	// Vector errors
	CodeDimensionMismatch = "DIMENSION_MISMATCH"
	CodeInvalidVector   = "INVALID_VECTOR"
	CodeIndexFull       = "INDEX_FULL"

	// Session errors
	CodeSessionNotFound = "SESSION_NOT_FOUND"
	CodeSessionExpired  = "SESSION_EXPIRED"
	CodeSessionLimitExceeded = "SESSION_LIMIT_EXCEEDED"

	// Graph errors
	CodeEntityNotFound = "ENTITY_NOT_FOUND"
	CodeRelationshipNotFound = "RELATIONSHIP_NOT_FOUND"
	CodeCommunityNotFound = "COMMUNITY_NOT_FOUND"
	CodeCyclicDependency = "CYCLIC_DEPENDENCY"

	// Backup/Recovery errors
	CodeBackupFailed    = "BACKUP_FAILED"
	CodeRestoreFailed   = "RESTORE_FAILED"
	CodeSnapshotCorrupt = "SNAPSHOT_CORRUPT"
	CodeWALCorrupt      = "WAL_CORRUPT"
)

// GibRAMError is the standard error type with code and context
type GibRAMError struct {
	Code    string
	Message string
	Cause   error
	Context map[string]interface{}
}

// Error implements the error interface
func (e *GibRAMError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap implements error unwrapping
func (e *GibRAMError) Unwrap() error {
	return e.Cause
}

// Is implements error comparison
func (e *GibRAMError) Is(target error) bool {
	t, ok := target.(*GibRAMError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// WithContext adds context to the error
func (e *GibRAMError) WithContext(key string, value interface{}) *GibRAMError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// New creates a new GibRAMError
func New(code, message string) *GibRAMError {
	return &GibRAMError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an error with a code and message
func Wrap(err error, code, message string) *GibRAMError {
	return &GibRAMError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

// Common error constructors

// ErrInternal creates an internal error
func ErrInternal(message string, cause error) *GibRAMError {
	return &GibRAMError{Code: CodeInternal, Message: message, Cause: cause}
}

// ErrInvalidInput creates an invalid input error
func ErrInvalidInput(message string) *GibRAMError {
	return &GibRAMError{Code: CodeInvalidInput, Message: message}
}

// ErrNotFound creates a not found error
func ErrNotFound(resource, id string) *GibRAMError {
	return &GibRAMError{
		Code:    CodeNotFound,
		Message: fmt.Sprintf("%s not found: %s", resource, id),
		Context: map[string]interface{}{"resource": resource, "id": id},
	}
}

// ErrAlreadyExists creates an already exists error
func ErrAlreadyExists(resource, id string) *GibRAMError {
	return &GibRAMError{
		Code:    CodeAlreadyExists,
		Message: fmt.Sprintf("%s already exists: %s", resource, id),
		Context: map[string]interface{}{"resource": resource, "id": id},
	}
}

// ErrQuotaExceeded creates a quota exceeded error
func ErrQuotaExceeded(quota string, limit, current int) *GibRAMError {
	return &GibRAMError{
		Code:    CodeQuotaExceeded,
		Message: fmt.Sprintf("%s quota exceeded: %d/%d", quota, current, limit),
		Context: map[string]interface{}{"quota": quota, "limit": limit, "current": current},
	}
}

// ErrResourceExhausted creates a resource exhausted error
func ErrResourceExhausted(resource string) *GibRAMError {
	return &GibRAMError{
		Code:    CodeResourceExhausted,
		Message: fmt.Sprintf("%s exhausted", resource),
		Context: map[string]interface{}{"resource": resource},
	}
}

// ErrOutOfMemory creates an out of memory error
func ErrOutOfMemory(requested, available int64) *GibRAMError {
	return &GibRAMError{
		Code:    CodeOutOfMemory,
		Message: fmt.Sprintf("out of memory: requested %d bytes, available %d bytes", requested, available),
		Context: map[string]interface{}{"requested": requested, "available": available},
	}
}

// ErrCorruption creates a data corruption error
func ErrCorruption(message string, cause error) *GibRAMError {
	return &GibRAMError{Code: CodeCorruption, Message: message, Cause: cause}
}

// ErrChecksumMismatch creates a checksum mismatch error
func ErrChecksumMismatch(expected, actual uint64) *GibRAMError {
	return &GibRAMError{
		Code:    CodeChecksumMismatch,
		Message: fmt.Sprintf("checksum mismatch: expected %d, got %d", expected, actual),
		Context: map[string]interface{}{"expected": expected, "actual": actual},
	}
}

// ErrDimensionMismatch creates a dimension mismatch error
func ErrDimensionMismatch(expected, actual int) *GibRAMError {
	return &GibRAMError{
		Code:    CodeDimensionMismatch,
		Message: fmt.Sprintf("dimension mismatch: expected %d, got %d", expected, actual),
		Context: map[string]interface{}{"expected": expected, "actual": actual},
	}
}

// ErrSessionNotFound creates a session not found error
func ErrSessionNotFound(sessionID string) *GibRAMError {
	return &GibRAMError{
		Code:    CodeSessionNotFound,
		Message: fmt.Sprintf("session not found: %s", sessionID),
		Context: map[string]interface{}{"session_id": sessionID},
	}
}

// ErrSessionExpired creates a session expired error
func ErrSessionExpired(sessionID string) *GibRAMError {
	return &GibRAMError{
		Code:    CodeSessionExpired,
		Message: fmt.Sprintf("session expired: %s", sessionID),
		Context: map[string]interface{}{"session_id": sessionID},
	}
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	var gibramErr *GibRAMError
	if errors.As(err, &gibramErr) {
		return gibramErr.Code == CodeRetryable ||
			gibramErr.Code == CodeTimeout ||
			gibramErr.Code == CodeTooManyRequests ||
			gibramErr.Code == CodeDeadlock ||
			gibramErr.Code == CodeConflict
	}
	return false
}

// IsNotFound checks if an error is a not found error
func IsNotFound(err error) bool {
	var gibramErr *GibRAMError
	if errors.As(err, &gibramErr) {
		return gibramErr.Code == CodeNotFound ||
			gibramErr.Code == CodeSessionNotFound ||
			gibramErr.Code == CodeEntityNotFound ||
			gibramErr.Code == CodeRelationshipNotFound ||
			gibramErr.Code == CodeCommunityNotFound
	}
	return false
}

// IsQuotaExceeded checks if an error is a quota exceeded error
func IsQuotaExceeded(err error) bool {
	var gibramErr *GibRAMError
	if errors.As(err, &gibramErr) {
		return gibramErr.Code == CodeQuotaExceeded ||
			gibramErr.Code == CodeResourceExhausted ||
			gibramErr.Code == CodeOutOfMemory ||
			gibramErr.Code == CodeSessionLimitExceeded
	}
	return false
}

// IsDataIntegrityError checks if an error is a data integrity error
func IsDataIntegrityError(err error) bool {
	var gibramErr *GibRAMError
	if errors.As(err, &gibramErr) {
		return gibramErr.Code == CodeCorruption ||
			gibramErr.Code == CodeChecksumMismatch ||
			gibramErr.Code == CodeVersionMismatch ||
			gibramErr.Code == CodeIntegrityViolation ||
			gibramErr.Code == CodeSnapshotCorrupt ||
			gibramErr.Code == CodeWALCorrupt
	}
	return false
}
