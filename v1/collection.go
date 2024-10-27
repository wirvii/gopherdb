package gopherdb

// Collection is a struct that holds the collection information.
type Collection struct {
	name     string
	database *Database
}

// Name is a function that returns the collection name.
func (c *Collection) Name() string {
	return c.name
}

// Database is a function that returns the database.
func (c *Collection) Database() *Database {
	return c.database
}
