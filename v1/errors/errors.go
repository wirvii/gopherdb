package errors

import "errors"

var (
	// ErrDatabaseNotFound is an error that is returned when the database is not found.
	ErrDatabaseNotFound = errors.New("database not found")
	// ErrDatabaseAlreadyExists is an error that is returned when the database already exists.
	ErrDatabaseAlreadyExists = errors.New("database already exists")
	// ErrInvalidDatabaseName is an error that is returned when the database name is invalid.
	ErrInvalidDatabaseName = errors.New("invalid database name")
	// ErrInvalidCollectionName is an error that is returned when the collection name is invalid.
	ErrInvalidCollectionName = errors.New("invalid collection name")
	// ErrCollectionNotFound is an error that is returned when the collection is not found.
	ErrCollectionNotFound = errors.New("collection not found")
	// ErrCollectionAlreadyExists is an error that is returned when the collection already exists.
	ErrCollectionAlreadyExists = errors.New("collection already exists")
)

// IsDatabaseNotFound is a function that checks if the error is a database not found error.
func IsDatabaseNotFound(err error) bool {
	return errors.Is(err, ErrDatabaseNotFound)
}

// IsDatabaseAlreadyExists is a function that checks if the error is a database already exists error.
func IsDatabaseAlreadyExists(err error) bool {
	return errors.Is(err, ErrDatabaseAlreadyExists)
}

// IsInvalidDatabaseName is a function that checks if the error is an invalid database name error.
func IsInvalidDatabaseName(err error) bool {
	return errors.Is(err, ErrInvalidDatabaseName)
}

// IsInvalidCollectionName is a function that checks if the error is an invalid collection name error.
func IsInvalidCollectionName(err error) bool {
	return errors.Is(err, ErrInvalidCollectionName)
}

// IsCollectionNotFound is a function that checks if the error is a collection not found error.
func IsCollectionNotFound(err error) bool {
	return errors.Is(err, ErrCollectionNotFound)
}

// IsCollectionAlreadyExists is a function that checks if the error is a collection already exists error.
func IsCollectionAlreadyExists(err error) bool {
	return errors.Is(err, ErrCollectionAlreadyExists)
}
