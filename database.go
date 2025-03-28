package gopherdb

import "github.com/wirvii/gopherdb/internal/storage"

// Database es una base de datos.
type Database struct {
	name    string
	storage storage.Storage
}

// NewDatabase crea una nueva instancia de Database.
func NewDatabase(name, path string) (*Database, error) {
	engine, err := storage.NewStorage(path)
	if err != nil {
		return nil, err
	}

	return &Database{
		name:    name,
		storage: engine,
	}, nil
}

// Collection devuelve una instancia de Collection para la base de datos
func (db *Database) Collection(name string) (*Collection, error) {
	return newCollection(db.storage, db.name, name)
}

// Close cierra la base de datos
func (db *Database) Close() error {
	return db.storage.Close()
}
