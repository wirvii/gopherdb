package gopherdb

import "errors"

var (
	// ErrMissingFieldForIndex is returned when a field is missing for an index.
	ErrMissingFieldForIndex = errors.New("missing field for index")
	// ErrEmptyIndexFields is returned when an index has no fields.
	ErrEmptyIndexFields = errors.New("empty index fields")
	// ErrDuplicateIndexField is returned when a duplicate index field is found.
	ErrDuplicateIndexField = errors.New("duplicate index field")
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
	// ErrDocumentNotFound is returned when a document is not found.
	ErrDocumentNotFound = errors.New("document not found")
	// ErrDocumentIsNil is returned when a document is nil.
	ErrDocumentIsNil = errors.New("document is nil")
	// ErrDocumentSliceIsNil is returned when a document slice is nil.
	ErrDocumentSliceIsNil = errors.New("document slice is nil")
	// ErrDocumentSlicePointerIsNil is returned when a document slice pointer is nil.
	ErrDocumentSlicePointerIsNil = errors.New("document slice pointer is nil")
	// ErrDocumentSliceEmpty is returned when a document slice is empty.
	ErrDocumentSliceEmpty = errors.New("document slice is empty")
	// ErrDocumentSliceTypeInvalid is returned when a document slice type is invalid.
	ErrDocumentSliceTypeInvalid = errors.New("document slice type is invalid")
	// ErrDocumentSliceElementTypeInvalid is returned when a document slice element type is invalid.
	ErrDocumentSliceElementTypeInvalid = errors.New("document slice element type is invalid")
	// ErrDocumentPointerIsNil is returned when a document pointer is nil.
	ErrDocumentPointerIsNil = errors.New("document pointer is nil")
	// ErrDocumentTypeInvalid is returned when a document type is invalid.
	ErrDocumentTypeInvalid = errors.New("document type is invalid")
)
