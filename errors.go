package gopherdb

import "errors"

var (
	// ErrMissingFieldForIndex is returned when a field is missing for an index.
	ErrMissingFieldForIndex = errors.New("missing field for index")
	// ErrUniqueIndexViolation is returned when a unique index is violated.
	ErrUniqueIndexViolation = errors.New("unique index violation")
	// ErrIndexAlreadyExists is returned when an index already exists.
	ErrIndexAlreadyExists = errors.New("index already exists")
	// ErrInvalidValueType is returned when an invalid value type is used.
	ErrInvalidValueType = errors.New("invalid value type")
	// ErrMapTypeConversionFailed is returned when a map type conversion fails.
	ErrMapTypeConversionFailed = errors.New("map type conversion failed")
	// ErrUnsupportedTypeForMapConversion is returned when an unsupported type is used for map conversion.
	ErrUnsupportedTypeForMapConversion = errors.New("unsupported type for map conversion")
	// ErrMustBePointer is returned when a value is not a pointer.
	ErrMustBePointer = errors.New("must be a pointer")
)
